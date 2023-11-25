package pair

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	apiv1 "yyyoichi/Collo-API/internal/api/v1"
)

var tconfig Config
var thandler = Handler{
	Err: func(err error) {
		fmt.Println(err)
		panic(err)
	},
	Resp: func(resp *apiv1.ColloStreamResponse) {
	},
	Done: func() {},
}

func TestMain(m *testing.M) {
	initConfigMock()
	os.Exit(m.Run())
}
func BenchmarkCase0(b *testing.B) {
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case0(context.Background()) {
		}
	}
}

func BenchmarkCase1(b *testing.B) {
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case1(context.Background()) {
		}
	}
}

func BenchmarkCase2(b *testing.B) {
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case2(context.Background()) {
		}
	}
}

func BenchmarkCase3(b *testing.B) {
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case3(context.Background()) {
		}
	}
}

func BenchmarkCase4(b *testing.B) {
	ps, _ := NewPairStore(tconfig, thandler)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range ps.stream_case4(context.Background()) {
		}
	}
}

func initConfigMock() {
	config := Config{}
	config.Search.Any = "自動車"
	l, _ := time.LoadLocation("Asia/Tokyo")
	config.Search.From = time.Date(2022, 3, 1, 0, 0, 0, 0, l)
	config.Search.Until = time.Date(2022, 5, 1, 0, 0, 0, 0, l)
	tconfig = CreateMockConfig(config)
	fmt.Println("init config")
}