package errcode

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/madlabx/pkgx/errcodex"
	"github.com/madlabx/pkgx/httpx"
	"github.com/madlabx/pkgx/log"
	"github.com/madlabx/pkgx/utils"
)

var _ errcodex.ErrorCodeIf = &ErrorCode{}

func init() {
	httpx.RegisterErrCodeDictionary(&ErrorCode{})
}

type ErrorCode struct {
	httpx.JsonResponse
}

func (ec *ErrorCode) ToStringWithStack() {

}

func (ec *ErrorCode) GetHttpStatus() int {
	return ec.Status
}

func (ec *ErrorCode) GetCode() string {
	return ec.Code
}

func (ec *ErrorCode) GetErrno() int {
	return ec.Errno
}

func (ec *ErrorCode) GetBadRequest() errcodex.ErrorCodeIf {
	return ErrBadRequest()
}

func (ec *ErrorCode) GetInternalError() errcodex.ErrorCodeIf {
	return ErrInternalServerError()
}

func (ec *ErrorCode) GetSuccess() errcodex.ErrorCodeIf {
	return OK
}

func (ec *ErrorCode) NewRequestId() string {
	return uuid.New().String()
}

func (ec *ErrorCode) ToCode(errno int) string {
	return trimHttpStatusText(errno)
}

func (ec *ErrorCode) ToHttpStatus(errno int) int {
	httpStatusText := http.StatusText(errno)
	if httpStatusText == "" {
		return 0
	} else {
		return errno
	}
}

func trimHttpStatusText(status int) string {
	trimmedSpace := strings.Replace(http.StatusText(status), " ", "", -1)
	trimmedSpace = strings.Replace(trimmedSpace, "-", "", -1)
	return trimmedSpace
}

// HttpCode required by interface httpx.JsonResponseError
func (ec *ErrorCode) HttpCode() int {
	//TODO support customized code
	return ec.Status
}

//func (ec *ErrorCode) Is(target error) bool {
//	te, ok := target.(*ErrorCode)
//	return ok && ec.Code == te.Code
//}

// HttpCode required by interface httpx.JsonResponseError
//
//	func (ec *ErrorCode) ToJsonResponseWithStack(depth int) *httpx.JsonResponse {
//		//TODO support customized code
//		return errors.WrapWithRelativeStackDepth()
//	}

var (
	once        sync.Once
	errCodeDict map[string]*ErrorCode
)

func Len() int {
	return len(errCodeDict)
}

func DumpErrorCodes() string {
	output := fmt.Sprintf("Total:%d\n", len(errCodeDict))
	return utils.ToPrettyString(errCodeDict) + output
}

func funcNewErrCode(errno int, opts ...string) func(errs ...error) *ErrorCode {
	once.Do(func() {
		errCodeDict = make(map[string]*ErrorCode)
	})

	errCode := newErrCode(errno, opts...)
	errCodeDict[errCode.Code] = errCode

	return func(errs ...error) *ErrorCode {
		//return clone
		ecopy := &ErrorCode{JsonResponse: *errCode.Copy()}
		if errs == nil {
			_ = ecopy.WithStack(2)
		} else {
			if len(errs) > 1 {
				log.Panicf("cannot accept more than one error")
			}
			//TODO only accept 1 error
			_ = ecopy.WithError(errs[0], 2)
		}
		return ecopy
	}
}

func newErrCode(errno int, opts ...string) *ErrorCode {
	httpStatusText := http.StatusText(errno)

	var httpStatus int
	if errno < constClientErrorBaseIndex && httpStatusText != "" {
		httpStatus = errno
	} else if errno >= constClientErrorBaseIndex && errno < constInternalErrorBaseIndex {
		httpStatus = http.StatusBadRequest
	} else if errno >= constInternalErrorBaseIndex && errno < constInvalidErrorBaseIndex {
		httpStatus = http.StatusInternalServerError
	} else {
		return ErrInvalidErrorCode
	}

	var code string
	if len(opts) > 0 {
		code = opts[0]
	} else {
		trimmedSpace := strings.Replace(httpStatusText, " ", "", -1)
		code = strings.Replace(trimmedSpace, "-", "", -1)
	}

	err := &ErrorCode{httpx.JsonResponse{
		Status: httpStatus,
		Errno:  errno,
		Code:   code,
	}}

	return err
}

const (
	constClientErrorBaseIndex     = 4000
	constInternalErrorBaseIndex   = 5000
	constDependencyErrorBaseIndex = 6000
	constInvalidErrorBaseIndex    = 9999
)

