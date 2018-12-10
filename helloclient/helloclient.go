package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TriggerMail/goclass_prpc_example/helloworld"
	"github.com/TriggerMail/luci-go/grpc/prpc"
)

func main() {

	fmt.Println("connecting to http://localhost:8080")

	opts := prpc.DefaultOptions()
	opts.Insecure = true

	client := &prpc.Client{Host: "localhost:8080", Options: opts, C: &http.Client{}}

	hwClient := helloworld.NewGreeterPRPCClient(client)

	req := &helloworld.HelloRequest{}
	req.Name = "foo"
	resp, err := hwClient.SayHello(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println("response:", resp.Message)
}
