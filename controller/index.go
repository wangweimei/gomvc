package controller

import (
	"bytes"
	"main/model"

	template "main/view/index/default"

	"github.com/valyala/fasthttp"
)

type IndexController struct{}

var Index IndexController

func (t *IndexController) Default(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html;charset=utf-8")

	r := model.Sample.GetList()

	buffer := new(bytes.Buffer)
	template.Body(r, buffer)
	ctx.Write(buffer.Bytes())
}
