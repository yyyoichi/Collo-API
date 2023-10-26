package app

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestRunStream(t *testing.T) {
	opt := CollocationServiceOptions{Any: "防災", From: time.Now().AddDate(0, -6, 0), Until: time.Now().AddDate(0, -4, 0)}
	service, err := NewCollocationService(opt)
	if err != nil {
		t.Error(err)
		return
	}
	wordLength := 0
	pairLength := 0
	for pr := range service.Stream(context.Background()) {
		wordLength += len(pr.WordByID)
		pairLength += len(pr.Pairs)
	}
	log.Println(wordLength)
	log.Println(pairLength)
}
