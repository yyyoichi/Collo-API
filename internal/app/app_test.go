package app

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestRunStream(t *testing.T) {
	opt := CollocationServiceOptions{Any: "防災", From: time.Now().AddDate(0, -6, 0), Until: time.Now()}
	service, err := NewCollocationService(opt)
	if err != nil {
		t.Error(err)
		return
	}
	for l := range service.Stream(context.Background()) {
		log.Printf("Get Stream Resp: %d pairs\n", len(l.Pairs))
	}
}
