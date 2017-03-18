package main

import (
	"fmt"
	"log"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/codegp/game-runner/turninformer"
)

const retryLimit int = 120

func newTurnInformerClient(addr string) (*turninformer.TurnInformerClient, error) {
	transport, err := thrift.NewTSocket(fmt.Sprintf("%s:9000", addr))
	if err != nil {
		return nil, err
	}

	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTCompactProtocolFactory()
	t := transportFactory.GetTransport(transport)

	retryCount := 0
	for retryCount <= retryLimit {
		if err = t.Open(); err == nil {
			return turninformer.NewTurnInformerClientFactory(t, protocolFactory), nil
		}
		log.Printf("Could not connect to %s:9000. Retrying in 1s: %v. ", addr, err)
		time.Sleep(time.Millisecond * time.Duration(1000))
		retryCount++
	}

	return nil, fmt.Errorf("Failed to start client at %s:9000, err:\n%v", addr, err)
}
