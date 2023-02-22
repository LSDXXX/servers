package wsmanager

import (
	"net/http"
	"sync"

	"github.com/LSDXXX/libs/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// WSManager manager
type WSManager struct {
	Group                   map[string]map[string]*WSClient
	groupCount, clientCount uint
	Lock                    sync.Mutex
	Register, UnRegister    chan *WSClient
	Message                 chan *MessageData
	GroupMessage            chan *GroupMessageData
	BroadCastMessage        chan *BroadCastMessageData
}

// WSEventHandler handler
type WSEventHandler interface {
	OnUpgrade(ctx *gin.Context, manager *WSManager) (groupID, clientID string, err error)
	OnClientRegister(*WSClient)
	OnClientDeregister(*WSClient)
	OnClientMessage(*WSClient, []byte) error
}

// WSClient client
type WSClient struct {
	Id, Group string
	Socket    *websocket.Conn
	Message   chan []byte
	Manager   *WSManager
	handler   WSEventHandler
}

// MessageData message info
type MessageData struct {
	Id, Group string
	Message   []byte
}

// GroupMessageData 组广播数据信息
type GroupMessageData struct {
	Group   string
	Message []byte
}

// BroadCastMessageData broad cast message data
type BroadCastMessageData struct {
	Message []byte
}

// Read description
// @receiver c
func (c *WSClient) Read() {
	defer func() {
		c.Manager.UnRegister <- c
		logrus.Infof("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			logrus.Infof("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		logrus.Infof("client [%s] receive message: %s", c.Id, string(message))
		c.handler.OnClientMessage(c, message)
	}
}

// Write description
// @receiver c
func (c *WSClient) Write() {
	defer func() {
		logrus.Infof("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			logrus.Infof("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			logrus.Infof("client [%s] write message: %s", c.Id, string(message))
			err := c.Socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				logrus.Infof("client [%s] writemessage err: %s", c.Id, err)
			}
		}
	}
}

// 启动 websocket 管理器
func (manager *WSManager) start() {
	logrus.Infof("websocket manage start")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			logrus.Infof("client [%s] connect", client.Id)
			logrus.Infof("register client [%s] to group [%s]", client.Id, client.Group)
			client.handler.OnClientRegister(client)

			manager.Lock.Lock()
			if manager.Group[client.Group] == nil {
				manager.Group[client.Group] = make(map[string]*WSClient)
				manager.groupCount += 1
			}
			manager.Group[client.Group][client.Id] = client
			manager.clientCount += 1
			manager.Lock.Unlock()

		// 注销
		case client := <-manager.UnRegister:
			logrus.Infof("unregister client [%s] from group [%s]", client.Id, client.Group)
			client.handler.OnClientDeregister(client)
			manager.Lock.Lock()
			if _, ok := manager.Group[client.Group]; ok {
				if _, ok := manager.Group[client.Group][client.Id]; ok {
					close(client.Message)
					delete(manager.Group[client.Group], client.Id)
					manager.clientCount -= 1
					if len(manager.Group[client.Group]) == 0 {
						//log.Printf("delete empty group [%s]", client.Group)
						delete(manager.Group, client.Group)
						manager.groupCount -= 1
					}
				}
			}
			manager.Lock.Unlock()

			// 发送广播数据到某个组的 channel 变量 Send 中
			//case data := <-manager.boardCast:
			//	if groupMap, ok := manager.wsGroup[data.GroupId]; ok {
			//		for _, conn := range groupMap {
			//			conn.Send <- data.Data
			//		}
			//	}
		}
	}
}

// 处理单个 client 发送数据
func (manager *WSManager) sendService() {
	for {
		select {
		case data := <-manager.Message:
			manager.Lock.Lock()
			if groupMap, ok := manager.Group[data.Group]; ok {
				if conn, ok := groupMap[data.Id]; ok {
					conn.Message <- data.Message
				}
			}
			manager.Lock.Unlock()
		}
	}
}

