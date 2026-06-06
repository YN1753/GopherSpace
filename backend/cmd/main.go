package main

import (
	"context"
	backend "gopherspace"
	"gopherspace/internal/core_ai"
	"gopherspace/internal/gateway"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	model, err := backend.Minimax(ctx)
	if err != nil {
		log.Fatal(err)
	}
	router := core_ai.NewRouterAgent(model)
	reviewer := core_ai.NewReviewerAgent(model)
	consultant := core_ai.NewConsultantAgent(model)
	tutor := core_ai.NewTutorAgent(model)
	graph, err := core_ai.NewChatGraph(ctx, router, reviewer, consultant, tutor, model)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	hub := gateway.NewHub(graph)
	go hub.Run()
	r.GET("/v1/ws", gateway.HttpToWSEntry(hub))
	r.Run(":8181")

}
