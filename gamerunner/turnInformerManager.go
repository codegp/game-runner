package main

import (
	"log"
	"sync"

	gobj "github.com/codegp/game-runner/gameobjects"
	"github.com/codegp/game-runner/turninformer"
	"github.com/codegp/kube-client"
)

type turnInformerManager struct {
	clients            map[int32]*turninformer.TurnInformerClient
	botCreationChans   map[int32]chan string
	creationChansMutex *sync.Mutex
	kubeClient         *kubeclient.KubeClient
}

func newTurnInformerManager() *turnInformerManager {
	kclient, err := kubeclient.NewClient()
	if err != nil {
		log.Fatalf("failed to instantiate kubeclient, %v", err)
	}

	return &turnInformerManager{
		clients:            map[int32]*turninformer.TurnInformerClient{},
		botCreationChans:   map[int32]chan string{},
		creationChansMutex: &sync.Mutex{},
		kubeClient:         kclient,
	}
}

func (t *turnInformerManager) startTurn(botID int32) error {
	return t.clients[botID].StartTurn()
}

func (t *turnInformerManager) newBotCreationChan(botID int32) {
	t.creationChansMutex.Lock()
	t.botCreationChans[botID] = make(chan string)
	t.creationChansMutex.Unlock()
}

func (t *turnInformerManager) getBotCreationChan(botID int32) chan string {
	t.creationChansMutex.Lock()
	defer t.creationChansMutex.Unlock()
	return t.botCreationChans[botID]
}

func (t *turnInformerManager) reportIP(botID int32, ip string) {
	t.getBotCreationChan(botID) <- ip
}

func (t *turnInformerManager) createNewClients(botsToCreate []*gobj.Bot) {
	t.botCreationChans = map[int32]chan string{}
	for _, bot := range botsToCreate {
		t.newBotCreationChan(bot.ID)
		go t.createBot(bot)
	}

	for _, bot := range botsToCreate {
		ip := <-t.getBotCreationChan(bot.ID)
		client, err := newTurnInformerClient(ip)
		if err != nil {
			log.Fatalf("Failed to start turn informer client for a new bot, err:\n%v", err)
		}
		t.clients[bot.ID] = client
	}
}

func (t *turnInformerManager) destroyClients(botsToDestroy []int32) {
	for _, botID := range botsToDestroy {
		go t.destroyBot(botID)
		delete(t.clients, botID)
	}
}

func (t *turnInformerManager) destroyAllClients() {
	for botID, client := range t.clients {
		err := client.Destroy()
		if err != nil {
			log.Printf("Failed to destroy bot %d, err:\n%v", botID, err)
		}
	}
}

func (t *turnInformerManager) createBot(bot *gobj.Bot) {
	pod, err := t.kubeClient.StartBot(getIP(), bot.ID, projects[bot.TeamID], game)
	if err != nil {
		log.Fatalf("Failed to start bot pod, err:\n%v", err)
	}
	pod, err = t.kubeClient.WatchToStartup(pod)
	if err != nil {
		log.Fatalf("Failed to start bot pod, err:\n%v", err)
	}

	t.reportIP(bot.ID, pod.Status.PodIP)
}

func (t *turnInformerManager) destroyBot(botID int32) {
	err := t.clients[botID].Destroy()
	if err != nil {
		log.Printf("Failed to destroy bot %d, err:\n%v", botID, err)
	}
}
