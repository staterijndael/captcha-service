package endpoint

import (
	"bytes"
	"github.com/valyala/fasthttp"
)

func (h *Handler) Endpoint(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()
	switch {
	case string(path) == "/generate_session":
		h.GenerateSession(ctx)
	case bytes.HasPrefix(path, []byte("/img/")):
		fasthttp.FSHandler("endpoint/captcha_images", 1)(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
