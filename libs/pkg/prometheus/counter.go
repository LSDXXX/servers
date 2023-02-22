package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// label 定义
const (
	labelFlowName = "flowName"
	labelFlowID   = "flowID"
	labelNodeName = "nodeName"
	labelNodeID   = "nodeID"
)

var (
	// Flow 接收事件计数
	flowReceiveCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "flow_receive_total",
		Help: "The number of received events by the flow (per event type).",
	}, []string{labelFlowName, labelFlowID})

	// Flow 错误事件计数
	flowErrorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "flow_error_total",
		Help: "The number of error events by the flow (per event type).",
	}, []string{labelFlowName, labelFlowID})

	// Node 接收事件计数
	nodeReceiveCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "node_receive_total",
		Help: "The number of received events by the node (per event type).",
	}, []string{labelNodeName, labelNodeID})

	// Node 错误事件计数
	nodeErrorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "node_error_total",
		Help: "The number of error events by the node (per event type).",
	}, []string{labelNodeName, labelNodeID})
)

// Increase for inc flow receive counter value
func FlowReceiveCounterInc(flowID, flowName string) {
	tagSet := map[string]string{}
	tagSet[labelFlowID] = flowID
	tagSet[labelFlowName] = flowName
	flowReceiveCounter.With(tagSet).Inc()
}

// Add for add flow receive counter value
func FlowReceiveCounterAdd(flowID, flowName string, cnt int64) {
	tagSet := map[string]string{}
	tagSet[labelFlowID] = flowID
	tagSet[labelFlowName] = flowName
	flowReceiveCounter.With(tagSet).Add(float64(cnt))
}

// Reset for reset flow receive counter value
func FlowReceiveCounterReset() {
	flowReceiveCounter.Reset()
}

// Increase for inc flow error counter value
func FlowErrorCounterInc(flowID, flowName string) {
	tagSet := map[string]string{}
	tagSet[labelFlowID] = flowID
	tagSet[labelFlowName] = flowName
	flowErrorCounter.With(tagSet).Inc()
}

// Reset for reset flow error counter value
func FlowErrorCounterReset() {
	flowErrorCounter.Reset()
}

// Increase for inc node receive counter value
func NodeReceiveCounterInc(nodeID, nodeName string) {
	tagSet := map[string]string{}
	tagSet[labelNodeID] = nodeID
	tagSet[labelNodeName] = nodeName
	nodeReceiveCounter.With(tagSet).Inc()
}

// Reset for reset node receive counter value
func NodeReceiveCounterReset() {
	nodeReceiveCounter.Reset()
}

// Increase for inc node error counter value
func NodeErrorCounterInc(nodeID, nodeName string) {
	tagSet := map[string]string{}
	tagSet[labelNodeID] = nodeID
	tagSet[labelNodeName] = nodeName
	nodeErrorCounter.With(tagSet).Inc()
}

// Reset for reset node error counter value
func NodeErrorCounterReset() {
	nodeErrorCounter.Reset()
}
