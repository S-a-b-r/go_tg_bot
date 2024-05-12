package main

import (
	"flag"
	"log"

	tgClient "telegram-bot/clients/telegram"
	event_consumer "telegram-bot/consumer/event-consumer"
	"telegram-bot/events/telegram"
	"telegram-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files-storage"
	batchSize   = 100
)

//7051466538:AAHEwSjegsleupL80Nbzcm9jiY4fm594WXw

func main() {
	tgClient := tgClient.New(tgBotHost, mustToken())

	eventsProcessor := telegram.New(tgClient, files.New(storagePath))

	log.Print("Service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"7051466538:AAHEwSjegsleupL80Nbzcm9jiY4fm594WXw",
		"token for access telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is empty")
	}

	return *token
}
