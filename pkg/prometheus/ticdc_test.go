package prometheus

import (
	"log"
	"testing"
)

const (
	addr = "http://172.16.6.191:9090"
)

func TestTiCDCUptime(t *testing.T) {
	metric, err := NewTiCDCUptime(addr)
	if err != nil {
		log.Panicf("create ticdc uptime metric failed, err: %+v", err)
	}

	resp, err := metric.Get()
	if err != nil {
		log.Panicf("get ticdc uptime metric failed, err: %+v", err)
	}

	log.Printf("resp = %+v", resp)
}
