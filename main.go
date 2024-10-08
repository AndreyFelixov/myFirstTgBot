package main

import (
	"flag"
	"log"
	tgClient "myFirstTgBot/clients/telegram"
	eventconsumer "myFirstTgBot/consumer/event-consumer"
	"myFirstTgBot/events/telegram"
	"myFirstTgBot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), files.New(storagePath))
	log.Print("service started")
	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is empty")
	}

	return *token
}
