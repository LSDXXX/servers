package constant

import (
	"time"

	"github.com/LSDXXX/libs/pkg/crudgen/helper"
)

const (
	// TraceIDHTTPHeader .
	TraceIDHTTPHeader = "X-B3-TraceId"

	// ProjectIDHTTPHeader .
	ProjectIDHTTPHeader = "x-welink-project-id"

	// AppIDHTTPHeader .
	//AppIDHTTPHeader = "x-welink-app-id"

	// SpanIDHTTPHeader .
	SpanIDHTTPHeader = "X-B3-SpanId"

	// RequestIDHTTPHeader .
	RequestIDHTTPHeader = "x-request-id"

	// TokenHTTPHeader .
	TokenHTTPHeader = "x-welink-token-id"

	// DisablePushLogHeader
	DisablePushLogHeader = "x-disable-push-log"

	// RedisPrefixKey .
	RedisPrefixKey = "logic_engine"

	// RPCTimeout .
	RPCTimeout = time.Second * 5

	// ZKLockPath .
	ZKLockPath = "/logic-engine/locks"

	OnlineCalculate  = 0
	OfflineCalculate = 1

	SubFlowKeyPrefix = "subFlow-"

	NodeCatalogOrgId = "node"

	SceneTemplateCatalogOrgId = "scene_template"

	CatalogNotBuiltIn = 2

	Mysql = "mysql"

	Pgsql = "pgsql"

	Tdsql = "tdsql"

	SceneSnapshot = 0

	QuotaSnapshot = 1

	SystemMaxId = helper.SystemMaxId

	CatalogMaxLevel = 3

	DeployOrPublishStatus = 1

	DeployOrPublishStatusFail = 2

	RunStatus = 1

	RunStatusFail = 2

	DefaultPageNum = 1

	DefaultPageSize = 10

	NotExist = -1

	JWTIdentityKey = "jwtIdentity"
)
