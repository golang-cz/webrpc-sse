package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		// No CLI arguments - subscribe & print all messages.
		if err := subscribe(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Some CLI arguments - send them as a message.
	if err := sendMessage(strings.Join(os.Args[1:], " ")); err != nil {
		log.Fatal(err)
	}
}

func subscribe() error {
	ctx := context.Background()
	chatbot := NewChatClient("http://localhost:4242", http.DefaultClient)

	msgs, err := chatbot.SubscribeMessages(ctx)
	if err != nil {
		return err
	}

	for msg := range msgs {
		fmt.Printf("%s: %s\n", msg.Author, msg.Msg)

	}

	return nil
}

func sendMessage(msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	chatbot := NewChatClient("http://localhost:4242", http.DefaultClient)

	author := "Go client"
	if user, err := user.Current(); err == nil {
		author = user.Username
	}

	_, err := chatbot.SendMessage(ctx, author, msg)
	if err != nil {
		return err
	}

	return nil
}
