package handler

import (
	"context"
	"fmt"
	"yyyoichi/Collo-API/internal/analyzer"
	apiv3 "yyyoichi/Collo-API/internal/api/v3"
	"yyyoichi/Collo-API/internal/api/v3/apiv3connect"
	"yyyoichi/Collo-API/internal/matrix"
	"yyyoichi/Collo-API/internal/ndl"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CoMatrixes []*matrix.CoMatrix

// implement apiv3connect.MintGreenServiceHandler
type V3Handler struct {
	apiv3connect.MintGreenServiceHandler
}

func (*V3Handler) NetworkStream(
	ctx context.Context,
	req *connect.Request[apiv3.NetworkStreamRequest],
	stream *connect.ServerStream[apiv3.NetworkStreamResponse],
) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	handleErr := func(err error) {
		select {
		case <-ctx.Done():
			return
		default:
			cancel(err)
		}
	}
	handleProcessResp := func(process float32) {
		resp := &apiv3.NetworkStreamResponse{
			Nodes:   []*apiv3.Node{},
			Edges:   []*apiv3.Edge{},
			Meta:    &apiv3.Meta{},
			Process: process,
		}
		select {
		case <-ctx.Done():
			return
		default:
			if err := stream.Send(resp); err != nil {
				cancel(err)
			}
		}
	}
	handleNetworkResp := func(nodes []*matrix.Node, edges []*matrix.Edge, meta *matrix.MultiDocMeta) {
		resp := &apiv3.NetworkStreamResponse{
			Nodes:   []*apiv3.Node{},
			Edges:   []*apiv3.Edge{},
			Meta:    &apiv3.Meta{},
			Process: 1,
		}
		resp.Meta = &apiv3.Meta{
			GroupId: meta.GroupID,
			From:    timestamppb.New(meta.From),
			Until:   timestamppb.New(meta.Until),
			Metas:   make([]*apiv3.DocMeta, len(meta.Metas)),
		}
		for i, dmeta := range meta.Metas {
			resp.Meta.Metas[i] = &apiv3.DocMeta{
				GroupId:     dmeta.GroupID,
				Key:         dmeta.Key,
				Name:        dmeta.Name,
				At:          timestamppb.New(dmeta.At),
				Description: dmeta.Description,
			}
		}
		for _, node := range nodes {
			resp.Nodes = append(resp.Nodes, &apiv3.Node{
				NodeId:   uint32(node.ID),
				Word:     string(node.Word),
				Rate:     float32(node.Rate),
				NumEdges: 0, // TODO
			})
		}
		for _, edge := range edges {
			resp.Edges = append(resp.Edges, &apiv3.Edge{
				EdgeId:  uint32(edge.ID),
				NodeId1: uint32(edge.Node1ID),
				NodeId2: uint32(edge.Node2ID),
				Rate:    float32(edge.Rate),
			})
		}
		select {
		case <-ctx.Done():
			return
		default:
			if err := stream.Send(resp); err != nil {
				cancel(err)
			}
		}
	}

	coMatrixes := NewCoMatrixes(
		ctx,
		ProcessHandler{
			Err:  handleErr,
			Resp: handleProcessResp,
		},
		NewConfig(req.Msg.Config),
	)
	select {
	case <-ctx.Done():
	default:
		if req.Msg.ForcusNodeId == uint32(0) {
			for _, cm := range coMatrixes {
				top1 := cm.NodeRank(0)
				top2 := cm.NodeRank(1)
				top3 := cm.NodeRank(2)
				nodes, edges := cm.CoOccurrences(top1.ID, top2.ID, top3.ID)
				nodes = append(nodes, top1, top2, top3)
				handleNetworkResp(nodes, edges, cm.Meta())
			}
		} else {
			for _, cm := range coMatrixes {
				nodes, edges := cm.CoOccurrences(uint(req.Msg.ForcusNodeId))
				handleNetworkResp(nodes, edges, cm.Meta())
			}
		}
		cancel(nil)
	}

	err := context.Cause(ctx)
	return responseError(err)
}

