package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	apiv3 "yyyoichi/Collo-API/internal/api/v3"
)

type (
	Storage struct {
		CoMatrixes  CoMatrixes `json:"m"`
		Words       []string   `json:"w"`
		new         func(context.Context, ProcessHandler, Config) CoMatrixes
		getFilename func(string) string
	}
	storagePermission struct {
		useStorage  bool
		saveStorage bool
		Config
	}
)

func (s *Storage) NewCoMatrixes(ctx context.Context, processHandler ProcessHandler, config storagePermission) CoMatrixes {
	s.init()
	var usedInStorage bool
	if config.useStorage {
		usedInStorage = s.readCoMatrixes(ctx, processHandler, config.Config)
	} else {
		s.CoMatrixes = s.new(ctx, processHandler, config.Config)
	}
	if len(s.CoMatrixes) > 0 {
		s.Words = s.CoMatrixes[0].Words
	}

	// savaする指定がありかつ、ストレージからデータを取得していなければ、ストレージを保存する
	if config.saveStorage && !usedInStorage {
		// save coMatrixes in /tmp
		if err := s.saveCoMatrixes(ctx, config.Config); err != nil {
			log.Println(err)
		}
	}

	return s.CoMatrixes
}

// Strageから読み込みを試みます。成功した場合trueを返します。
func (s *Storage) readCoMatrixes(ctx context.Context, processHandler ProcessHandler, config Config) bool {
	filename := s.getFilename(config.ToString())
	f, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	if err != nil && err != io.EOF {
		s.CoMatrixes = s.new(ctx, processHandler, config)
		return false
	}

	var b bytes.Buffer
	if _, err := io.Copy(&b, f); err != nil {
		return false
	}
	if err := json.Unmarshal(b.Bytes(), s); err != nil {
		s.CoMatrixes = s.new(ctx, processHandler, config)
		return false
	}
	// Wordは各行列に保存されていないのでセット
	for _, cm := range s.CoMatrixes {
		cm.Words = s.Words
	}
	slog.InfoContext(ctx, "read from storage", slog.String("filename", filename))
	return true
}
func (s *Storage) saveCoMatrixes(ctx context.Context, config Config) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	filename := s.getFilename(config.ToString())
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, bytes.NewReader(b)); err != nil {
		return err
	}
	slog.InfoContext(ctx, "save in storage", slog.String("filename", filename))
	return nil
}

func (Storage) PermitNetworkStreamRequest(req *apiv3.NetworkStreamRequest) storagePermission {
	var p storagePermission
	if req.ForcusNodeId != 0 {
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

func (Storage) PermitNodeRateStreamRequest(req *apiv3.NodeRateStreamRequest) storagePermission {
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

func (s *Storage) init() {
	if s.getFilename == nil {
		s.getFilename = func(s string) string {
			hasher := sha256.New()
			hasher.Write([]byte(s))
			hashBytes := hasher.Sum(nil)
			// バイト列を16進数文字列に変換
			return "/tmp/" + hex.EncodeToString(hashBytes)
		}
	}
	if s.new == nil {
		s.new = NewCoMatrixes
	}
}
