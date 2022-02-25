package middleware

import (
	"errors"
	tele "gopkg.in/telebot.v3"
	"regexp"
)

func ValidationGet(next tele.HandlerFunc) tele.HandlerFunc{
	return func(c tele.Context) error {
		tags := c.Args()
		if len(tags) != 2{
			c.Send("You need to send 2 arguments: URL and login.\n" +
				" Send /help to see the format of commands\n")
			return errors.New("middleware.Validation1: Number of arguments not equal 3")
		}
		if matched, _ := regexp.MatchString(`^[a-z0-9]{1,30}\.[a-z]{1,3}$`, tags[0]); !matched{
			c.Send("Your URL is invalid. Try something like example.com, mywebsite.net or etc\n")
			return errors.New("middleware.Validation1: Invalid URL")
		}

		for _, tag := range tags {
			if len(tag) > 63{
				c.Send("Your arguments must have length in range from 1 to 62 symbols")
			}
		}
		return next(c) // continue execution chain
	}
}

func ValidationNew(next tele.HandlerFunc) tele.HandlerFunc{
	return func(c tele.Context) error {
		tags := c.Args()
		if len(tags) != 3{
			c.Send("You need to send 3 arguments: URL, login and password.\n" +
				" Send /help to see the format of commands\n")
			return errors.New("middleware.Validation2: Number of arguments not equal 3")
		}
		if matched, _ := regexp.MatchString(`^[a-z0-9]{1,30}\.[a-z]{1,3}$`, tags[0]); !matched{
			c.Send("Your URL is invalid. Try something like example.com, mywebsite.net or etc\n")
			return errors.New("middleware.Validation2: Invalid URL")
		}

		for _, tag := range tags {
			if len(tag) > 63{
				c.Send("Your arguments must have length in range from 1 to 62 symbols")
			}
		}
		return next(c) // continue execution chain
	}
}