package main

import(
  "fmt"
  "github.com/slack-go/slack"
)

func main(){
  api := slack.New("xoxb-23846488290-1288404371285-bizzoEHUDRO94eowq4GOykIR")

  // If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// slack.New("YOUR_TOKEN_HERE", slack.OptionDebug(true))


  channelID, timestamp, err := api.PostMessage(
    "bot-benchmarks",
    slack.MsgOptionText("test", false),
    slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
  )
  if err != nil {
    fmt.Printf("%s\n", err)
    return
  }
  fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

}


