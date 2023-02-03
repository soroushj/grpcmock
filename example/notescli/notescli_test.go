package notescli_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/soroushj/grpcmock"
	"github.com/soroushj/grpcmock/example/notes"
	"github.com/soroushj/grpcmock/example/notescli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	port = "5050"
)

var (
	mock = grpcmock.New()
)

func TestMain(m *testing.M) {
	// Create a gRPC server using the mock interceptor
	server := grpc.NewServer(grpc.UnaryInterceptor(mock.UnaryServerInterceptor()))
	// Register the generated Unimplemented implementation of the Notes server
	notes.RegisterNotesServer(server, &notes.UnimplementedNotesServer{})
	// Run the server
	addr := fmt.Sprintf(":%v", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen at %v: %v", addr, err)
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve at %v: %v", addr, err)
		}
	}()
	// Run the tests
	code := m.Run()
	// Stop the server and exit
	server.Stop()
	os.Exit(code)
}

func TestGetNoteText(t *testing.T) {
	// Create a Notes client
	addr := fmt.Sprintf("localhost:%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial %v: %v", addr, err)
	}
	defer conn.Close()
	client := notes.NewNotesClient(conn)
	// Create a NotesCLI to be tested
	nc := notescli.New(client)
	// Set a mock handler for the GetNote method on the Notes server
	mock.SetHandler("GetNote", func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(*notes.GetNoteRequest)
		if r.Id == "1" {
			return &notes.GetNoteResponse{
				Note: &notes.Note{
					Id:   "1",
					Text: "test",
				},
			}, nil
		}
		return nil, status.Error(codes.NotFound, "note not found")
	})
	// Run the subtests
	testCases := []struct {
		name    string
		id      string
		text    string
		errCode codes.Code
	}{
		{"Note exists", "1", "test", codes.OK},
		{"Note does not exist", "2", "", codes.NotFound},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			text, err := nc.GetNoteText(tc.id)
			if text != tc.text {
				t.Errorf("text: got %q want %q", text, tc.text)
			}
			if err, _ := status.FromError(err); err.Code() != tc.errCode {
				t.Errorf("err: got %v want %v", err.Code(), tc.errCode)
			}
		})
	}
	// Set a mock response for another test case
	mock.SetResponse("GetNote", &grpcmock.UnaryResponse{
		Resp: &notes.GetNoteResponse{},
	})
	t.Run("Bad response", func(t *testing.T) {
		text, err := nc.GetNoteText("3")
		if text != "" {
			t.Errorf("text: got %q want %q", text, "")
		}
		if err != notescli.ErrBadResponse {
			t.Errorf("err: got %v want %v", err, notescli.ErrBadResponse)
		}
	})
}
