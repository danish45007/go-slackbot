package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Krognol/go-wolfram"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go"
)

var wolframClient *wolfram.Client

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
	slackBot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	witClient := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_TOKEN")}
	go printCommandEvents(slackBot.CommandEvents())
	slackBot.Command("query for bot - <message>", &slacker.CommandDefinition{
		Description: "send any question to wolfram",
		Example:     "what is the size of earth ?",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")
			witResponse, err := witClient.Parse(&witai.MessageRequest{
				Query: query,
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			data, _ := json.MarshalIndent(witResponse, "", "    ")
			stringifyData := string(data[:])
			targetData := gjson.Get(stringifyData, "entities.wolfram_search_query.0.value")
			stringTargetData := targetData.String()
			wolframResponse, err := wolframClient.GetSpokentAnswerQuery(stringTargetData, wolfram.Metric, 1000)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			response.Reply(wolframResponse)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := slackBot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
