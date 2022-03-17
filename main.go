package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hellflame/argparse"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

var Login = new(FutilUser)

type FutilUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const LoginFIle = "login.json"

var mutex sync.Mutex

func WriteJSONToFile() {
	database, err := os.OpenFile(LoginFIle, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	data, err := json.MarshalIndent(Login, "", "\t")

	err = ioutil.WriteFile(LoginFIle, data, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadJSONToStruct() {
	content, _ := ioutil.ReadFile(LoginFIle)
	if len(content) <= 1 {
		mainModel, _ := json.MarshalIndent(FutilUser{}, "", "\t")
		err := ioutil.WriteFile(LoginFIle, mainModel, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err := json.Unmarshal(content, &Login)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

type User struct {
	ID               uuid.UUID `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
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

	parser := argparse.NewParser("futil", "Internal Flappy Bird moderation tool. If you do not have permission to use this, you should not.", nil)

	id := parser.String("i", "id", &argparse.Option{Required: false, Help: "Get the uuid from a name"})
	list := parser.Flag("l", "list", &argparse.Option{Required: false, Help: "Lists all of the users, and their score"})

	ban := parser.Strings("b", "ban", &argparse.Option{Required: false, Help: "Ban a user with reason", Positional: true})
	unban := parser.String("u", "unban", &argparse.Option{Required: false, Help: "Unban a user"})
	restoreScore := parser.Strings("r", "restore", &argparse.Option{Required: false, Help: "Restore a users score, [name] [score]", Positional: true})
	logs := parser.Flag("", "log", &argparse.Option{Required: false, Help: "Shows the server log"})
	jailbroken := parser.Flag("", "list-jail", &argparse.Option{Required: false, Help: "Shows jailbroken users"})

	delete := parser.String("d", "delete", &argparse.Option{Required: false, Help: "Delete a user"})

	admin := parser.String("a", "admin", &argparse.Option{Required: false, Help: "Make a user a admin"})

	loginFIle, err := os.OpenFile(LoginFIle, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer loginFIle.Close()

	ReadJSONToStruct()

	if len(Login.Username) < 1 || len(Login.Password) < 1 {
		log.Fatal("You must add a username and password in login.json")
	}

	err = parser.Parse(os.Args)
	if err != nil || len(os.Args) < 2 {
		parser.PrintHelp()
		return
	}

	if len(*ban) > 1 {
		args := *ban
		fmt.Println(fmt.Sprintf("Banning %s Reason: %s", args[2], args[3]))
		callApi(fmt.Sprintf("auth/ban/%s/%s", GetID(args[2], Login.Username, Login.Password), url.QueryEscape(args[3])), Login.Username, Login.Password)
		fmt.Println("Banned", args[2])
	}

	if hasStrArg(unban) {
		fmt.Println("Unbanning", *unban)
		callApi("auth/unban/"+GetID(*unban, Login.Username, Login.Password), Login.Username, Login.Password)
		fmt.Println("Unbanned", *unban)
	}

	if hasStrArg(delete) {
		fmt.Println("Deleting", *delete)
		callApi("auth/delete/"+GetID(*delete, Login.Username, Login.Password), Login.Username, Login.Password)
		fmt.Println("Deleted", *delete)
	}

	if len(*restoreScore) > 1 {
		args := *restoreScore
		score, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatal("Could not parse the score")
			return
		}
		fmt.Println("Setting", args[2]+"'s score to", score)
		callApi(fmt.Sprintf("auth/restoreScore/%s/%d", GetID(args[2], Login.Username, Login.Password), score), Login.Username, Login.Password)
		fmt.Println("Set", args[2]+"'s score to", score)
	}

	if hasStrArg(id) {
		fmt.Println("Getting", *id+"'s id")
		fmt.Println("ID:", GetID(*id, Login.Username, Login.Password))
	}

	if *list {
		var users []User
		body := callApi("auth/internalUsers", Login.Username, Login.Password)
		json.Unmarshal(body, &users)
		for _, user := range users {
			fmt.Println("------------------------------------------------------------------")
			fmt.Println(user.Name + "\t\tScore: " + strconv.Itoa(user.Score) + "\t\tDeaths: " + strconv.Itoa(user.Deaths))
		}
		fmt.Println("------------------------------------------------------------------")
	}

	if *jailbroken {
		var users []User
		body := callApi("auth/internalUsers", Login.Username, Login.Password)
		json.Unmarshal(body, &users)
		var hasAtLestOneUser bool
		for _, user := range users {
			if user.Jailbroken {
				hasAtLestOneUser = true
				fmt.Println("------------------------------------------------------------------")
				fmt.Println(user.Name + "\t\tScore: " + strconv.Itoa(user.Score) + "\t\tDeaths: " + strconv.Itoa(user.Deaths))
			}
		}
		if hasAtLestOneUser {
			fmt.Println("------------------------------------------------------------------")
		} else {
			fmt.Println("No jailbroken users")
		}
	}

	if *logs {
		body := callApi("auth/logs", Login.Username, Login.Password)
		fmt.Println(string(body))
	}

	if hasStrArg(admin) {
		fmt.Println("Making", *admin, "a admin")
		callApi("auth/makeAdmin/"+GetID(*admin, Login.Username, Login.Password), Login.Username, Login.Password)
		fmt.Println("Made", *admin, "a admin")
	}
}
