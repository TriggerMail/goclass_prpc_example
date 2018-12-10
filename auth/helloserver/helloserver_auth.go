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

// matches prpc.constants.AUTHORIZATION_PREFIX
const authorizationPrefix = "Bearer "

type tokenAuthenticator struct {
	token string
}

func (t *tokenAuthenticator) Authenticate(ctx context.Context, req *http.Request) (
	context.Context, error) {
	fmt.Println("checking token:", req.Header.Get("Authorization"))
	authorized := req.Header.Get("Authorization") == authorizationPrefix+t.token
	if !authorized {
		return nil, fmt.Errorf("Unauthenticated")
	}
	return ctx, nil
}

func main() {

	hs := &helloServer{}
	prpcServer := &prpc.Server{
		Authenticator: &tokenAuthenticator{"token"},
	}
	helloworld.RegisterGreeterServer(prpcServer, hs)
	r := router.New()
	middleware := router.NewMiddlewareChain()
	prpcServer.InstallHandlers(r, middleware)

	fmt.Println("serving prpc on http://localhost:8080/")
	err := http.ListenAndServe("127.0.0.1:8080", r)
	if err != nil {
		panic(err)
	}
}
