package handler

import (
	"context"
	"log"
	apiv3 "yyyoichi/Collo-API/internal/api/v3"
)

type (
	storage           struct{}
	storagePermission struct {
		useStorage  bool
		saveStorage bool
	}
)

func (p *storagePermission) NewCoMatrixes(ctx context.Context, processHandler ProcessHandler, config Config) CoMatrixes {
	var coMatrixes CoMatrixes
	var usedInStorage bool
	if p.useStorage {
		// get from config.ToString()
		// if found foundInStorage = true
		// else coMatrixes = NewCoMatrixes(ctx, processHandler, config)
	} else {
		coMatrixes = NewCoMatrixes(ctx, processHandler, config)
	}

	// savaする指定がありかつ、ストレージからデータを取得していなければ、ストレージを保存する
	if p.saveStorage && !usedInStorage {
		// save coMatrixes in /tmp
		log.Printf("save at /tmp/%s", config.ToString())
	}

	return coMatrixes
}

func (storage) PermitNetworkStreamRequest(req *apiv3.NetworkStreamRequest) storagePermission {
	var p storagePermission
	if req.ForcusNodeId == 0 {
		// ノードをフォーカスするなら初回リクエストではなく
		// おそらくストレージが存在するので、ストレージを使ってほしい
		p.useStorage = true
		p.saveStorage = false
		return p
	}
	if req.Config.ForcusGroupId != "" {
		// 特定のグループにフォーカスするなら初回リクエストではなく
		// おそらくストレージが存在するので、ストレージを使ってほしい
		p.useStorage = true
		p.saveStorage = false
		return p
	}

	if !req.Config.UseNdlCache {
		// APIのキャッシュを利用せず最新の内容を取得したいのであれば、
		// ストレージは使用せず、結果を保存してほしい
		p.useStorage = false
		p.saveStorage = true
	}

	// 初回リクエストっぽいので、
	// ストレージがあれば使用してほしいし
	// 結果を保存しておいてほしい
	p.useStorage = true
	p.saveStorage = true
	return p
}

func (storage) PermitNodeRateStreamRequest(req *apiv3.NodeRateStreamRequest) storagePermission {
	var p storagePermission
	if req.Config.ForcusGroupId != "" {
		// 特定のグループにフォーカスするなら初回リクエストではなく
		// おそらくストレージが存在するので、ストレージを使ってほしい
		p.useStorage = true
		p.saveStorage = false
		return p
	}

	if !req.Config.UseNdlCache {
		// APIのキャッシュを利用せず最新の内容を取得したいのであれば、
		// ストレージは使用せず、結果を保存してほしい
		p.useStorage = false
		p.saveStorage = true
	}

	// 初回リクエストっぽいので、
	// ストレージがあれば使用してほしいし
	// 結果を保存しておいてほしい
	p.useStorage = true
	p.saveStorage = true
	return p
}
