package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	queuer "github.com/defryheryanto/job-queuer"
)

type Simple struct {
	Title string
}

func (s *Simple) GetTitle() string {
	return s.Title
}

func (s *Simple) Do(ctx context.Context) error {
	time.Sleep(1 * time.Second)
	return nil
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	q := queuer.NewQueuer(2)
	q.Run(context.Background())

	q.Push(&Simple{Title: "Task 1"})
	q.Push(&Simple{Title: "Task 2"})
	q.Push(&Simple{Title: "Task 3"})
	q.Push(&Simple{Title: "Task 4"})
	q.Push(&Simple{Title: "Task 5"})

	<-c
}
