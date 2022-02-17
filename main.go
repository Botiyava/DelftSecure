package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	//Logging system
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

	//Telegram connection
	botToken := viper.GetString("botSettings.token")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil{
		logrus.Fatal("Error with bot authorization:", err.Error())
	}
	fmt.Printf("Successfully authorized on [ %s ]\n", bot.Self.UserName)
	//bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID))

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
