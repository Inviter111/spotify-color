package handlers

import (
	"fmt"

	"spotify-color/colorpicker"

	"github.com/valyala/fasthttp"
)

type colorResponse struct {
	CoverColor string `json:"hex"`
}

// ColorPickerHandler handle /getColor
func ColorPickerHandler(ctx *fasthttp.RequestCtx) {
	url := string(ctx.URI().QueryArgs().Peek("url"))
	if len(url) < 1 {
		fmt.Fprintf(ctx, "missing param url")
		return
	}
	hex := colorpicker.GetImageColor(url)
	response := colorResponse{
		CoverColor: hex,
	}
	JSONResponse(ctx, response)
}
