package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"flag"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"main/controller"

	"main/pkg/base"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var (
	json     = jsoniter.ConfigCompatibleWithStandardLibrary
	listener net.Listener
	graceful = flag.Bool("graceful", false, "listen on fd open 3 (internal use only)")
)

func main() {
	// 建议程序开启多核支持
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd := os.Args
	if len(cmd) > 2 {
		switch cmd[1] {
		case "cli":
			var f = map[string]func(){
				"Test": Test,
			}
			f[cmd[2]]()
			break
		}
		return
	} else {
		srv := fasthttp.Server{
			Handler: index,
		}

		var err error
		flag.Parse()
		if *graceful {
			base.Log("Listening to existing file descriptor 3.")
			// cmd.ExtraFiles: If non-nil, entry i becomes file descriptor 3+i.
			// when we put socket FD at the first entry, it will always be 3(0+3)
			f := os.NewFile(3, "")
			listener, err = net.FileListener(f)
		} else {
			base.Log("Listening on a new file descriptor.")
			if _, ok := base.Config["tls"]; ok {
				tlsc := base.Config["tls"].([]interface{})
				tlsConfig := &tls.Config{}
				tlsConfig.Certificates = make([]tls.Certificate, len(tlsc))
				for k, v := range tlsc {
					t := v.(map[string]interface{})
					tlsConfig.Certificates[k], _ = tls.LoadX509KeyPair(t["certFile"].(string), t["keyFile"].(string))
				}
				tlsConfig.BuildNameToCertificate()
				listener, err = tls.Listen("tcp", base.Config["tlsAddress"].(string), tlsConfig)
			} else {
				listener, err = net.Listen("tcp", base.Config["netAddress"].(string))
			}
		}

		if err != nil {
			base.Log("listener error:=" + err.Error())
		}

		go func() { srv.Serve(listener) }()

		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
		sig := <-c
		switch sig {
		case syscall.SIGUSR2:
			base.Log("http reload")
			err := reload()
			if err != nil {
				base.Log("http reload err:" + err.Error())
			}
			break
		default:
			base.Log("http shutdown")
			break
		}
		if err := srv.Shutdown(); err != nil {
			base.Log("http shutdown err:" + err.Error())
		}
	}
}

func index(ctx *fasthttp.RequestCtx) {
	u := string(ctx.RequestURI())
	if _, ok := base.Config["debug"]; ok {
		if !base.Config["debug"].(bool) {
			defer func() {
				//全局异常处理
				if err := recover(); err != nil {
					base.Log("系统异常:" + err.(string) + "||uri:" + u + "||data:" + ctx.PostArgs().String())
				}
			}()
		}
	}
	var c = map[string]func(ctx *fasthttp.RequestCtx){
		"/": controller.Index.Default,
	}
	if _, ok := c[u]; ok {
		c[u](ctx)
	} else {
		ctx.SetStatusCode(403)
	}
}

type ViewFunc func(*fasthttp.RequestCtx)

func BasicAuth(f ViewFunc, user, passwd []byte) ViewFunc {
	return func(ctx *fasthttp.RequestCtx) {
		basicAuthPrefix := "Basic "

		// 获取 request header
		auth := string(ctx.Request.Header.Peek("Authorization"))
		// 如果是 http basic auth
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// 解码认证信息
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):],
			)
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 && bytes.Equal(pair[0], user) &&
					bytes.Equal(pair[1], passwd) {
					// 执行被装饰的函数
					f(ctx)
					return
				}
			}
		}

		// 认证失败，提示 401 Unauthorized
		// Restricted 可以改成其他的值
		ctx.Response.Header.Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// 401 状态码
		ctx.SetStatusCode(401)
	}
}

func reload() error {
	tl, ok := listener.(*net.TCPListener)
	if !ok {
		return errors.New("listener is not tcp listener")
	}

	f, err := tl.File()
	if err != nil {
		return err
	}

	args := []string{"-graceful"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// put socket FD at the first entry
	cmd.ExtraFiles = []*os.File{f}
	return cmd.Start()
}
