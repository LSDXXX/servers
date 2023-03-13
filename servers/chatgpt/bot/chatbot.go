package bot

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/libs/service"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ConversationContainer Conversation

type Chatbot struct {
	auth     *service.OpenAIAuth
	ctx      context.Context
	email    string
	password string
	proxy    string
}

func NewChatbot(email, password, proxy string) *Chatbot {
	out := &Chatbot{
		auth:     service.NewOpenAIAuth(email, password, proxy),
		email:    email,
		password: password,
		proxy:    proxy,
	}
	err := out.auth.Login()
	if err != nil {
		panic("failed to login: " + err.Error())
	}
	return out
}

func (c *Chatbot) WithContext(ctx context.Context) *Chatbot {
	out := *c
	out.ctx = ctx
	return &out
}

func (c *Chatbot) doAsk(content, convId, preConvId string, retry int) (res *http.Response, err error) {
	if retry > 3 {
		return nil, errors.New("failed to ask")
	}
	headers := map[string][]string{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", c.auth.AccessToken())},
		"Accept":        {"text/event-stream"},
		"Connection":    {"close"},
		"Referer":       {"https://chat.openai.com/chat"},
		"Origin":        {"https://chat.openai.com"},
		// "User-Agent":                {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"},
		"X-OpenAI-Assistant-App-Id": {""},
	}

	data, _ := json.Marshal(askReq{
		Action: "next",
		Messages: []askMessage{
			{
				Id:   uuid.NewString(),
				Role: "user",
				Content: struct {
					ContentType string   "json:\"content_type\""
					Parts       []string "json:\"parts\""
				}{
					ContentType: "text",
					Parts:       []string{content},
				},
			},
		},
		ConversationId:  convId,
		ParentMessageId: preConvId,
		Model:           "text-davinci-002-render-sha",
	})

	req, err := http.NewRequest(http.MethodPost,
		"https://bypass.duti.tech/api/conversation", bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	req.Header = headers
	httpClient := http.Client{}
	res, err = httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do http request")
	}
	if res.StatusCode != 200 {
		data, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		log.WithContext(c.ctx).Errorf("res: %s", string(data))
		err := c.auth.Login()
		if err != nil {
			return nil, errors.Wrap(err, "login")
		}
		return c.doAsk(content, convId, preConvId, retry+1)
	}
	return res, nil
}

func (c *Chatbot) Ask(content, convId, preConvId string) (id, parent, response string, err error) {
	res, err := c.doAsk(content, convId, preConvId, 0)
	if err != nil {
		return "", "", "", err
	}
	data, _ := ioutil.ReadAll(res.Body)
	lines := strings.Split(string(data), "\n")
	fmt.Println(string(data), res.StatusCode)

	var message, conversationId, parentId string

	for _, line := range lines {
		line = strings.Trim(line, " ")
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		line = strings.Trim(strings.TrimPrefix(line, "data:"), " ")
		if line == "[DONE]" {
			break
		}
		var resData ResponseMessage
		err := json.Unmarshal([]byte(line), &resData)
		if err != nil {
			continue
		}
		message = resData.Message.Content.Parts[0]
		parentId = resData.Message.ID
		conversationId = resData.ConversationID
	}
	return conversationId, parentId, message, nil
}

func (c *Chatbot) AskStream(content, convId, preConvId string) (<-chan ResponseMessage, error) {
	res, err := c.doAsk(content, convId, preConvId, 0)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(res.Body)
	ch := make(chan ResponseMessage)
	go func() {
		defer close(ch)
		i := 0
		for {
			data, _, err := reader.ReadLine()
			i++
			if i <= 3 {
				continue
			}
			if err != nil {
				break
			}
			line := string(data)
			line = strings.Trim(line, " ")
			if !strings.HasPrefix(line, "data:") {
				continue
			}
			line = strings.Trim(strings.TrimPrefix(line, "data:"), " ")
			if line == "[DONE]" {
				break
			}
			var resData ResponseMessage
			err = json.Unmarshal([]byte(line), &resData)
			if err != nil {
				continue
			}
			ch <- resData
		}
	}()
	return ch, nil
}

type Conversation struct {
	conversations sync.Map
}

func (c *Conversation) Add(id, parent string) {
	c.conversations.Store(id, parent)
}

func (c *Conversation) Del(id string) {
	c.conversations.Delete(id)
}

func (c *Conversation) GetParentId(id string) string {
	v, ok := c.conversations.Load(id)
	if !ok {
		return ""
	}
	return v.(string)
}
