package main

import "strconv"

// botStartInformerHandler implements the botStartInformer interface defined by thrift
type botStartInformerHandler struct {
	manager *turnInformerManager
}

// newFromBotHandler creates a new instance of the botStartInformerHandler
func newBotStartInformerHandler(manager *turnInformerManager) *botStartInformerHandler {
	return &botStartInformerHandler{
		manager: manager,
	}
}

// Started sends the new bots ip to the turnInformerManager so it can
// create a turnInformerClient
func (bs *botStartInformerHandler) Started(id, ip string) error {
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return err
	}
	bs.manager.reportIP(int32(i), ip)
	return nil
}
