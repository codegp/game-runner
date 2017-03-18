package main

import (
	"fmt"
	"log"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/codegp/game-runner/api"
	"github.com/codegp/game-runner/gamestate"
)

func startServer(gameStateUtils *gamestate.GameStateUtils) *thrift.TSimpleServer {
	server, err := getServer(gameStateUtils)
	if err != nil {
		log.Fatalf("Error starting thrift server for game api %v", err)
	}

	go func() {
		log.Printf("Starting API server")
		e := server.Serve()
		if e != nil {
			log.Fatalf("the errz!!! %v", e)
		}
	}()
	return server
}

func getServer(gameStateUtils *gamestate.GameStateUtils) (*thrift.TSimpleServer, error) {
	//transport is a thrift.TServerTransport
	transport, err := thrift.NewTServerSocket(fmt.Sprintf("%s:9000", getIP()))
	if err != nil {
		return nil, err
	}

	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTCompactProtocolFactory()

	apiHandler := newAPIHandler(gameStateUtils)
	apiProc := api.NewAPIProcessor(apiHandler)
	server := thrift.NewTSimpleServer4(apiProc, transport, transportFactory, protocolFactory)

	return server, nil
}
