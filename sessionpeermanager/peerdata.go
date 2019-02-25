package sessionpeermanager

import (
	"fmt"
	"time"

	"github.com/ipfs/go-cid"
)

const (
	newLatencyWeight = 0.5
)

type peerData struct {
	hasLatency bool
	latency    time.Duration
	lt         *latencyTracker
}

func newPeerData() *peerData {
	return &peerData{
		hasLatency: false,
		lt:         newLatencyTracker(),
		latency:    0,
	}
}

func (pd *peerData) AdjustLatency(k cid.Cid, hasFallbackLatency bool, fallbackLatency time.Duration) {

	latency, hasLatency := pd.lt.RecordResponse(k)
	if !hasLatency {
		if hasFallbackLatency {
			fmt.Printf("Fallback latency used: %f\n", fallbackLatency.Seconds())
		}
		latency, hasLatency = fallbackLatency, hasFallbackLatency
	} else {
		fmt.Println("Had local latency")
	}
	if hasLatency {
		if pd.hasLatency {
			pd.latency = time.Duration(float64(pd.latency)*(1.0-newLatencyWeight) + float64(latency)*newLatencyWeight)
		} else {
			fmt.Printf("New Latency: %f\n", latency.Seconds())
			pd.latency = latency
			pd.hasLatency = true
		}
	}
}