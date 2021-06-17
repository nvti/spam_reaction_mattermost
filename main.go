package main

import (
	"errors"
	"flag"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
)

var client *model.Client4
var wg sync.WaitGroup

func main() {
	email := flag.String("email", "", "Mattermost login email")
	password := flag.String("pass", "", "Mattermost login password")
	post := flag.String("post", "", "Link of post you want to spam reaction :))")
	flag.Parse()

	if *email == "" || *password == "" || *post == "" {
		flag.Usage()
		os.Exit(1)
	}

	u, err := url.Parse(*post)
	if err != nil {
		log.Fatal("Wrong post link")
	}

	mattermost_server := u.Scheme + "://" + u.Host
	params := strings.Split(u.Path, "/")
	if len(params) == 0 {
		log.Fatal("Wrong post link")
	}

	postID := params[len(params)-1]

	client = model.NewAPIv4Client(mattermost_server)
	pingServer()
	userID := login(*email, *password)

	getAllEmoji()
	reactAll(userID, postID)

	wg.Wait()
	log.Print("Done")
}

func pingServer() {
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		log.Fatal("There was a problem pinging the Mattermost server.")
	} else {
		log.Print("Server detected and is running version " + props["Version"])
	}
}

func login(email string, pass string) string {
	user, resp := client.Login(email, pass)
	if resp.Error != nil {
		log.Fatal("There was a problem logging into the Mattermost server.")
	}
	return user.Id
}

func getAllEmoji() {
	support_emojis = append(support_emojis, getCustomEmoji()...)
}

func getCustomEmoji() (emojis []string) {
	emoji, resp := client.GetEmojiList(0, 60)
	if resp.Error != nil {
		return
	}
	for _, e := range emoji {
		emojis = append(emojis, e.Name)
	}

	return
}

func reactAll(userID string, postID string) {
	for _, emoji := range support_emojis {
		wg.Add(1)
		go react(userID, postID, emoji)
	}
}

func react(userID string, postID string, emoji string) error {
	defer wg.Done()
	react := &model.Reaction{
		UserId:    userID,
		PostId:    postID,
		EmojiName: emoji,
	}
	react.PreSave()
	_, resp := client.SaveReaction(react)

	if resp.Error != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("")
	}
	return nil
}
