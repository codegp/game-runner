package main

import (
	"fmt"
	"log"

	"github.com/codegp/cloud-persister/models"
	gobj "github.com/codegp/game-runner/gameobjects"
	"github.com/codegp/game-runner/turninformer"
)

type turnInformerManager struct {
	clients map[int64]*turninformer.TurnInformerClient
}

func newTurnInformerManager(projects []*models.Project) *turnInformerManager {
	return &turnInformerManager{
		clients: createNewClients(projects), // team number to client
	}
}

func createNewClients(projects []*models.Project) map[int64]*turninformer.TurnInformerClient {
	clients := map[int64]*turninformer.TurnInformerClient{}
	var err error
	for _, project := range projects {
		clients[project.ID], err = newTurnInformerClient(fmt.Sprint("team-runner-", project.ID)) // todo: figure out how to do in kube
		if err != nil {
			panic(err)
		}
	}
	return clients
}

func (t *turnInformerManager) destroy() {
	for teamID, client := range t.clients {
		err := client.Destroy()
		if err != nil {
			log.Printf("Failed to destroy teamID %d, err:\n%v", teamID, err)
		}
	}
}

func (t *turnInformerManager) createBots(botsToCreate []*gobj.Bot) {
	for _, bot := range botsToCreate {
		t.clients[bot.TeamID].CreateBot(bot.ID)
	}
}

func (t *turnInformerManager) destroyBots(botsToDestroy []*gobj.Bot) {
	for _, bot := range botsToDestroy {
		t.clients[bot.TeamID].DestroyBot(bot.ID)
	}
}

func (t *turnInformerManager) startTurn(bot *gobj.Bot) error {
	return t.clients[bot.TeamID].StartTurn(bot.ID)
}
