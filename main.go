package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"gopkg.in/yaml.v2"
)

// mattermost client
var client *model.Client4

// waitgroup for wait all go routine
var wg sync.WaitGroup

// user structure for parse user data from file
type user struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

func main() {
	// Create command line argument parser
	email := flag.String("email", "", "Mattermost login email")
	password := flag.String("pass", "", "Mattermost login password")
	file := flag.String("file", "", "File contain mattermost login email and password. Support json and yaml type")
	flag.Usage = usage
	flag.Parse()

	// Check post link is exist
	if flag.NArg() != 1 {
		println("Missing post link")
		flag.Usage()
		os.Exit(1)
	}
	post := flag.Args()[0]

	// Check cert file or email+pass is exist
	if *file == "" && (*email == "" || *password == "") {
		println("Missing certificate file or email/pass parametter")
		flag.Usage()
		os.Exit(1)
	}

	// Parse cert file if listed
	if *file != "" {
		var u user
		if filepath.Ext(*file) == ".json" {
			u = parseJson(*file)
		} else if filepath.Ext(*file) == ".yaml" || filepath.Ext(*file) == ".yml" {
			u = parseYaml(*file)
		} else {
			println("Un-support file type ", filepath.Ext(*file))
			flag.Usage()
			os.Exit(1)
		}

		email = &u.Email
		password = &u.Pass
	}

	// Parse mattermost server url and post id from post link
	u, err := url.Parse(post)
	if err != nil {
		println("Wrong post link")
		os.Exit(1)
	}

	mattermost_server := u.Scheme + "://" + u.Host
	params := strings.Split(u.Path, "/")
	if len(params) == 0 {
		println("Wrong post link")
		os.Exit(1)
	}
	postID := params[len(params)-1]

	// Connect and login
	client = model.NewAPIv4Client(mattermost_server)
	pingServer()
	userID := login(*email, *password)

	// get list of supported emoji and add reaction
	getAllEmoji()
	reactAll(userID, postID)

	// Wait all go routine
	wg.Wait()
	println("Done")
}

// Custom usage function
func usage() {
	println("Usage: ", os.Args[0], "[options] post_link")
	println("post_link: Link of post you want to spam reaction =))")
	println("Options:")
	flag.PrintDefaults()
	println("Login using certificate file or email+pass")
}

// Parse json file
func parseJson(file string) (u user) {
	jsonFile, err := os.Open(file)
	if err != nil {
		println("Can't open file: ", err.Error())
		os.Exit(1)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		println("Can't open file: ", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(byteValue, &u)
	if err != nil {
		println("Can't parse json file: ", err.Error())
		os.Exit(1)
	}

	return
}

// Parse yaml file
func parseYaml(file string) (u user) {
	yamlFile, err := os.Open(file)
	if err != nil {
		println("Can't open file: ", err.Error())
		os.Exit(1)
	}
	// defer the closing of our yamlFile so that we can parse it later on
	defer yamlFile.Close()

	byteValue, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		println("Can't open file: ", err.Error())
		os.Exit(1)
	}

	err = yaml.Unmarshal(byteValue, &u)
	if err != nil {
		println("Can't parse json file: ", err.Error())
		os.Exit(1)
	}

	return
}

// Check server is running
func pingServer() {
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.")
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}

// Login to server with email+pass
func login(email string, pass string) string {
	user, resp := client.Login(email, pass)
	if resp.Error != nil {
		println("There was a problem logging into the Mattermost server.")
		os.Exit(1)
	}
	return user.Id
}

// Get all emoji
func getAllEmoji() {
	support_emojis = append(support_emojis, getCustomEmoji()...)
}

// Get list of custom emoji in the server
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

// add all emoji to the post
func reactAll(userID string, postID string) {
	for _, emoji := range support_emojis {
		wg.Add(1)
		go react(userID, postID, emoji)
	}
}

// add one emoji to the post
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
