package main

import (
	"os"

	"github.com/matmazurk/file-store/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 3 {
		panic("invalid number or args")
	}
	conn, err := grpc.Dial(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	c := client.New(conn)
	err = c.SendFile(os.Args[2])
	if err != nil {
		panic(err)
	}
}
