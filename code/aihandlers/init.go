package aihandlers

import (
	"context"
	"fmt"
	"start-feishubot/initialization"
	"start-feishubot/logger"
	"start-feishubot/services/openai"
	"start-feishubot/utils"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type MessageHandlerInterface interface {
	msgReceivedHandler(ctx context.Context, event *larkim.P2MessageReceiveV1) error
	CardHandler(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error)
	Handler(ctx context.Context, event *larkim.P2MessageReceiveV1) error
}

type HandlerType string

var sdhandlers utils.SafeDict

func init() {
	sdhandlers = *utils.NewSafeDict(make(map[string]interface{}))
}

const (
	GroupHandler      = "group"
	UserHandler       = "personal"
	CacheAi302Handler = "CACHE_AI302_HANDLER"
)

func InitHandlers(gpt *openai.ChatGPT, config initialization.Config, tokenMappingID int, botName string) MessageHandlerInterface {
	key := fmt.Sprintf("%s_%d", CacheAi302Handler, tokenMappingID)
	sdhandlers.Get(key)
	if sdhandlers.Exists(key) {
		result, _ := sdhandlers.Get(key)
		return result.(MessageHandlerInterface)
	} else {
		result := NewMessageHandler(gpt, config, botName)
		sdhandlers.Put(key, result)
		return result
	}
}

// func  Handler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
// 	return msgReceivedHandler(ctx, event)
// }

func ReadHandler(ctx context.Context, event *larkim.P2MessageReadV1) error {
	readerId := event.Event.Reader.ReaderId.OpenId
	//fmt.Printf("msg is read by : %v \n", *readerId)
	logger.Debugf("msg is read by : %v \n", *readerId)

	return nil
}

/*
func CardHandler() func(
	ctx context.Context,
	cardAction *larkcard.CardAction) (interface{}, error) {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		//handlerType := judgeCardType(cardAction)
		return handlers.cardHandler(ctx, cardAction)
	}
}
*/

func judgeCardType(cardAction *larkcard.CardAction) HandlerType {
	actionValue := cardAction.Action.Value
	chatType := actionValue["chatType"]
	//fmt.Printf("chatType: %v", chatType)
	if chatType == "group" {
		return GroupHandler
	}
	if chatType == "personal" {
		return UserHandler
	}
	return "otherChat"
}

func judgeChatType(event *larkim.P2MessageReceiveV1) HandlerType {
	chatType := event.Event.Message.ChatType
	if *chatType == "group" {
		return GroupHandler
	}
	if *chatType == "p2p" {
		return UserHandler
	}
	return "otherChat"
}
