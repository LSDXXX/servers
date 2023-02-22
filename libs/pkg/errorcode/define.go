package errorcode

import "errors"

var (
	defines map[error]int

	// ErrUnknownError .
	ErrUnknownError = errors.New("未知参数")

	// ErrParameterInvalid .
	ErrParameterInvalid = errors.New("参数不合理")

	// ErrNotRegistered .
	ErrNotRegistered = errors.New("未注册用户")

	// ErrIncompleteParameters .
	ErrIncompleteParameters = errors.New("参数缺失")

	// ErrServiceUnavailable .
	ErrServiceUnavailable = errors.New("服务不可用")

	// ErrNotFound .
	ErrNotFound = errors.New("资源不存在")

	// ErrInternalServerError .
	ErrInternalServerError = errors.New("内部服务器错误或异常")

	// ErrNotImplemented .
	ErrNotImplemented = errors.New("功能未实现")

	// ErrMissingNode worker自己error
	ErrMissingNode = errors.New("节点缺失")

	// ErrUpdateFlowError .
	ErrUpdateFlowError = errors.New("部署flow出错")

	// ErrFlowIsStarted .
	ErrFlowIsStarted = errors.New("flow已经部署")

	// ErrFlowIsStoped .
	ErrFlowIsStoped = errors.New("flow已经停止")

	// ErrInvalidNodeConfig .
	ErrInvalidNodeConfig = errors.New("节点配置信息错误")

	// ErrHTTPInNodeExist .
	ErrHTTPInNodeExist = errors.New("httpin 相同uri的节点已经存在, 一个项目只能存在一个")

	// ErrMSGPayloadCastError .
	ErrMSGPayloadCastError = errors.New("消息体类型转换失败")

	// ErrHTTPRequestTimeOut .
	ErrHTTPRequestTimeOut = errors.New("请求超时")

	// ErrInvalidFlow .
	ErrInvalidFlow = errors.New("flow编排错误")

	// ErrFlowNeedsToBeDeployed .
	ErrFlowNeedsToBeDeployed = errors.New("flow需要先部署")

	// ErrMissingDBConnector .
	ErrMissingDBConnector = errors.New("数据库连接节点缺失")

	ErrInvalidConnector = errors.New("获取连接器失败")

	ErrBindTimer = errors.New("绑定定时器失败")

	ErrFlowIsDeployed = errors.New("已部署最新版本")

	ErrFlowIsNotDeployed = errors.New("流程未部署")

	ErrInvalidSignature = errors.New("sig 校验失败")

	ErrNoWorker = errors.New("no worker")

	ErrBlankFlowTemplate = errors.New("不能发布空流模版")

	ErrSystemResourceCanNotChange = errors.New("内置资源不能更改")

	NameHasBeenUsed = errors.New("名称已被使用")

	MissingParameters = errors.New("缺少必填参数")

	NameCanNotBeEmpty = "名称不能为空"

	InvalidRequestParameters = "请求参数异常"

	InternalError = "内部异常，请重试"

	TemplateNotPublish = "模版未发布"

	ErrCatalogNotExist = "目录不存在"

	CatalogBuiltInDeleteError = "内置目录不允许删除"

	CatalogBuiltInUpdateError = "内置目录不允许编辑"

	ErrCatalogLevelExceed = errors.New("目录层级超出最大限制")

	// CoordinatorRequest .
	CoordinatorRequest = "CoordinatorRequest"

	// HandleCoordinatorError .
	HandleCoordinatorError = "HandleCoordinatorError"

	// HTTPInRequest .
	HTTPInRequest = "HTTPInRequest"

	// HandleHTTPInError .
	HandleHTTPInError = "HandleHTTPInError"

	// NodeReceive .
	NodeReceive = "NodeReceive"

	// SendAlarm .
	SendAlarm = "SendAlarm"

	// SendToApp .
	SendToApp = "SendToApp"

	// ReceiveAppMessage .
	ReceiveAppMessage = "ReceiveAppMessage"

	// HandleAppMessage .
	HandleAppMessage = "HandleAppMessage"

	// SendToDevice .
	SendToDevice = "SendToDevice"

	// ReceiveDeviceMessage .
	ReceiveDeviceMessage = "ReceiveDeviceMessage"

	// HandleDeviceMessage .
	HandleDeviceMessage = "HandleDeviceMessage"
)

func init() {
	defines = map[error]int{
		ErrUnknownError:         1102000,
		ErrParameterInvalid:     1102025,
		ErrNotRegistered:        1102026,
		ErrIncompleteParameters: 1102028,
		ErrServiceUnavailable:   1102036,
		ErrNotFound:             1102032,
		ErrInternalServerError:  1102034,
		ErrNotImplemented:       1102035,
		ErrNoWorker:             1102036,
		ErrInvalidSignature:     1102037,

		ErrMissingDBConnector:    1200001,
		ErrFlowNeedsToBeDeployed: 1200002,

		ErrFlowIsDeployed:    30001,
		ErrFlowIsNotDeployed: 30003,
		ErrInvalidNodeConfig: 30004,
	}
}

// Code get error code
func Code(err error) int {
	tmp := err
	// extract original error
	for {
		tmp = errors.Unwrap(tmp)
		if tmp == nil {
			break
		}
		err = tmp
	}
	if code, ok := defines[err]; ok {
		return code
	}
	return -1
}
