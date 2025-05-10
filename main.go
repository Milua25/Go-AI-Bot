package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
	"log"
	"os"
)

var wolframClient *wolfram.Client

func handleError(err error) {
	if err != nil {
		fmt.Sprintf("Something went wrong:%s", err.Error())
	}
}

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Event)
		fmt.Println(event.Parameters)
		fmt.Println(event.Command)
		fmt.Println()
	}
}
func main() {
	fmt.Println("Hello GO😍!!!")

	// load the environment variables
	err := godotenv.Load(".env")
	handleError(err)

	// create a new bot
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	// Wit client
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))

	// WolframClient
	wolframClient = &wolfram.Client{
		AppID: os.Getenv("WOLFRAM_APP_ID"),
	}
	// print command events
	go printCommandEvents(bot.CommandEvents())

	bot.Command("query for bot - <message>", &slacker.CommandDefinition{
		Description: "Send any question to wolfram",
		Examples:    []string{"who is the president of nigeria", "who owns telsa"},
		Handler: func(botContext slacker.BotContext, request slacker.Request, writer slacker.ResponseWriter) {
			query := request.Param("message")
			fmt.Println(query)
			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			data, _ := json.MarshalIndent(msg, "", "  ")
			rough := string(data[:])
			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			answer := value.String()
			res, err := wolframClient.GetSpokentAnswerQuery(answer, wolfram.Metric, 1000)
			if err != nil {
				fmt.Println("There is an error")
			}
			fmt.Println(value)
			writer.Reply(res)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
