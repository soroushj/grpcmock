package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/soroushj/grpcmock/example/notes"
	"github.com/soroushj/grpcmock/example/notescli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: notescli server_addr")
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	addr := flag.Arg(0)
	ctx := context.Background()
	err := run(ctx, addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := notes.NewNotesClient(conn)
	nc := notescli.New(client)
	r := bufio.NewReader(os.Stdin)
	fmt.Println("enter \\q to quit")
	for {
		fmt.Print("\nid=")
		id, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		id = strings.TrimSpace(id)
		if id == "\\q" {
			return nil
		}
		text, err := nc.GetNoteText(ctx, id)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(text)
		}
	}
}
