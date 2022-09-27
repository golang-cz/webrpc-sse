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
	"github.com/rs/cors"
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

	rpcServer := &RPC{
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
	}

	webrpcHandler := NewChatServer(rpcServer)
	r.Handle("/*", webrpcHandler)

	go func() {
		// Generate random messages.
		for {
			time.Sleep(time.Second * time.Duration(rand.Intn(15)))
			_, _ = rpcServer.SendMessage(context.Background(), "Test", fmt.Sprintf("Random message"))
		}
	}()

	return http.ListenAndServe(":4242", r)
}

type RPC struct {
	// Store with all messages.
	msgLock sync.RWMutex
	msgId   uint64
	msgs    []*Message

	// Subscriptions - each SSE client can subscribe to new messages.
	subsLock sync.RWMutex
	subs     []chan *Message
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

	s.msgId++
	message := &Message{
		ID:        s.msgId,
		Author:    author,
		CreatedAt: time.Now(),
		Msg:       msg,
	}

	s.msgLock.Lock()
	defer s.msgLock.Unlock()
	s.msgs = append(s.msgs, message)

	s.subsLock.RLock()
	defer s.subsLock.RUnlock()
	for _, subscriber := range s.subs {
		subscriber <- message
	}

	return true, nil
}

func (s *RPC) SubscribeMessages(ctx context.Context) (chan *Message, error) {
	msgs := make(chan *Message, 100)

	go func() {
		func() {
			s.msgLock.RLock()
			defer s.msgLock.RUnlock()

			// Print all messages.
			for _, msg := range s.msgs {
				msgs <- msg
			}
		}()

		// Subscribe to new messages.
		newMsgs, unsubscribe := s.subscribe()

		for {
			select {
			case msg := <-newMsgs:
				msgs <- msg

			case <-ctx.Done():
				unsubscribe()
				return
			}
		}
	}()

	return msgs, nil
}

// Subscribe returns a channel with all new messages and unsubscribe() function.
func (s *RPC) subscribe() (chan *Message, func()) {
	sub := make(chan *Message, 10)

	s.subsLock.Lock()
	defer s.subsLock.Unlock()

	s.subs = append(s.subs, sub)

	return sub, func() {
		s.subsLock.Lock()
		defer s.subsLock.Unlock()

		for i, subscription := range s.subs {
			if sub == subscription {
				// Remove subscription.
				s.subs = append(s.subs[:i], s.subs[i+1:]...)
			}
		}

		close(sub)
	}
}
