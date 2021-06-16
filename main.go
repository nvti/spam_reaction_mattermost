package main

import (
	"os"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	MATTERMOST_SERVER = ""
	USER_EMAIL        = ""
	USER_PASSWORD     = ""
)

var client *model.Client4

var botUser *model.User

// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func main() {

	client = model.NewAPIv4Client(MATTERMOST_SERVER)

	// Lets test to see if the mattermost server is up and running
	MakeSureServerIsRunning()

	// lets attempt to login to the Mattermost server as the bot user
	// This will set the token required for all future calls
	// You can get this token with client.AuthToken
	LoginAsTheBotUser()

	ReactMess("qr63i64ynbdwims4nwps6o4p5r")
	// You can block forever with
	select {}
}

func MakeSureServerIsRunning() {
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}

func LoginAsTheBotUser() {
	if user, resp := client.Login(USER_EMAIL, USER_PASSWORD); resp.Error != nil {
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		os.Exit(1)
	} else {
		botUser = user
	}
}

func ReactMess(postID string) {
	react := &model.Reaction{
		UserId:    botUser.Id,
		PostId:    postID,
		EmojiName: "grinning",
	}

	react.PreSave()

	if _, resp := client.SaveReaction(react); resp.Error != nil {
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
	}
}
