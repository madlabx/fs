package controller

import (
	"fmt"
	"runtime/debug"

	"github.com/labstack/echo"
	"github.com/madlabx/pkgx/errors"
	"github.com/madlabx/pkgx/httpx"
	"github.com/madlabx/pkgx/log"
)

type HandleFuncCtx func(ctx echo.Context) error

func HandleCtx(fn HandleFuncCtx) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Recovered: %v", r)
				log.Errorf("receive callstack: %v", string(debug.Stack()))
				_ = httpx.SendResp(ctx, errors.New(fmt.Sprintf("%v", r)))
			}
		}()

		resp := fn(ctx)

		if resp != nil {
			jr := &httpx.JsonResponse{}
			ok := errors.As(resp, &jr)
			if !ok {
				log.Errorf("http handle error:%+v", resp)
			} else if jr.Unwrap() != nil {
				log.Errorf("http handle error:%+v", jr.Unwrap())
			}
		}

		return httpx.SendResp(ctx, resp)
	}
}
