package main

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"regexp"
	"telegramBot/repository"
	"time"
)

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
		Token: botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfully connected to [", b.Me.Username, "].")
	defer b.Close()
	// Database connection
	db, err := repository.NewPostgresDB(repository.Config{
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


	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello! Here you can store your password " +
			"and get it at any time.\n\n" +
			"[+] Save new password:\n" +
			"    /new <url> <login> <password>\n" +
			"Example:\n    /new site.com username mypassword")

	})
	b.Handle("/new", func(c tele.Context) error {
		tags := c.Args()
		if len(tags) != 3{
			c.Send("You need to send 3 arguments: URL, login and password.\n" +
				" Send /help to see the format of commands\n")
			return errors.New("")
		}
		if matched, _ := regexp.MatchString(`^[a-z0-9]{1,30}\.[a-z]{1,3}$`, tags[0]); !matched{
			c.Send("Your URL is invalid. Try something like example.com, mywebsite.net or etc\n")
		}
		for _, tag := range tags {
			fmt.Println(tag)
		}
		//TODO check if len of login and password strings <= len of the same fields in DB

		return nil
	})


	b.Start()


}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
