package transfer

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/madlabx/pkgx/errors"
	"github.com/madlabx/pkgx/httpx"
	"github.com/madlabx/pkgx/utils"
	"github.com/marusama/semaphore/v2"

	"github.com/madlabx/fs/common/cfg"
	"github.com/madlabx/fs/common/errcode"
)

// saveUploadedFile uploads the form file to specific dst.
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return errors.Wrap(err)
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return errors.Wrap(err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err)
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return errors.Wrap(err)
}

var sem semaphore.Semaphore
var once sync.Once

func generateFilePath(name string, expireDays int) string {
	return filepath.Join(strconv.Itoa(expireDays), utils.Md5Sum(name+strconv.FormatInt(time.Now().Unix(), 10)), filepath.Base(name))
}

const (
	ConstDownloadPath = "/v1/fs/files"
)

func generateFullUrl(path string) string {
	return cfg.Get().Sys.Domain + ConstDownloadPath + "/" + path
}

func Upload(file *multipart.FileHeader, expireDays int) (string, error) {

	once.Do(func() {
		sem = semaphore.New(cfg.Get().Sys.MaxUploadParallel)
	})

	if !sem.TryAcquire(1) {
		return "", httpx.StatusResp(http.StatusTooManyRequests)
	}
	defer sem.Release(1)

	notEnough, err := IsDiskSpaceNotEnough(uint64(file.Size))
	if err != nil {
		return "", httpx.Wrap(err)
	}
	if notEnough {
		return "", errcode.ErrInsufficientStorage().WithErrorf("No enough space")
	}
	dstFilePath := generateFilePath(file.Filename, expireDays)
	if err = saveUploadedFile(file, filepath.Join(cfg.Get().Sys.Root, dstFilePath)); err != nil {
		return "", errors.Wrap(err)
	}

	return generateFullUrl(dstFilePath), nil
}

func setContentDisposition(w http.ResponseWriter, r *http.Request, fileName string) {
	if r.URL.Query().Get("inline") == "true" {
		w.Header().Set("Content-Disposition", "inline")
	} else {
		// As per RFC6266 section 4.3
		w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+url.PathEscape(fileName))
	}
}

func Download(ctx echo.Context, path string) error {

	fullPath := filepath.Join(cfg.Get().Sys.Root, path)
	fs, err := os.Stat(fullPath)
	if errors.Is(err, os.ErrNotExist) {
		return errcode.ErrObjectNotExist(err)
	}
	if err != nil {
		return errors.Wrap(err)
	}

	if fs.IsDir() {
		return errcode.ErrBadRequest().WithErrorf("Dir not allowed, %v", fullPath)
	}

	fd, err := os.Open(fullPath)
	if err != nil {
		return errors.Wrap(err)
	}

	setContentDisposition(ctx.Response(), ctx.Request(), fs.Name())
	httpx.ServeContent(ctx.Response(), ctx.Request(), fs.Name(), fs.ModTime(), fs.Size(), fd)

	return nil
}
