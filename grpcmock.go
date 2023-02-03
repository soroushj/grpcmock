// Package grpcmock provides functionality for mocking gRPC servers dynamically.
//
// Currently, only unary RPCs are supported.
//
// All functions in this package are safe for concurrent use by multiple goroutines.
package grpcmock

import (
	"context"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

// UnaryResponse represents values returned by a unary RPC.
type UnaryResponse struct {
	Resp interface{}
	Err  error
}

// rh represents a mock response/handler.
type rh struct {
	response *UnaryResponse
	handler  grpc.UnaryHandler
}

// GRPCMock provides a gRPC interceptor and a set of methods for mocking gRPC servers.
type GRPCMock struct {
	// methods stores mock response/handlers, where the map key is the short method name.
	methods struct {
		sync.RWMutex
		m map[string]*rh
	}
}

// New returns a new [GRPCMock].
func New() *GRPCMock {
	gm := &GRPCMock{}
	gm.methods.m = make(map[string]*rh)
	return gm
}

// SetResponse sets a mock response and removes any mock handler for method.
func (gm *GRPCMock) SetResponse(method string, response *UnaryResponse) {
	gm.methods.Lock()
	gm.methods.m[method] = &rh{response: response}
	gm.methods.Unlock()
}

// SetHandler sets a mock handler and removes any mock response for method.
func (gm *GRPCMock) SetHandler(method string, handler grpc.UnaryHandler) {
	gm.methods.Lock()
	gm.methods.m[method] = &rh{handler: handler}
	gm.methods.Unlock()
}

// Unset removes any mock response or handler for method.
func (gm *GRPCMock) Unset(method string) {
	gm.methods.Lock()
	delete(gm.methods.m, method)
	gm.methods.Unlock()
}

// Clear removes any mock response or handler for all methods.
func (gm *GRPCMock) Clear() {
	gm.methods.Lock()
	for method := range gm.methods.m {
		delete(gm.methods.m, method)
	}
	gm.methods.Unlock()
}

// UnaryServerInterceptor returns a gRPC unary server interceptor which handles methods
// using mock responses or handlers. If no mock response or handler is set for a method,
// the registered handler will be used.
func (gm *GRPCMock) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Extract the short method name
		n := strings.LastIndex(info.FullMethod, "/")
		method := info.FullMethod[n+1:]
		// Use a mock response/handler if available
		gm.methods.RLock()
		rh := gm.methods.m[method]
		gm.methods.RUnlock()
		if rh != nil {
			if rh.response != nil {
				return rh.response.Resp, rh.response.Err
			}
			if rh.handler != nil {
				return rh.handler(ctx, req)
			}
		}
		// Use the registered handler
		return handler(ctx, req)
	}
}
