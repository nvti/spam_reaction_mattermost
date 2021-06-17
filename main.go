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
	email := flag.String("email", "", "Mattermost login email")
	password := flag.String("pass", "", "Mattermost login password")
	file := flag.String("file", "", "File contain mattermost login email and password. Support json and yaml type")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		println("Missing post link")
		flag.Usage()
		os.Exit(1)
	}
	post := flag.Args()[0]

	if *file == "" && (*email == "" || *password == "") {
		println("Missing certificate file or email/pass parametter")
		flag.Usage()
		os.Exit(1)
	}
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

	client = model.NewAPIv4Client(mattermost_server)
	pingServer()
	userID := login(*email, *password)

	getAllEmoji()
	reactAll(userID, postID)

	wg.Wait()
	println("Done")
}

func usage() {
	println("Usage: ", os.Args[0], "[options] post_link")
	println("post_link: Link of post you want to spam reaction =))")
	println("Options:")
	flag.PrintDefaults()
	println("Login using certificate file or email+pass")
}

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

func pingServer() {
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.")
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}

func login(email string, pass string) string {
	user, resp := client.Login(email, pass)
	if resp.Error != nil {
		println("There was a problem logging into the Mattermost server.")
		os.Exit(1)
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
