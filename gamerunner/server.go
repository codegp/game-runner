package main

import (
	"fmt"
	"log"

	"github.com/codegp/game-runner/api"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/codegp/game-runner/gamestate"
)

func startServer(gameStateUtils *gamestate.GameStateUtils, tiManager *turnInformerManager) *thrift.TSimpleServer {
	server, err := getServer(gameStateUtils, tiManager)
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

func getServer(gameStateUtils *gamestate.GameStateUtils, tiManager *turnInformerManager) (*thrift.TSimpleServer, error) {
	//transport is a thrift.TServerTransport
	transport, err := thrift.NewTServerSocket(fmt.Sprintf("%s:9000", getIP()))
	if err != nil {
		return nil, err
	}

	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTCompactProtocolFactory()

	apiHandler := newAPIHandler(gameStateUtils)
	apiProc := api.NewAPIProcessor(apiHandler)

	// botStartInformerHandler := newBotStartInformerHandler(tiManager)
	// botStartInformerProc := botstartinformer.NewBotStartInformerProcessor(botStartInformerHandler)

	// mProc := thrift.NewTMultiplexedProcessor()
	// mProc.RegisterProcessor("API", apiProc)
	// mProc.RegisterProcessor("BotStartInformer", botStartInformerProc)
	//
	// server := thrift.NewTSimpleServer4(mProc, transport, transportFactory, protocolFactory)

	server := thrift.NewTSimpleServer4(apiProc, transport, transportFactory, protocolFactory)

	return server, nil
}
