package router

import (
	"spotify-color/router/handlers"
	"spotify-color/ws"

	"github.com/valyala/fasthttp"
)

// RequestHandler is a simple router
func RequestHandler(ctx *fasthttp.RequestCtx) {
	setHeaders(ctx)
	switch string(ctx.Path()) {
	case "/get-color":
		handlers.ColorPickerHandler(ctx)
	case "/ws":
		ws.ServeWs(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func setHeaders(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
