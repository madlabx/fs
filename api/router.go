package api

import (
	"github.com/labstack/echo"
	"github.com/madlabx/pkgx/log"

	"github.com/madlabx/fs/api/controller"
	"github.com/madlabx/fs/module/transfer"
)

func registerRouter(gw *Gateway) {
	handleWrap := func(fn controller.HandleFuncCtx) echo.HandlerFunc {
		return controller.HandleCtx(fn)
	}

	gw.HTTPErrorHandler = controller.OnError

	fsUri := "/v1/fs"
	gw.GET(fsUri+"/health", handleWrap(controller.OnHealth))
	gw.GET(fsUri+"/config", handleWrap(controller.OnConfig))

	gw.POST(fsUri+"/files", handleWrap(controller.OnUpload))
	gw.GET(transfer.ConstDownloadPath+"/:expire/:md5/:name", handleWrap(controller.OnDownload))

	log.Infof("Got routes:%v", gw.RoutesToString())
}
