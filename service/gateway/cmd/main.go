package main

import (
	"github.com/stack-labs/stack-rpc"
	"github.com/stack-labs/stack-rpc/util/log"

	"github.com/stack-labs/stack-rpc-plugins/service/gateway/api"
)

func main() {
	svc := stack.NewService()

	// gateway server
	apiServer := api.NewServer(svc)
	svc.Init(apiServer.Options()...)

	// run service
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
