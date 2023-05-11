package main

import (
	"time"

	keylogger "github.com/alex-grimes/go-klogger"
	"github.com/sirupsen/logrus"
)

func main() {

	keyboard := keylogger.FindKeyboardDevice()

	if len(keyboard) <= 0 {
		logrus.Fatal("No keyboard found...you will need to provide the path manually")
		return
	}

	logrus.Println("Found a keyboard at", keyboard)

	k, err := keylogger.New(keyboard)
	if err != nil {
		logrus.Fatal(err)
		return
	}
	defer k.Close()

	go func() {
		time.Sleep(5 * time.Second)

		keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "ENTER"}
		for _, key := range keys {
			k.WriteOnce(key)
		}
	}()

	events := k.Read()

	for e := range events {
		switch e.Type {
		case keylogger.EvKey:
			if e.KeyPress() {
				logrus.Printf("Key press detected %s", e.KeyString())
			}
			if e.KeyRelease() {
				logrus.Printf("Key release detected %s", e.KeyString())
			}
			break
		}
	}
}
