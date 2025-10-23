package notescli_test

import (
	"context"
	"net"
	"testing"

	"github.com/soroushj/grpcmock"
	"github.com/soroushj/grpcmock/example/notes"
	"github.com/soroushj/grpcmock/example/notescli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestNotesCLI(t *testing.T) {
	// Create a gRPC server using a grpcmock interceptor
	mock := grpcmock.New()
	server := grpc.NewServer(grpc.UnaryInterceptor(mock.UnaryServerInterceptor()))
	defer server.Stop()
	// Register the generated Unimplemented implementation of the Notes server
	notes.RegisterNotesServer(server, &notes.UnimplementedNotesServer{})
	// Run the server on an available port
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			t.Errorf("serve at %v: %v", lis.Addr(), err)
		}
	}()
	// Connect to the mock Notes server
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("connect to %v: %v", lis.Addr(), err)
	}
	defer conn.Close()
	// Create a Notes client
	client := notes.NewNotesClient(conn)
	// Create a NotesCLI to be tested
	nc := notescli.New(client)
	ctx := context.Background()
	t.Run("GetNoteText", func(t *testing.T) {
		// Set a mock handler for the GetNote method on the Notes server
		mock.SetHandler("GetNote", func(ctx context.Context, req any) (any, error) {
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
		// Run the subtests
		testCases := []struct {
			name    string
			id      string
			text    string
			errCode codes.Code
		}{
			{"Found", "1", "a", codes.OK},
			{"Not found", "2", "", codes.NotFound},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				text, err := nc.GetNoteText(ctx, tc.id)
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
			text, err := nc.GetNoteText(ctx, "3")
			if text != "" {
				t.Errorf("text: got %q want %q", text, "")
			}
			if err != notescli.ErrBadResponse {
				t.Errorf("err: got %v want %v", err, notescli.ErrBadResponse)
			}
		})
	})
}
