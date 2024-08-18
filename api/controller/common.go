package controller

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/madlabx/pkgx/httpx"

	"github.com/madlabx/fs/common/cfg"
	"github.com/madlabx/fs/pkg/buildcontext"
)

func OnError(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = httpx.SendResp(c, httpx.StatusResp(code))
		} else {
			//index.ServeHTTP(c.Response().Writer, c.Request())
			//err = nil
			err = httpx.SendResp(c, httpx.Wrap(err))
		}
	}
}

func OnConfig(c echo.Context) error {
	return httpx.SuccessResp(cfg.Get())
}

func OnHealth(c echo.Context) error {
	return httpx.SuccessResp(buildcontext.Get())
}