func (*V3Handler) NodeRateStream(
	ctx context.Context,
	req *connect.Request[apiv3.NodeRateStreamRequest],
	stream *connect.ServerStream[apiv3.NodeRateStreamResponse],
) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	handleErr := func(err error) {
		select {
		case <-ctx.Done():
			return
		default:
			cancel(err)
		}
	}
	handleProcessResp := func(process float32) {
		resp := &apiv3.NodeRateStreamResponse{
			Nodes:   []*apiv3.Node{},
			Meta:    &apiv3.Meta{},
			Process: process,
		}
		select {
		case <-ctx.Done():
			return
		default:
			if err := stream.Send(resp); err != nil {
				cancel(err)
			}
		}
	}
	handleNodeRateResp := func(resp *apiv3.NodeRateStreamResponse) {
		select {
		case <-ctx.Done():
			return
		default:
			if err := stream.Send(resp); err != nil {
				cancel(err)
			}
		}
	}
	coMatrixes := NewCoMatrixes(
		ctx,
		ProcessHandler{
			Err:  handleErr,
			Resp: handleProcessResp,
		},
		NewConfig(req.Msg.Config),
	)
	for _, cm := range coMatrixes {
		resp := &apiv3.NodeRateStreamResponse{
			Nodes:   []*apiv3.Node{},
			Meta:    &apiv3.Meta{},
			Num:     uint32(cm.LenNodes()),
			Next:    0,
			Count:   0,
			Process: 1,
		}
		meta := cm.Meta()
		resp.Meta = &apiv3.Meta{
			GroupId: meta.GroupID,
			From:    timestamppb.New(meta.From),
			Until:   timestamppb.New(meta.Until),
			Metas:   make([]*apiv3.DocMeta, len(meta.Metas)),
		}
		for i, dmeta := range meta.Metas {
			resp.Meta.Metas[i] = &apiv3.DocMeta{
				GroupId:     dmeta.GroupID,
				Key:         dmeta.Key,
				Name:        dmeta.Name,
				At:          timestamppb.New(dmeta.At),
				Description: dmeta.Description,
			}
		}

		offset := 0
		if req.Msg.Offset > 0 {
			offset = int(req.Msg.Offset)
		}
		limit := 100
		if req.Msg.Limit != 100 {
			limit = int(req.Msg.Limit)
		}
		for rank := offset; rank < cm.LenNodes(); rank++ {
			node := cm.NodeRank(rank)
			if node == nil {
				break
			}
			resp.Nodes = append(resp.Nodes, &apiv3.Node{
				NodeId:   uint32(node.ID),
				Word:     string(node.Word),
				Rate:     float32(node.Rate),
				NumEdges: 0, // TODO
			})
			if len(resp.Nodes) >= limit {
				break
			}
		}
		resp.Count = uint32(len(resp.Nodes))
		if resp.Num > uint32(offset)+resp.Count {
			resp.Next = uint32(offset) + resp.Count + 1
		}
		handleNodeRateResp(resp)
	}
	cancel(nil)

	err := context.Cause(ctx)
	return responseError(err)
}

func responseError(err error) error {
	if err == nil || err == context.Canceled {
		return nil
	}
	switch err.(type) {
	case ndl.NdlError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録データの取得に失敗しました。; %s", err.Error()),
		)
	case analyzer.AnalysisError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("議事録を形態素解析結果中にエラーが発生しました。; %s", err.Error()),
		)
	case matrix.MatrixError:
		return connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("共起関係の計算に失敗しました。; %s", err.Error()),
		)
	default:
		return connect.NewError(
			connect.CodeUnknown,
			fmt.Errorf("予期せぬエラーが発生しました。; %s", err.Error()),
		)
	}
}
