package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TriggerMail/go_prpc_example/helloworld"
	"go.chromium.org/luci/grpc/prpc"
	"go.chromium.org/luci/server/router"
)

type helloServer struct {
	callCount int
}

func (h *helloServer) SayHello(ctx context.Context, request *helloworld.HelloRequest) (
	*helloworld.HelloReply, error) {
	h.callCount += 1
	resp := &helloworld.HelloReply{
		Message: fmt.Sprintf("hello %s; call number %d", request.Name, h.callCount),
	}
	fmt.Println("called with", request.Name, h.callCount)
	return resp, nil
}

func main() {

	// create server structs
	hs := &helloServer{}
	prpcServer := &prpc.Server{
		Authenticator: prpc.NoAuthentication,
	}

	// register server
	helloworld.RegisterGreeterServer(prpcServer, hs)

	// create new router
	r := router.New()
	middleware := router.NewMiddlewareChain()
	prpcServer.InstallHandlers(r, middleware)

	fmt.Println("serving prpc on http://:8080/")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
