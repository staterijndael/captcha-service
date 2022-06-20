package endpoint

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

func SendSuccessResponse(ctx *fasthttp.RequestCtx, response interface{}) {
	marshalledMessage, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	ctx.SetBody(marshalledMessage)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func SendErrorResponse(ctx *fasthttp.RequestCtx, errorCode int, error []byte) {
	ctx.SetBody(error)
	ctx.SetStatusCode(errorCode)
}
