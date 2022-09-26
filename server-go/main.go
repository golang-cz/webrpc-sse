package main

import (
	"context"
	"log"
	"net/http"
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

	webrpcHandler := NewChatbotServer(&RPC{})
	r.Handle("/*", webrpcHandler)

	return http.ListenAndServe(":4242", r)
}

type RPC struct {
	msgLock sync.RWMutex
	msgId   uint64
	msgs    []*Message
}

func (s *RPC) SendMessage(ctx context.Context, author string, msg string) (bool, error) {
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

func (s *RPC) SubscribeMessages(ctx context.Context) ([]*Message, error) {
	s.msgLock.RLock()
	defer s.msgLock.RUnlock()

	return s.msgs, nil
}
