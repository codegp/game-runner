package main

import (
	"fmt"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/codegp/game-runner/turninformer"
)

const retryLimit int = 20

func newTurnInformerClient(ip string) (*turninformer.TurnInformerClient, error) {
	transport, err := thrift.NewTSocket(fmt.Sprintf("%s:9000", ip))
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
		time.Sleep(time.Millisecond * time.Duration(100*retryCount))
		retryCount++
	}

	return nil, fmt.Errorf("Failed to start client at %s:9000, err:\n%v", ip, err)
}
