package grpcmock_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/soroushj/grpcmock"
	"google.golang.org/grpc"
)

type (
	key string
)

const (
	servicePrefix = "/package.Service/"
	methodA       = "MethodA"
	methodB       = "MethodB"
	ctxKey        = key("ctx-key")
	ctxVal        = "ctx-val"
	reqVal        = "req"
	realVal       = "real"
	mockVal       = "mock"
)

var (
	infoA    = &grpc.UnaryServerInfo{FullMethod: servicePrefix + methodA}
	infoB    = &grpc.UnaryServerInfo{FullMethod: servicePrefix + methodB}
	respReal = fmt.Sprintf("%v/%v/%v", realVal, ctxVal, reqVal)
	respMock = fmt.Sprintf("%v/%v/%v", mockVal, ctxVal, reqVal)
	errReal  = errors.New("err-real")
	errMock  = errors.New("err-mock")
)

func handlerReal(ctx context.Context, req any) (any, error) {
	return fmt.Sprintf("%v/%v/%v", realVal, ctx.Value(ctxKey), req), errReal
}

func handlerMock(ctx context.Context, req any) (any, error) {
	return fmt.Sprintf("%v/%v/%v", mockVal, ctx.Value(ctxKey), req), errMock
}

func TestNew(t *testing.T) {
	mock := grpcmock.New()
	interceptor := mock.UnaryServerInterceptor()
	ctx := context.WithValue(context.Background(), ctxKey, ctxVal)
	resp, err := interceptor(ctx, reqVal, infoA, handlerReal)
	if resp != respReal {
		t.Errorf("resp: got %q want %q", resp, respReal)
	}
	if err != errReal {
		t.Errorf("err: got %v want %v", err, errReal)
	}
}

func TestSetResponse(t *testing.T) {
	mock := grpcmock.New()
	interceptor := mock.UnaryServerInterceptor()
	mock.SetResponse(methodA, &grpcmock.UnaryResponse{Resp: mockVal, Err: errMock})
	ctx := context.WithValue(context.Background(), ctxKey, ctxVal)
	resp, err := interceptor(ctx, nil, infoA, handlerReal)
	if resp != mockVal {
		t.Errorf("resp: got %q want %q", resp, mockVal)
	}
	if err != errMock {
		t.Errorf("err: got %v want %v", err, errMock)
	}
}

func TestSetHandler(t *testing.T) {
	mock := grpcmock.New()
	interceptor := mock.UnaryServerInterceptor()
	mock.SetHandler(methodA, handlerMock)
	ctx := context.WithValue(context.Background(), ctxKey, ctxVal)
	resp, err := interceptor(ctx, reqVal, infoA, handlerReal)
	if resp != respMock {
		t.Errorf("resp: got %q want %q", resp, respMock)
	}
	if err != errMock {
		t.Errorf("err: got %v want %v", err, errMock)
	}
}

func TestUnsetResponse(t *testing.T) {
	mock := grpcmock.New()
	interceptor := mock.UnaryServerInterceptor()
	mock.SetResponse(methodA, &grpcmock.UnaryResponse{Resp: mockVal, Err: errMock})
	mock.Unset(methodA)
	ctx := context.WithValue(context.Background(), ctxKey, ctxVal)
	resp, err := interceptor(ctx, reqVal, infoA, handlerReal)
	if resp != respReal {
		t.Errorf("resp: got %q want %q", resp, respReal)
	}
	if err != errReal {
		t.Errorf("err: got %v want %v", err, errReal)
	}
}

func TestUnsetHandler(t *testing.T) {
	mock := grpcmock.New()
	interceptor := mock.UnaryServerInterceptor()
	mock.SetHandler(methodA, handlerMock)
	mock.Unset(methodA)
	ctx := context.WithValue(context.Background(), ctxKey, ctxVal)
	resp, err := interceptor(ctx, reqVal, infoA, handlerReal)
	if resp != respReal {
		t.Errorf("resp: got %q want %q", resp, respReal)
	}
	if err != errReal {
		t.Errorf("err: got %v want %v", err, errReal)
	}
}

func TestClear(t *testing.T) {
	mock := grpcmock.New()
	interceptor := mock.UnaryServerInterceptor()
	mock.SetResponse(methodA, &grpcmock.UnaryResponse{Resp: mockVal, Err: errMock})
	mock.SetHandler(methodB, handlerMock)
	mock.Clear()
	t.Run("MethodA", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxKey, ctxVal)
		resp, err := interceptor(ctx, reqVal, infoA, handlerReal)
		if resp != respReal {
			t.Errorf("resp: got %q want %q", resp, respReal)
		}
		if err != errReal {
			t.Errorf("err: got %v want %v", err, errReal)
		}
	})
	t.Run("MethodB", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxKey, ctxVal)
		resp, err := interceptor(ctx, reqVal, infoB, handlerReal)
		if resp != respReal {
			t.Errorf("resp: got %q want %q", resp, respReal)
		}
		if err != errReal {
			t.Errorf("err: got %v want %v", err, errReal)
		}
	})
}
