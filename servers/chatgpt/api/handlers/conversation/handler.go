package conversation

import (
	"github.com/LSDXXX/libs/pkg/handlergen/helper"
	"github.com/LSDXXX/libs/service"
	"github.com/LSDXXX/servers/chatgpt/bot"
)

//@RequestMapping(/api)
//go:generate handlergentool -f $GOFILE -op ./ -pkg $GOPACKAGE
type ConversationHandler interface {
	helper.InjectServices2[*bot.Chatbot, *service.Conversation]

	//@RequestMapping(/ask, POST)
	//@BindBody(req)
	Ask(req AskReq) (string, error)
}
