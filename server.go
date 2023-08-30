package main

import (
	"context"
	"fmt"
	"log"
	"time"
	collov1 "yyyoichi/Collo-API/gen/proto/collo/v1"

	"connectrpc.com/connect"
)

type ColloServer struct{}

func (*ColloServer) ColloStream(cxt context.Context, req *connect.Request[collov1.ColloRequest], str *connect.ServerStream[collov1.ColloStreamResponse]) error {
	log.Printf("Get Request: %s\n", req.Header())
	for i := 0; i < 5; i++ {
		i32 := int32(i)
		// make word dummy
		words := make(map[int32]string)
		words[i32] = fmt.Sprintf("New Word: 'word%d'", i)
		// make pairs dummy
		pairs := make([]*collov1.Pair, 1)
		pairs = append(pairs, &collov1.Pair{Values: []int32{i32, i32}})
		if err := str.Send(&collov1.ColloStreamResponse{Words: words, Pairs: pairs}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}
