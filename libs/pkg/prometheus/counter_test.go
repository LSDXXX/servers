package prometheus

import (
	"testing"
)

func TestFlowReceiveCounterInc(t *testing.T) {
	tests := []struct {
		name     string
		flowID   string
		flowName string
	}{
		{"case-1", "1", "flow-1"},
		{"case-2", "2", "flow-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FlowReceiveCounterInc(tt.flowID, tt.flowName)
		})
	}
}

func TestFlowErrorCounterInc(t *testing.T) {
	tests := []struct {
		name     string
		flowID   string
		flowName string
	}{
		{"case-1", "1", "flow-1"},
		{"case-2", "2", "flow-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FlowErrorCounterInc(tt.flowID, tt.flowName)
		})
	}
}

func TestNodeReceiveCounterInc(t *testing.T) {
	tests := []struct {
		name     string
		nodeID   string
		nodeName string
	}{
		{"case-1", "1", "flow-1"},
		{"case-2", "2", "flow-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NodeReceiveCounterInc(tt.nodeID, tt.nodeName)
		})
	}
}

func TestNodeErrorCounterInc(t *testing.T) {
	tests := []struct {
		name     string
		nodeID   string
		nodeName string
	}{
		{"case-1", "1", "flow-1"},
		{"case-2", "2", "flow-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NodeErrorCounterInc(tt.nodeID, tt.nodeName)
		})
	}
}
