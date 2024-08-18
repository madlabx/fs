package controller

import (
	"time"

	"github.com/labstack/echo"
	"github.com/madlabx/pkgx/httpx"
	"github.com/madlabx/pkgx/typex"

	"github.com/madlabx/fs/common/errcode"
	"github.com/madlabx/fs/module/transfer"
)

type reqDownload struct {
	Path string `hx_place:"path" hx_must:"true"`
}

func OnDownload(ctx echo.Context) error {
	//var req reqDownload
	//
	//if err := httpx.BindAndValidate(ctx, &req); err != nil {
	//	return errcode.ErrBadRequest(err)
	//}

	expire := ctx.Param("expire")

	md5 := ctx.Param("md5")
	name := ctx.Param("name")
	return transfer.Download(ctx, expire+"/"+md5+"/"+name)
}

type reqUpload struct {
	ExpireDays int `hx_place:"query" hx_default:"1"` // -1, no expire
	//Path       string `hx_place:"query"`
}

func (req *reqUpload) Validate() error {
	if req.ExpireDays < 0 || req.ExpireDays > 7 {
		return errcode.ErrBadRequest().WithErrorf("ExpireDays should between -1 and 7, -1 expect for no expiration")
	}

	return nil
}

func OnUpload(ctx echo.Context) error {
	var req reqUpload
	if err := httpx.BindAndValidate(ctx, &req); err != nil {
		return errcode.ErrBadRequest(err)
	}

	inFile, err := ctx.FormFile("upload")
	if err != nil {
		return errcode.ErrBadRequest(err)
	}

	url, err := transfer.Upload(inFile, req.ExpireDays)
	if err != nil {
		return err
	}

	expireAt := "-"
	if req.ExpireDays != 0 {
		expireAt = time.Now().Add(time.Hour * 24 * time.Duration(req.ExpireDays)).Format("2006-01-02 15:04:05")
	}

	return httpx.SuccessResp(typex.JsonMap{
		"ExpireAt":   expireAt,
		"ExpireDays": req.ExpireDays,
		"Url":        url,
	})
}
