# grpcmock: Mock gRPC servers dynamically in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/soroushj/grpcmock.svg)](https://pkg.go.dev/github.com/soroushj/grpcmock)
[![CI](https://github.com/soroushj/grpcmock/actions/workflows/ci.yml/badge.svg)](https://github.com/soroushj/grpcmock/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/soroushj/grpcmock)](https://goreportcard.com/report/github.com/soroushj/grpcmock)

See [docs and examples on pkg.go.dev](https://pkg.go.dev/github.com/soroushj/grpcmock).

See the [example](./example) directory for a complete example.

## Usage

```go
import (
	"context"
	"log"
	"net"

	"github.com/soroushj/grpcmock"
	"github.com/soroushj/grpcmock/example/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Example() {
	// Create a gRPC server using a grpcmock interceptor
	mock := grpcmock.New()
	server := grpc.NewServer(grpc.UnaryInterceptor(mock.UnaryServerInterceptor()))

	// Register an implementation of your server; the example Notes server in this case.
	// This typically should be the generated Unimplemented implementation.
	notes.RegisterNotesServer(server, &notes.UnimplementedNotesServer{})

	// Run the server on an available port
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("serve at %v: %v", lis.Addr(), err)
		}
	}()

	// At this point, if you call any method from the running server, you will get an Unimplemented error.
	// Let's change this behavior for the GetNote method of this server.

	// This is how you can set a mock response for a method; an error in this case.
	// After this, you will get a NotFound error from GetNote instead of an Unimplemented error.
	mock.SetResponse("GetNote", &grpcmock.UnaryResponse{
		Err: status.Error(codes.NotFound, "note not found"),
	})

	// Similarly, you can set a non-error response.
	// After this, you will get the response below from GetNote instead of an error.
	mock.SetResponse("GetNote", &grpcmock.UnaryResponse{
		Resp: &notes.GetNoteResponse{
			Note: &notes.Note{
				Id:   "1",
				Text: "a",
			},
		},
	})

	// If you need something more than a simple response, you can set a handler.
	// After this, any call to GetNote will be handled using the function below.
	mock.SetHandler("GetNote", func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(*notes.GetNoteRequest)
		if r.Id == "1" {
			return &notes.GetNoteResponse{
				Note: &notes.Note{
					Id:   "1",
					Text: "a",
				},
			}, nil
		}
		return nil, status.Error(codes.NotFound, "note not found")
	})

	// You can remove any previously-set response or handler for a method.
	// After this, GetNote will return an Unimplemented error.
	mock.Unset("GetNote")

	// You can also remove any response or handler for all methods
	mock.Clear()

	// Stop the server
	server.Stop()
}
```
