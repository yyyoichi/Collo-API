package matrix

import (
	"sort"
	"time"
)

type (
	MultiDocMeta struct {
		GroupID string         `json:"i"` // グループ識別子
		From    time.Time      `json:"f"` // 開始日
		Until   time.Time      `json:"u"` // 終了日
		Metas   []DocumentMeta `json:"m"`
	}

	DocumentMeta struct {
		Key         string    `json:"k"` // 識別子
		Name        string    `json:"n"` // 任意の名前
		At          time.Time `json:"a"` // 日付
		Description string    `json:"d"` // 説明
	}
)

func NewMultiDocMeta(id string, metas []DocumentMeta) MultiDocMeta {
	var meta MultiDocMeta
	meta.GroupID = id
	if len(metas) == 0 {
		return meta
	}
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].At.Before(metas[j].At)
	})
	meta.From = metas[0].At
	meta.Until = metas[len(metas)-1].At
	meta.Metas = metas
	return meta
}
