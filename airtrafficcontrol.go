package main

import (
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

var configFilePath = flag.String("cfgFile", "./config.yml", "Path to deploy config")

func main() {
	flag.Parse()

	var deployConfig DeployConfig
	deployConfig.loadConfig(*configFilePath)

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
			log.Printf("Message: %v\n", ev)

			info := rtm.GetInfo()

			text := ev.Text
			text = strings.TrimSpace(text)
			text = strings.ToLower(text)

			matched, _ := regexp.MatchString("deploy", text)

			if ev.User != info.User.ID && matched {
				err := runDeploy(&deployConfig)
				if err != nil {
					rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Deployment failed: %s", err), ev.Channel))
				} else {
					rtm.SendMessage(rtm.NewOutgoingMessage("Successfully deployed!", ev.Channel))
				}
			}

		default:
			// Ignore other events..

		}
	}
}
