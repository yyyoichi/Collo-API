package datalake

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	var config Config
	config.Search.Any = "能登半島地震"
	config.Search.From = parseTime("2024-01-31")
	config.Search.Until = parseTime("2024-01-31")
	config.CacheDir = "/tmp/datalake-test"
	config.NoCache = false
	client := NewClient(config)
	ctx := context.Background()
	for m, err := range client.IterMeeting(ctx) {
		assert.NoError(t, err)
		assert.NotNil(t, m)
	}
}
