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
	if len(os.Args) > 2 {
		if err := sendMessage(strings.Join(os.Args[1:], " ")); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := subscribe(); err != nil {
		log.Fatal(err)
	}
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

func subscribe() error {
	ctx := context.Background()

	chatbot := NewChatClient("http://localhost:4242", http.DefaultClient)
	msgs, err := chatbot.SubscribeMessages(ctx)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		fmt.Printf("%s: %s", msg.Author, msg.Msg)
	}

	return nil
}
