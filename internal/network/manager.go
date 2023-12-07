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
	data: make(map[string]*struct {
		*Network
		expiration time.Time
	}),
	dir:      "/tmp/collo-network",
	hFn:      md5Hash,
	mu:       sync.Mutex{},
	ttl:      time.Duration(time.Minute * 2),
	tick:     time.Duration(time.Second),
	stopChan: make(chan struct{}),
}

type NetworkManager struct {
	data map[string]*struct {
		*Network
		expiration time.Time
	}
	dir      string
	hFn      func(string) string
	mu       sync.Mutex
	ttl      time.Duration
	tick     time.Duration
	stopChan chan struct{}
}

func (m *NetworkManager) Set(key string, network *Network) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key = m.hFn(key)
	m.data[key] = &struct {
		*Network
		expiration time.Time
	}{
		Network:    network,
		expiration: time.Now().Add(m.ttl),
	}

	err := m.setToStrage(key, network)
	return err
}

func (m *NetworkManager) Get(key string) (*Network, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key = m.hFn(key)
	if data, found := m.data[key]; found {
		data.expiration = time.Now().Add(m.ttl)
		return data.Network, true
	}

	if network := m.getFromStrage(key); network != nil {
		return network, true
	} else {
		return nil, false
	}
}

func (m *NetworkManager) setToStrage(key string, network *Network) error {
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

func (m *NetworkManager) getFromStrage(key string) *Network {
	bytes, err := os.ReadFile(
		filepath.Join(m.dir, fmt.Sprintf("%s%s", key, ".json")),
	)
	if err != nil {
		return nil
	}

	network := NewNetwork()
	if err := json.Unmarshal(bytes, &network); err != nil {
		return nil
	}
	network.refreshMap()
	m.data[key] = &struct {
		*Network
		expiration time.Time
	}{
		Network:    network,
		expiration: time.Now().Add(m.ttl),
	}
	return network
}

func (m *NetworkManager) StartCleanup() {
	go func() {
		ticker := time.NewTicker(m.tick)
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
