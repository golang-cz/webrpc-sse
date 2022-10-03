package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"time"
)

func main() {
	// Subscribe to chat server and print all messages to STDOUT.
	go func() {
		if err := subscribe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Read each line from STDIN and send it as a message to the chat server.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := sendMessage(scanner.Text()); err != nil {
			log.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
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
		fmt.Printf("\t%s:\t%s\n", msg.Author, msg.Msg)
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