var (
	OK = newErrCode(http.StatusOK)

	ErrInvalidErrorCode = &ErrorCode{httpx.JsonResponse{
		Status: 0,
		Errno:  constInvalidErrorBaseIndex,
		Code:   "FatalErrorInvalidErrorCode",
	}}

	// Code == http status
	ErrSuccess    = funcNewErrCode(http.StatusOK)
	ErrBatchError = funcNewErrCode(http.StatusOK, "InternalServerError")
	ErrPartlyOk   = funcNewErrCode(http.StatusOK, "PartlyOk")
	ErrErrorFound = funcNewErrCode(http.StatusOK, "ErrorFound")

	ErrInsufficientStorage = funcNewErrCode(http.StatusInsufficientStorage)
	ErrInternalServerError = funcNewErrCode(http.StatusInternalServerError)
	ErrTimeout             = funcNewErrCode(http.StatusGatewayTimeout)
	ErrForbidden           = funcNewErrCode(http.StatusForbidden)
	ErrNotFound            = funcNewErrCode(http.StatusNotFound)
	ErrConflict            = funcNewErrCode(http.StatusConflict)
	ErrUnauthorized        = funcNewErrCode(http.StatusUnauthorized)
	ErrPreconditionFailed  = funcNewErrCode(http.StatusPreconditionFailed)
	ErrTooManyRequests     = funcNewErrCode(http.StatusTooManyRequests)
	ErrNotImplemented      = funcNewErrCode(http.StatusNotImplemented)

	ErrFailedDependency    = funcNewErrCode(http.StatusFailedDependency, "No storage")
	ErrMaxBindingLimited   = funcNewErrCode(http.StatusBadRequest, "MaxBindingLimited")
	ErrBadRequest          = funcNewErrCode(http.StatusBadRequest)
	ErrUnsupportedDevice   = funcNewErrCode(http.StatusBadRequest, "UnsupportedDevice")
	ErrMissPublicParameter = funcNewErrCode(http.StatusBadRequest, "MissPublicParameter")
	ErrInvalidSumAlgo      = funcNewErrCode(http.StatusBadRequest, "InvalidSumAlgo")
	ErrIsDirectory         = funcNewErrCode(http.StatusBadRequest, "TargetIsDirectory")
	ErrSourceIsParent      = funcNewErrCode(http.StatusBadRequest, "SourceIsParent")
	ErrExpiredRequest      = funcNewErrCode(http.StatusBadRequest, "ExpiredRequest")
	ErrObjectExist         = funcNewErrCode(http.StatusBadRequest, "ObjectExist")
	ErrObjectNotExist      = funcNewErrCode(http.StatusBadRequest, "ObjectNotExist")
	ErrSessionNotExist     = funcNewErrCode(http.StatusBadRequest, "SessionNotExist")
	ErrNoNeedUpgrade       = funcNewErrCode(http.StatusBadRequest, "NoNeedUpgrade")
	ErrUserExist           = funcNewErrCode(http.StatusBadRequest, "UserExist")
	ErrWrongSign           = funcNewErrCode(http.StatusBadRequest, "WrongSign")
	ErrInvalidVrfCode      = funcNewErrCode(http.StatusBadRequest, "InvalidSmsCode")
	ErrExpiredVrfCode      = funcNewErrCode(http.StatusBadRequest, "ExpiredSmsCode")
	ErrInvalidDeviceId     = funcNewErrCode(http.StatusBadRequest, "InvalidDeviceId")
	ErrAlreadyFormated     = funcNewErrCode(http.StatusBadRequest, "AlreadyFormated")

	ErrDeviceAlreadyBound = funcNewErrCode(http.StatusBadRequest, "DeviceAlreadyBound")
	ErrDeviceOffline      = funcNewErrCode(http.StatusBadRequest, "DeviceOffline")

	ErrInvalidJwt      = funcNewErrCode(http.StatusUnauthorized, "InvalidJwt")
	ErrInvalidPassword = funcNewErrCode(http.StatusUnauthorized, "InvalidPassword")
	ErrInvalidNonce    = funcNewErrCode(http.StatusUnauthorized, "InvalidNonce")
	ErrExpiredNonce    = funcNewErrCode(http.StatusUnauthorized, "ExpiredNonce")
	ErrInvalidSign     = funcNewErrCode(http.StatusUnauthorized, "InvalidSign")
	ErrExpiredToken    = funcNewErrCode(http.StatusUnauthorized, "ExpiredToken")
	ErrInvalidToken    = funcNewErrCode(http.StatusUnauthorized, "InvalidToken")
	ErrInvalidIssuer   = funcNewErrCode(http.StatusUnauthorized, "InvalidIssuer")

	//4XXX 调用方错误
	//5XXX 服务器内部错误
	ErrConflictConnection = funcNewErrCode(http.StatusInternalServerError, "ConflictConnection")
	ErrInvalidDataType    = funcNewErrCode(http.StatusInternalServerError, "InvalidDataType")
	ErrEmptyKey           = funcNewErrCode(http.StatusInternalServerError, "EmptyKey")
	ErrEmptyPassword      = funcNewErrCode(http.StatusInternalServerError, "EmptyPassword")
	ErrRootUserDeletion   = funcNewErrCode(http.StatusInternalServerError, "RootUserDeletion")

	//6XXX 依赖方错误

)

func ErrHttpStatus(status int) *ErrorCode {
	return funcNewErrCode(status)()
}

func ErrStrResp(status int, b any, format string, a ...any) error {
	err := funcNewErrCode(status)()
	return err.WithErrorf(format, a...)
}
