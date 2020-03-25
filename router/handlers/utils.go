package handlers

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type responseInterface interface{}

func JSONResponse(ctx *fasthttp.RequestCtx, r responseInterface) {
	ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))

	if err := json.NewEncoder(ctx).Encode(r); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
