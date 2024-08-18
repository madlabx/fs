package api

import (
	"context"

	"github.com/madlabx/pkgx/errors"
	"github.com/madlabx/pkgx/httpx"
	"github.com/sirupsen/logrus"
)

type Gateway struct {
	*httpx.ApiGateway
	ctx context.Context
}

func New(pCtx context.Context,
	conf *httpx.LogConfig,
	logFormat logrus.Formatter) (*Gateway, error) {
	agw, err := httpx.NewApiGateway(pCtx, conf, logFormat)
	if err != nil {
		return nil, errors.Wrapf(err, "failure in NewApiGateway")
	}

	gw := &Gateway{
		ctx:        context.WithoutCancel(pCtx),
		ApiGateway: agw,
	}

	registerRouter(gw)

	return gw, nil
}
