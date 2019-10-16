package main

import (
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var configFilePath = flag.String("cfgFile", "./config.yml", "Path to deploy config")

func main() {
	flag.Parse()

	var dCfg deployConfig
	dCfg.loadConfig(*configFilePath)

	err := godotenv.Load()
	if err != nil {
		return
	}

	token := os.Getenv("slackkey")
	api := slack.New(
		token,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			if len(ev.Attachments) > 0 {
				text := ev.Attachments[0].Text
				log.Printf("Message: %v\n", text)

				matched := strings.Contains(text, "jsmithedin/overmyhouse@master by Jamie Smith passed")

				if matched {
					err := runDeploy(&dCfg)
					if err != nil {
						rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Deployment failed: %s", err), ev.Channel))
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage("Successfully deployed!", ev.Channel))
					}
				}
			}

		default:
			// Ignore other events..

		}
	}
}
