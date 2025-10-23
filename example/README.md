# Notes CLI: An example for grpcmock

This example demonstrates how grpcmock can be used for writing integration tests.

In this example, we define a *Notes* gRPC service for retrieving notes and implement a CLI client for this service.
We will use the grpcmock package to write integration tests for the CLI.

## Example structure

- [`notes/`](./notes/) defines a *Notes* gRPC service.
  - [`notes/notes.proto`](./notes/notes.proto) is the protobuf definition.
  - [`notes/gen.go`](./notes/gen.go) contains a `go:generate` comment for generating Go code from the protobuf definition.
  - The `notes/*.pb.go` files are automatically generated and define the Go types for the *Notes* service.
- [`notescli/`](./notescli/) implements a CLI client for the *Notes* server.
  - [`notescli/notescli.go`](./notescli/notescli.go) implements a `NotesCLI` type that takes a *Notes* client as a dependency.
  - [`notescli/notescli_test.go`](./notescli/notescli_test.go) is the **interesting part.**
    In this test, first we create a mock *Notes* server using grpcmock, then use the mock server to create a real *Notes* client, and finally use this client to create a `NotesCLI` to test.
  - [`notescli/cmd/notescli/main.go`](./notescli/cmd/notescli/main.go) uses `NotesCLI` to implement a CLI program.
    You don't need to read this file to understand the example.

## Go generate

As mentioned in the previous section, the `notes/*.pb.go` files are automatically generated.
You can run `go generate ./...` to re-generate them.
This requires having [`protoc`](https://protobuf.dev/installation/), [`protoc-gen-go`](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go), and [`protoc-gen-go-grpc`](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) installed:

```bash
# Install protoc
# on Ubuntu
sudo apt install -y protobuf-compiler
# on macOS
brew install protobuf
# on Windows
winget install protobuf

# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
