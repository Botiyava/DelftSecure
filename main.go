package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))
	file, err := os.OpenFile("log/errors.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	//logger.SetOutput(file)
	logrus.SetOutput(file)
	if err := initConfig(); err != nil {
		logrus.Fatal("Error initializing configs: ", err.Error())
	}

	a := viper.GetString("botSettings.token")
	fmt.Println(a)
	fmt.Println("lel")


}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("confkig")
	return viper.ReadInConfig()
}
