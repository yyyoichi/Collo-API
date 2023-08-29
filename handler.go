package main

import (
	"fmt"
	"time"
	"yyyoichi/Collo-API/proto/collo"
)

type Server struct {
	collo.UnimplementedColloServiceServer
}

func (*Server) ColloStream(req *collo.ColloRequest, stream collo.ColloService_ColloStreamServer) error {
	for i := 0; i < 5; i++ {
		i32 := int32(i)
		// make word dummy
		words := make(map[int32]string)
		words[i32] = fmt.Sprintf("New Word: 'word%d'", i)
		// make pairs dummy
		pairs := make([]*collo.Pair, 1)
		pairs = append(pairs, &collo.Pair{Values: []int32{i32, i32}})

		resp := &collo.ColloResponse{
			Words: words,
			Pairs: pairs,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}
