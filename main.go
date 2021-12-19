package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name             string `json:"name" gorm:"unique"`
	Score            int    `json:"score"`
	Deaths           int    `json:"deaths"`
	PasswordHash     string `json:"passwordHash"`
	Jailbroken       bool   `json:"jailbroken"`
	HasHackedTools   bool   `json:"hasHackedTools"`
	RanInEmulator    bool   `json:"ranInEmulator"`
	HasModifiedScore bool   `json:"hasModifiedScore"`
	IsBanned         bool   `json:"isBanned"`
	BanReason        string `json:"banReason"`
	Admin            bool   `json:"admin"`
	Owner            bool   `json:"owner"`
}

const (
	apiUrl = "https://flappybird.brandonplank.org/v1/"
)

func CraftBasicAuthHeader(username string, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

func callApi(endpoint string, username string, password string) []byte {
	request, err := http.NewRequest("GET", apiUrl+endpoint, nil)
	if len(username) > 0 && len(password) > 0 {
		request.Header.Set("Authorization", "Basic "+CraftBasicAuthHeader(username, password))
	}
	if err != nil {
		log.Fatal(err)
	}
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)

	if resp.StatusCode == 401 {
		log.Fatal("You are not authorized to perform this action")
		os.Exit(-1)
	}

	if resp.StatusCode == 500 {
		log.Fatal("This was not supposed to happen, internal server error")
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func GetID(name string, username string, password string) string {
	body := callApi("getID/"+name, username, password)
	return string(body)
}

func hasStrArg(s *string) bool {
	return len(*s) > 0
}

func main() {

	parser := argparse.NewParser("futil", "Internal Flappy Bird moderation tool\\nIf you do not have permission to use this, you should not.")

	id := parser.String("i", "id", &argparse.Options{Required: false, Help: "Get the uuid from a name"})
	list := parser.Flag("l", "list", &argparse.Options{Required: false, Help: "Lists all of the users, and their score"})

	ban := parser.String("b", "ban", &argparse.Options{Required: false, Help: "Ban a user"})
	unban := parser.String("u", "unban", &argparse.Options{Required: false, Help: "Unban a user"})
	restoreScore := parser.StringList("r", "restore", &argparse.Options{Required: false, Help: "Restore a users score, [name] [score]"})

	delete := parser.String("d", "delete", &argparse.Options{Required: false, Help: "Delete a user"})

	admin := parser.String("a", "admin", &argparse.Options{Required: false, Help: "Make a user a admin"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, please set USERNAME and PASSWORD")
		return
	}

	username := os.Getenv("USERNAME")
	if len(username) < 1 {
		log.Fatal("You must set your username in the .env file!")
		return
	}
	password := os.Getenv("PASSWORD")
	if len(password) < 1 {
		log.Fatal("You must set your password in the .env file!")
		return
	}

	if hasStrArg(ban) {
		fmt.Println("Banning", *ban)
		callApi("auth/ban/"+GetID(*ban, username, password), username, password)
		fmt.Println("Banned", *ban)
	}

	if hasStrArg(unban) {
		fmt.Println("Unbanning", *unban)
		callApi("auth/unban/"+GetID(*unban, username, password), username, password)
		fmt.Println("Unbanned", *unban)
	}

	if hasStrArg(delete) {
		fmt.Println("Deleting", *delete)
		callApi("auth/delete/"+GetID(*delete, username, password), username, password)
		fmt.Println("Deleted", *delete)
	}

	if len(*restoreScore) > 1 && len(*restoreScore) < 3 {
		args := *restoreScore
		score, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal("Could not parse the score")
			return
		}
		fmt.Println("Setting", args[0]+"'s score to", score)
		callApi("auth/restoreScore/"+GetID(args[0], username, password), username, password)
		fmt.Println("Set", args[0]+"'s score to", score)
	}

	if hasStrArg(id) {
		fmt.Println("Getting", *id+"'s id")
		fmt.Println("ID:", GetID(*id, username, password))
	}

	if *list {
		var users []User
		body := callApi("auth/internalUsers", username, password)
		json.Unmarshal(body, &users)
		for i := 0; i < len(users); i++ {
			fmt.Println("------------------------------------------------------------------")
			fmt.Println(users[i].Name + "\t\tScore: " + strconv.Itoa(users[i].Score) + "\t\tDeaths: " + strconv.Itoa(users[i].Deaths))
		}
		fmt.Println("------------------------------------------------------------------")
	}

	if hasStrArg(admin) {
		fmt.Println("Making", *admin, "a admin")
		callApi("auth/makeAdmin/"+GetID(*admin, username, password), username, password)
		fmt.Println("Made", *admin, "a admin")
	}
}
