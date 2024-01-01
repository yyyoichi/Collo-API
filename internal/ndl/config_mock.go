package ndl

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateMeetingConfigMock(config Config, dir string) Config {
	if dir == "" {
		dir = "/tmp/test-meeting-result"
	}
	if _, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, 0700); err != nil {
			panic(err)
		}
	}

	config.init()
	DoGet := config.DoGet
	storeDoGet := func(url string) (body []byte, err error) {
		file := filepath.Join(fmt.Sprintf("%s/%s.json", dir, md5Hash(url)))
		if b, err := os.ReadFile(file); err == nil || errors.Is(err, io.EOF) {
			return b, err
		}
		b, err := DoGet(url)
		if err != nil {
			panic(err)
		}
		f, err := os.OpenFile(
			file,
			os.O_RDWR|os.O_CREATE,
			0600,
		)
		if err != nil {
			panic(err)
		}

		if _, err := f.Write(b); err != nil {
			panic(err)
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
		return b, nil
	}
	config.DoGet = storeDoGet
	client := NewClient(config)
	_, resultCh := client.GenerateResult(context.Background())
	for range resultCh {
	}
	return config
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
