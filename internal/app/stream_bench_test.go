package app

import (
	"context"
	"testing"
	"time"
)

var opt = CollocationServiceOptions{Any: "防災", From: time.Now().AddDate(0, -6, 0), Until: time.Now().AddDate(0, -4, 0)}

func BenchmarkStream(b *testing.B) {
	service, err := NewCollocationService(opt)
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range service.Stream(context.Background()) {
		}
	}
}

func BenchmarkStreamFun(b *testing.B) {
	service, err := NewCollocationService(opt)
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range service.StreamFun(context.Background()) {
		}
	}
}

func BenchmarkDemultiParseStream(b *testing.B) {
	service, err := NewCollocationService(opt)
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range service.DemultiParseStream(context.Background()) {
		}
	}
}

func BenchmarkDemultiFunParseStream(b *testing.B) {
	service, err := NewCollocationService(opt)
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range service.DemultiFunParseStream(context.Background()) {
		}
	}
}

func BenchmarkDemultiFunFetchParseStream(b *testing.B) {
	service, err := NewCollocationService(opt)
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range service.DemultiFunFetchParseStream(context.Background()) {
		}
	}
}

func BenchmarkDemultiFunStream(b *testing.B) {
	service, err := NewCollocationService(opt)
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range service.DemultiFunStream(context.Background()) {
		}
	}
}
