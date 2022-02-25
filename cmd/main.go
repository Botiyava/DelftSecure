package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"telegramBot/pkg/middleware"
	"telegramBot/pkg/repository"
	"time"
)

var db *sqlx.DB

func main() {
	// Logging system
	logrus.SetFormatter(new(logrus.JSONFormatter))

	file, err := os.OpenFile("log/errors.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	logrus.SetOutput(file)
	if err := initConfig(); err != nil {
		logrus.Fatal("Error initializing configs: ", err.Error())
	}

	// Telegram bot
	botToken := viper.GetString("botSettings.token")
	pref := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfully connected to [", b.Me.Username, "].")
	defer b.Close()
	b.SetCommands("/new", "/start")

	// Database connection
	db, err = repository.NewPostgresDB(repository.Config{
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize db: %s", err.Error())
	}
	fmt.Println("Successfully connected to [ database ].")
	fmt.Println("Now bot is working...")
	defer db.Close()

	b.Handle("/start", startHandler)
	b.Handle("/help", startHandler)
	b.Handle("/new", newHandler, middleware.ValidationNew)
	b.Handle("/get", getHandler, middleware.ValidationGet)

	b.Start()

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func startHandler(c tele.Context) error {
	return c.Send("Hello! Here you can store your password " +
		"and get it at any time.\n\n" +
		"[+] Save new password:\n" +
		"    /new <url> <login> <password>\n" +
		"Example:\n    /new site.com username mypassword")

}

type dbRecord struct {
	userID   int64
	url      string
	login    string
	password string
}

func newHandler(c tele.Context) error {
	tags := c.Args()
	userID := c.Sender().ID
	userURL := tags[0]
	userLogin := tags[1]
	userPassword := tags[2]
	var exists bool
	//if user has old password for this url and login
	err := db.QueryRow("SELECT exists (SELECT * FROM password_storage WHERE userid = $1 AND url = $2 AND login = $3)",
		userID, userURL, userLogin).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		logrus.Fatalf("newHandler: Failed check row existance: %s", err.Error())
	}

	// if user already has record linked to required url && login
	if exists {
		_, err = db.Exec("UPDATE password_storage SET password = $1 WHERE userid = $2 AND url = $3 AND login = $4",
			userPassword, userID, userURL, userLogin)
		if err != nil {
			logrus.Fatalf("newHandler: Failed to update new record in db: %s", err.Error())
			c.Send("Sorry, I can't update password to " + userLogin + " on " + userURL + "\nPlease try again later.")
			return err
		}
		c.Send("Successfully updated password to user " + userLogin + " on " + userURL)
	} else {
		fmt.Printf("%v, %s, %s, %s\n", c.Sender().ID, userURL, userLogin, userPassword)
		_, err = db.Exec("INSERT INTO password_storage(userid, url, login, password) VALUES ($1,$2,$3,$4)",
			userID, userURL, userLogin, userPassword)
		if err != nil {
			logrus.Fatalf("newHandler: Failed to add new record in db: %s", err.Error())
			c.Send("Sorry, I can't add your record\nnPlease try again later.")
			return err
		}
		c.Send("Successfully saved password to user " + userLogin + " on " + userURL)
	}
	c.DeleteAfter(10 * time.Second)
	return nil
}
func getHandler(c tele.Context) error{
	tags := c.Args()
	userID := c.Sender().ID
	userURL := tags[0]
	userLogin := tags[1]

	var password string
	err := db.QueryRow("SELECT password FROM password_storage WHERE userid = $1 AND url = $2 AND login = $3",
		userID, userURL, userLogin).Scan(&password)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(password)
	c.Send("Login: " + userLogin + " Password: " + password)
	c.DeleteAfter(10 * time.Second)
return nil
}