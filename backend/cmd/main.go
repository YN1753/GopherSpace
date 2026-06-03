package main

import (
	"context"
	backend "gopherspace"
	"gopherspace/internal/core_ai"
	"gopherspace/internal/gateway"
	"log"

	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

func main() {
	messages := []*schema.Message{
		schema.SystemMessage("You are a helpful assistant."),
		schema.UserMessage("Hi, how are you?"),
	}
	ctx := context.Background()
	model, err := backend.Minimax(ctx)
	if err != nil {
		log.Fatal(err)
	}
	service := core_ai.NewCoreAiService(model)
	err = service.ChatWithModel(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	hub := gateway.NewHub(core_ai.NewCoreAiService(model))
	go hub.Run()
	r.GET("/v1/ws", gateway.HttpToWSEntry(hub))
	r.Run(":8181")

}
