package vk

import (
	"encoding/json"
	"time"
)

type Timestamp struct {
	time.Time
}

func (p *Timestamp) UnmarshalJSON(bytes []byte) error {
	var raw int64
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}
	*&p.Time = time.Unix(raw, 0)
	return nil
}
