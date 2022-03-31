package vk_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kandayo/vkscraper/pkg/vk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimestampUnmarshalJSON(t *testing.T) {
	t.Parallel()
	unixTimestamp := 1648556794
	data, _ := json.Marshal(unixTimestamp)
	var timestamp vk.Timestamp
	err := json.Unmarshal(data, &timestamp)
	require.Nil(t, err)
	assert.Equal(t, vk.Timestamp{time.Unix(1648556794, 0)}, timestamp)
}
