package network

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var NManager = &NetworkManager{
	data: make(map[string]struct {
		*Network
		expiration time.Time
	}),
	dir:      "/tmp",
	hFn:      md5Hash,
	mu:       sync.Mutex{},
	ttl:      time.Duration(time.Minute * 2),
	stopChan: make(chan struct{}),
}

type NetworkManager struct {
	data map[string]struct {
		*Network
		expiration time.Time
	}
	dir      string
	hFn      func(string) string
	mu       sync.Mutex
	ttl      time.Duration
	stopChan chan struct{}
}

func (m *NetworkManager) Set(key string, network *Network) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key = m.hFn(key)
	m.data[key] = struct {
		*Network
		expiration time.Time
	}{
		Network:    network,
		expiration: time.Now().Add(m.ttl),
	}

	bytes, err := json.Marshal(network)
	if err != nil {
		return err
	}

	// 永続化
	file, err := os.OpenFile(
		filepath.Join(m.dir, fmt.Sprintf("%s%s", key, ".json")),
		os.O_RDWR|os.O_CREATE,
		0600,
	)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteAt(bytes, 0)

	return err
}

func (m *NetworkManager) Get(key string) (*Network, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key = m.hFn(key)
	if data, found := m.data[key]; found {
		return data.Network, true
	}

	bytes, err := os.ReadFile(
		filepath.Join(m.dir, fmt.Sprintf("%s%s", key, ".json")),
	)
	if err != nil {
		return nil, false
	}

	network := NewNetwork()
	if err := json.Unmarshal(bytes, &network); err != nil {
		return nil, false
	}
	network.refreshMap()
	return network, true
}

func (m *NetworkManager) StartCleanup() {
	go func() {
		ticker := time.NewTicker(time.Duration(time.Second))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.cleanupExpiredData()
			case <-m.stopChan:
				return
			}
		}
	}()
}

func (m *NetworkManager) cleanupExpiredData() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, data := range m.data {
		if time.Now().After(data.expiration) {
			delete(m.data, key)
		}
	}
}

func (m *NetworkManager) StopCleanup() {
	close(m.stopChan)
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
