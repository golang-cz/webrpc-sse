package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	webrpcHandler := NewChatbotServer(&RPC{
		msgId: 3,
		msgs: []*Message{
			{
				ID:     1,
				Author: "Test",
				Msg:    "Message 1",
			},
			{
				ID:     2,
				Author: "Test",
				Msg:    "Message 2",
			},
			{
				ID:     3,
				Author: "Test",
				Msg:    "Message 3",
			},
		},
	})
	r.Handle("/*", webrpcHandler)

	return http.ListenAndServe(":4242", r)
}

type RPC struct {
	msgLock sync.RWMutex
	msgId   uint64
	msgs    []*Message
}

func (s *RPC) generateMessage() *Message {
	s.msgLock.Lock()
	defer s.msgLock.Unlock()

	s.msgId++
	msg := &Message{
		ID:     s.msgId,
		Author: "Test",
		Msg:    fmt.Sprintf("Message %v", s.msgId),
	}
	s.msgs = append(s.msgs, msg)
	return msg
}

func (s *RPC) SendMessage(ctx context.Context, author string, msg string) (bool, error) {
	author = strings.TrimSpace(author)
	if msg == "" {
		return false, ErrorInvalidArgument("author", "empty author")
	}

	msg = strings.TrimSpace(msg)
	if msg == "" {
		return false, ErrorInvalidArgument("msg", "empty message")
	}

	s.msgLock.RLock()
	defer s.msgLock.RUnlock()

	s.msgId++
	s.msgs = append(s.msgs, &Message{
		ID:        s.msgId,
		Author:    author,
		CreatedAt: time.Now(),
		Msg:       msg,
	})
	return true, nil
}

func (s *RPC) SubscribeMessages(ctx context.Context) (chan *Message, error) {
	msgs := make(chan *Message, 100)

	go func() {
		func() {
			s.msgLock.RLock()
			defer s.msgLock.RUnlock()

			for _, msg := range s.msgs {
				msgs <- msg
			}
		}()

		for {
			time.Sleep(time.Second * time.Duration(rand.Intn(15)))
			msgs <- s.generateMessage()
		}
	}()

	return msgs, nil
}
