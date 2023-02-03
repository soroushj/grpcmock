package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/soroushj/grpcmock/example/notes"
	"github.com/soroushj/grpcmock/example/notescli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: notescli server_addr note_id")
	}
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}
	addr := flag.Arg(0)
	id := flag.Arg(1)
	err := printNoteText(addr, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printNoteText(addr, id string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := notes.NewNotesClient(conn)
	nc := notescli.New(client)
	text, err := nc.GetNoteText(id)
	if err != nil {
		return err
	}
	fmt.Println(text)
	return nil
}
