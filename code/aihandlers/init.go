package aihandlers

import (
	"context"
	"fmt"
	"start-feishubot/initialization"
	"start-feishubot/logger"
	"start-feishubot/services/openai"
	"start-feishubot/utils"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type MessageHandlerInterface interface {
	msgReceivedHandler(ctx context.Context, event *larkim.P2MessageReceiveV1) error
	cardHandler(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error)
	Handler(ctx context.Context, event *larkim.P2MessageReceiveV1) error
}

type HandlerType string

var sdhandlers utils.SafeDict
var sdLarkClient utils.SafeDict

func init() {
	sdhandlers = *utils.NewSafeDict(make(map[string]interface{}))
	sdLarkClient = *utils.NewSafeDict(make(map[string]interface{}))
}

const (
	GroupHandler         = "group"
	UserHandler          = "personal"
	CacheAi302Handler    = "CACHE_AI302_HANDLER"
	CacheAi302LarkClient = "CACHE_AI302_LARKCLIENT"
)

func InitHandlers(gpt *openai.ChatGPT, config initialization.Config, tokenMappingID int, botName string, feishuAppID, feishuAppSecret string) MessageHandlerInterface {
	key := fmt.Sprintf("%s_%d", CacheAi302Handler, tokenMappingID)
	if sdhandlers.Exists(key) {
		result, _ := sdhandlers.Get(key)
		return result.(MessageHandlerInterface)
	} else {
		result := NewMessageHandler(gpt, config, botName, tokenMappingID)
		sdhandlers.Put(key, result)
		clientKey := fmt.Sprintf("%s_%d", CacheAi302LarkClient, tokenMappingID)
		larkClient := newLarkClient(config, feishuAppID, feishuAppSecret)
		sdLarkClient.Put(clientKey, larkClient)
		return result
	}
}

func GetAi302Handler(tokenMappingID int) *MessageHandler {
	key := fmt.Sprintf("%s_%d", CacheAi302Handler, tokenMappingID)
	if sdhandlers.Exists(key) {
		result, _ := sdhandlers.Get(key)
		return result.(*MessageHandler)
	}
	return nil
}

func GetLarkClient(tokenMappingID int) *lark.Client {
	key := fmt.Sprintf("%s_%d", CacheAi302LarkClient, tokenMappingID)
	if sdLarkClient.Exists(key) {
		result, _ := sdLarkClient.Get(key)
		return result.(*lark.Client)
	}
	return nil
}

func newLarkClient(config initialization.Config, feishuAppID string, feishuAppSecret string) *lark.Client {
	options := []lark.ClientOptionFunc{
		lark.WithLogLevel(larkcore.LogLevelDebug),
	}
	if config.FeishuBaseUrl != "" {
		options = append(options, lark.WithOpenBaseUrl(config.FeishuBaseUrl))
	}
	larkClient := lark.NewClient(feishuAppID, feishuAppSecret, options...)
	return larkClient
}

func UpdateHandler(id int) {
	key := fmt.Sprintf("%s_%d", CacheAi302Handler, id)
	if sdhandlers.Exists(key) {
		sdhandlers.Delete(key)
	}
	return
}

/*
func Handler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	return msgReceivedHandler(ctx, event)
}
*/

func ReadHandler(ctx context.Context, event *larkim.P2MessageReadV1) error {
	readerId := event.Event.Reader.ReaderId.OpenId
	//fmt.Printf("msg is read by : %v \n", *readerId)
	logger.Debugf("msg is read by : %v \n", *readerId)

	return nil
}

func CardHandler(tokenMappingID int) func(
	ctx context.Context,
	cardAction *larkcard.CardAction) (interface{}, error) {
	logger.Info("aaa%d", tokenMappingID)
	h := GetAi302Handler(tokenMappingID)
	logger.Info("bbb%+v", h)
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		//handlerType := judgeCardType(cardAction)
		logger.Info("11111")
		return h.cardHandler(ctx, cardAction)
	}
}

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
