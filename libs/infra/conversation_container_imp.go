package infra

import (
	"sync"

	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/repo"
	"github.com/spf13/cast"
)

func init() {
	AppendInitFunc(func() {
		_ = container.Singleton(NewConversationHandlerImp)
	})
}

type ConversationContainerImp struct {
	conversations sync.Map
}

func NewConversationHandlerImp() repo.ConversationContainer {
	return &ConversationContainerImp{}
}

func (c *ConversationContainerImp) SetParentId(id, parent string) {
	c.conversations.Store(id, parent)
}

func (c *ConversationContainerImp) Del(id string) {
	c.conversations.Delete(id)
}

func (c *ConversationContainerImp) GetParentId(id string) (string, bool) {
	v, ok := c.conversations.Load(id)
	return cast.ToString(v), ok
}
