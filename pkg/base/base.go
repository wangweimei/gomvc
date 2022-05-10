package base

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"main/pkg/database"
)

var (
	DB        *gorm.DB
	Redis     *redigo.Pool
	json      = jsoniter.ConfigCompatibleWithStandardLibrary
	Config    map[string]interface{}
	dbFactory = database.Factory{}
)

func init() {
	//读取配置文件
	filepath := "./config.json"
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	c := make(map[string]interface{})
	json.Unmarshal(f, &c)
	Config = c

	DB = dbFactory.Init(c)

	//redis
	if _, ok := c["redis"]; ok {
		r := c["redis"].(map[string]interface{})
		pool := redigo.NewPool(func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", r["host"].(string)+":"+r["port"].(string))
			if err != nil {
				return nil, err
			}
			if r["pwd"] != "" {
				if _, err := c.Do("AUTH", r["pwd"].(string)); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		}, int(r["maxIdle"].(float64)))
		Redis = pool
	}
}

func Log(v ...interface{}) {
	//日志
	now := time.Now()
	logpath := "./log/" + now.Format("2006") + "/" + now.Format("01") + "/"
	createFile(logpath)
	file := logpath + now.Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	Log := log.New(logFile, "", log.LstdFlags|log.Lshortfile) // 将文件设置为loger作为输出
	defer logFile.Close()
	Log.Output(2, fmt.Sprintln(v...))
}

func RedisGet(key string) (d string) {
	c := Redis.Get()
	d, _ = redigo.String(c.Do("GET", key))
	c.Close()
	return
}

func PushQueue(queueName string, data string) (err error) {
	con := Redis.Get()
	defer con.Close()
	_, err = con.Do("lpush", queueName, data)
	return
}

//调用os.MkdirAll递归创建文件夹
func createFile(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//加载redis队列
func LoadQueue(queueFunc map[string]func(string)) {
	if _, ok := Config["redis"]; !ok {
		return
	}
	for ke, va := range queueFunc {
		k := ke
		v := va
		go func() {
			for {
				con := Redis.Get()
				defer con.Close()
				nameAndData, err := redigo.Strings(con.Do("brpop", k, 0))
				if err != nil {
					if err == redigo.ErrNil {
						err = nil
						continue
					}
					continue
				}
				if len(nameAndData) > 1 {
					data := nameAndData[1]
					v(data)
				}
			}
		}()
	}
}
