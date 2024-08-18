package errcode

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/madlabx/pkgx/errors"
	"github.com/madlabx/pkgx/httpx"
	"github.com/stretchr/testify/require"
)

func TestWrapItself(t *testing.T) {
	err := ErrInvalidNonce(ErrInvalidNonce())
	fmt.Printf("%+v\n", err.Unwrap())
	require.NotNil(t, err)
}

func TestIsItself(t *testing.T) {
	err := errors.Wrap(ErrObjectNotExist())

	ok := errors.Is(err, ErrNotFound())
	require.False(t, ok)
}

func TestStack(t *testing.T) {
	err := fmt.Errorf("testerror")
	ec := ErrHttpStatus(http.StatusTooManyRequests).WithError(err)

	fmt.Printf("ec.Unwrap()%%+v:%+v\n", ec.Unwrap())

	ec2 := ErrHttpStatus(http.StatusTooManyRequests).WithErrorf("err from WithErrorf, %v", err)

	fmt.Printf("%+v\n", ec2.(*httpx.JsonResponse).Unwrap())

	ec1 := ErrBadRequest().WithErrorf("err from WithErrorf")

	fmt.Printf("ec1.Unwrap()%%+v:%+v\n", ec1.(*httpx.JsonResponse).Unwrap())
	fmt.Printf("%s\n", ec1)

	wrapErr := errors.Wrapf(ec1, "newError from errors.Wrapf")

	fmt.Printf("%%#v:%#v\n", wrapErr)
	fmt.Printf("%+v\n", wrapErr)
	fmt.Printf("%s\n", wrapErr)
	fmt.Printf("wrapErr.Error():%s\n", wrapErr.Error())
}

func TestWithErrorf(t *testing.T) {
	e := ErrObjectExist().WithErrorf("dir existing")
	fmt.Printf("err:%v\n", e)

	e1 := ErrObjectExist()
	//should not impact e
	require.NotEqual(t, e, e1)
	fmt.Printf("err:%v\n", e1)

	e2 := ErrObjectExist().WithErrorf("dir existing")
	fmt.Printf("e2:%v\n", e2)
	require.Equal(t, e.Error(), e2.Error())
	require.NotEqual(t, e1.Error(), e2.Error())

	e3 := e2.(*httpx.JsonResponse).WithErrorf("dir existing")
	require.Equal(t, e3.Error(), e2.Error())
	fmt.Printf("e3:%v\n", e3)
	fmt.Printf("e2:%v\n", e2)

	e4 := ErrObjectExist().WithErrorf("dir existing1")
	require.NotEqual(t, e3.Error(), e4.Error())

	fmt.Printf("e4:%v\n", e4)
	fmt.Println(Len())
}

func TestRawWrap(t *testing.T) {
	e := ErrObjectExist(errors.New("dir existing"))
	fmt.Printf("err:%v\n", e)

	e1 := ErrObjectExist()
	//should not impact e
	require.NotEqual(t, e, e1)
	fmt.Printf("err:%v\n", e1)

	e2 := ErrObjectExist(errors.New("dir existing"))
	fmt.Printf("e2:%v\n", e2)
	require.Equal(t, e.Error(), e2.Error())
	require.NotEqual(t, e1.Error(), e2.Error())

	e3 := e2.JsonResponse.WithErrorf("dir existing")
	require.Equal(t, e3.Error(), e2.Error())
	fmt.Printf("e3:%v\n", e3)
	fmt.Printf("e2:%v\n", e2)

	e4 := ErrObjectExist(errors.New("dir existing1"))
	require.NotEqual(t, e3.Error(), e4.Error())

	fmt.Printf("e4:%v\n", e4)
	fmt.Println(Len())
}

func TestWithMessagef(t *testing.T) {
	e := ErrObjectExist(errors.New("dir existing"))
	e1 := ErrObjectExist()
	//should not impact e
	require.NotEqual(t, e, e1)

	e2 := ErrObjectExist(errors.New("dir existing"))
	require.Equal(t, e.Error(), e2.Error())
	require.NotEqual(t, e1.Error(), e2.Error())

	e3 := e2.WithMessagef("dir existing")
	require.Equal(t, e3.Error(), e2.Error())
	require.True(t, e2.Is(e3))
	require.Equal(t, e2.IsOK(), e3.(*httpx.JsonResponse).IsOK())

	e4 := ErrSuccess()
	msgString := "i am test"
	e5 := ErrSuccess().WithMessagef(msgString).(*httpx.JsonResponse)
	require.True(t, e4.Is(e5))
	require.True(t, e5.IsOK())

	require.Equal(t, msgString, e5.Message)
	require.NotEqual(t, msgString, e4.Error())

	require.Equal(t, e5.Error(), msgString)
	require.NotEqual(t, e5.JsonString(), e4.JsonString())

}
