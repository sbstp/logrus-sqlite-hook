package main

import (
	"log"

	logrussqlitehook "github.com/sbstp/logrus-sqlite-hook"
	"github.com/sirupsen/logrus"
)

func main() {
	hook, err := logrussqlitehook.New("log.db")
	if err != nil {
		log.Fatal(err)
	}
	defer hook.Close()
	logrus.AddHook(hook)

	logrus.WithField("geez", "goo").Info("kachow my dude")
}
