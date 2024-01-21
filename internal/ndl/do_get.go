package ndl

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type cachedDoGet struct {
	useCache    bool // 検索時キャッシュ利用
	createCache bool // 検索後キャッシュ作成
	dir         string
}

func (d *cachedDoGet) DoGet(url string) ([]byte, error) {
	filename := d.filename(url)
	if d.useCache {
		if b, err := os.ReadFile(filename); err == nil || errors.Is(err, io.EOF) {
			return b, err
		}
	}
	b, err := d.doget(url)
	if err != nil {
		return nil, err
	}

	if !d.createCache {
		return b, nil
	}

	// create cache
	f, err := os.OpenFile(
		filename,
		os.O_RDWR|os.O_CREATE,
		0600,
	)
	if err != nil {
		return b, nil
	}

	if _, err := f.Write(b); err != nil {
		return b, nil
	}
	if err := f.Close(); err != nil {
		return b, nil
	}
	return b, nil
}

func (d *cachedDoGet) doget(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (d *cachedDoGet) filename(url string) string {
	return filepath.Join(fmt.Sprintf("%s/%s.json", d.dir, md5Hash(url)))
}

func md5Hash(input string) string {
	// MD5ハッシュ関数を作成
	hasher := md5.New()

	hasher.Write([]byte(input))

	// ハッシュを取得し、16進数文字列に変換
	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString
}
