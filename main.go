package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	godotenv.Load(".env")
	slackBot := slacker.NewClient("xoxb-2323039060352-2305584906260-ooJ6dODDUm5oGLWMhrEJB8cW",
		"xapp-1-A029608QS2V-2323047138656-b26fc207c5e54902f6188014f78e9351eb88f4bd6b05d696845d31f64c1832f0")
	go printCommandEvents(slackBot.CommandEvents())
	slackBot.Command("query for bot -,<message>", &slacker.CommandDefinition{
		Description: "send any question to wolfram",
		Example:     "what is the size of earth ?",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")
			fmt.Println(query)
			response.Reply("received")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := slackBot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
