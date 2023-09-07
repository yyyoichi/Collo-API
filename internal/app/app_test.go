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
	if wordLength != 6236 {
		t.Errorf("Expected wordLength is 6236, but got='%d'", wordLength)
	}
	if pairLength != 18408185 {
		t.Errorf("Expected pairLength is 18408185, but got='%d'", pairLength)
	}
}

func BenchmarkStream(b *testing.B) {
	opt := CollocationServiceOptions{Any: "防災", From: time.Now().AddDate(0, -6, 0), Until: time.Now().AddDate(0, -4, 0)}
	service, err := NewCollocationService(opt)
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < 2; i++ {
		for range service.Stream(context.Background()) {
		}
	}
}

func BenchmarkStreamFun(b *testing.B) {
	opt := CollocationServiceOptions{Any: "防災", From: time.Now().AddDate(0, -6, 0), Until: time.Now().AddDate(0, -4, 0)}
	service, err := NewCollocationService(opt)
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < 2; i++ {
		for range service.StreamFun(context.Background()) {
		}
	}
}