// 处理 group 广播数据
func (manager *WSManager) sendGroupService() {
	for {
		select {
		// 发送广播数据到某个组的 channel 变量 Send 中
		case data := <-manager.GroupMessage:
			manager.Lock.Lock()
			if groupMap, ok := manager.Group[data.Group]; ok {
				for _, conn := range groupMap {
					conn.Message <- data.Message
				}
			}
			manager.Lock.Unlock()
		}
	}
}

// 处理广播数据
func (manager *WSManager) sendAllService() {
	for {
		select {
		case data := <-manager.BroadCastMessage:
			manager.Lock.Lock()
			for _, v := range manager.Group {
				for _, conn := range v {
					conn.Message <- data.Message
				}
			}
			manager.Lock.Unlock()
		}
	}
}

// Send send
//  @receiver manager
//  @param id
//  @param group
//  @param message
func (manager *WSManager) Send(id string, group string, message []byte) {
	data := &MessageData{
		Id:      id,
		Group:   group,
		Message: message,
	}
	manager.Message <- data
}

// SendGroup group send
//  @receiver manager
//  @param group
//  @param message
func (manager *WSManager) SendGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group:   group,
		Message: message,
	}
	manager.GroupMessage <- data
}

// SendAll 广播
//  @receiver manager
//  @param message
func (manager *WSManager) SendAll(message []byte) {
	data := &BroadCastMessageData{
		Message: message,
	}
	manager.BroadCastMessage <- data
}

// RegisterClient 注册
//  @receiver manager
//  @param client
func (manager *WSManager) RegisterClient(client *WSClient) {
	manager.Register <- client
}

// UnRegisterClient 注销
//  @receiver manager
//  @param client
func (manager *WSManager) UnRegisterClient(client *WSClient) {
	manager.UnRegister <- client
}

// LenGroup num groups
//  @receiver manager
//  @return uint
func (manager *WSManager) LenGroup() uint {
	return manager.groupCount
}

// LenClient num clients
//  @receiver manager
//  @return uint
func (manager *WSManager) LenClient() uint {
	return manager.clientCount
}

// Info manager info
//  @receiver manager
//  @return map
func (manager *WSManager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}

// NewWSManager new
//  @return *WSManager
func New() *WSManager {
	manager := &WSManager{
		Group:            make(map[string]map[string]*WSClient),
		Register:         make(chan *WSClient, 128),
		UnRegister:       make(chan *WSClient, 128),
		GroupMessage:     make(chan *GroupMessageData, 128),
		Message:          make(chan *MessageData, 128),
		BroadCastMessage: make(chan *BroadCastMessageData, 128),
		groupCount:       0,
		clientCount:      0,
	}
	go manager.start()
	go manager.sendGroupService()
	go manager.sendService()
	go manager.sendAllService()
	return manager
}

// BuildHTTPHandler  build http handler
//  @receiver manager
//  @param handler
//  @return *gin.Context
//  @return func(*gin.Context)
func (manager *WSManager) BuildHTTPHandler(handler WSEventHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		groupID, clientID, err := handler.OnUpgrade(c, manager)
		if err != nil {
			log.WithContext(ctx).Errorf("handle upgrade event error: %+v", err)
			return
		}
		if len(clientID) == 0 {
			clientID = uuid.NewString()
		}
		upGrader := websocket.Upgrader{
			// cross origin domain
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			// 处理 Sec-WebSocket-Protocol Header
			Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
		}

		conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logrus.Infof("websocket connect error: %s", c.Param("channel"))
			return
		}

		client := &WSClient{
			Id:      clientID,
			Group:   groupID,
			Socket:  conn,
			Manager: manager,
			Message: make(chan []byte, 1024),
			handler: handler,
		}

		manager.RegisterClient(client)
		go client.Read()
		go client.Write()
	}
}
