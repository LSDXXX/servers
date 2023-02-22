package conversation

import "github.com/google/uuid"

func (imp *ConversationHandlerImp) Ask(req AskReq) (string, error) {
	//TODO:

	parent, ok := imp.Conversation.GetParentId(req.ConversationId)
	if !ok {
		parent = uuid.NewString()
	}
	id, parentId, res, err := imp.Chatbot.Ask(req.Content, req.ConversationId, parent)
	if err != nil {
		return "", err
	}
	imp.Conversation.SetParentId(id, parentId)
	return res, nil
}
