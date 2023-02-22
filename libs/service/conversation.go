package service

import (
	"context"

	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/LSDXXX/libs/repo"
)

func init() {
	util.PanicWhenError(container.Singleton(NewConversation))
}

type Conversation struct {
	container repo.ConversationContainer `container:"type"`
	ctx       context.Context
}

func NewConversation() *Conversation {
	out := Conversation{}
	util.PanicWhenError(container.Fill(&out))

	return &out
}

func (c *Conversation) WithContext(ctx context.Context) *Conversation {
	out := *c
	out.ctx = ctx
	return &out
}

func (c *Conversation) GetParentId(convId string) (string, bool) {
	return c.container.GetParentId(convId)
}

func (c *Conversation) SetParentId(convId string, parentId string) {
	c.container.SetParentId(convId, parentId)
}
