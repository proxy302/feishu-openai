package aihandlers

import (
	"context"
	"fmt"
	"start-feishubot/logger"
	"strings"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"

	"start-feishubot/initialization"
	"start-feishubot/services"
	"start-feishubot/services/openai"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 责任链
func chain(data *ActionInfo, actions ...Action) bool {
	for _, v := range actions {
		if !v.Execute(data) {
			return false
		}
	}
	return true
}

type MessageHandler struct {
	sessionCache   services.SessionServiceCacheInterface
	msgCache       services.MsgCacheInterface
	gpt            *openai.ChatGPT
	config         initialization.Config
	BotName        string
	TokenMappingID int
}

func (m MessageHandler) cardHandler(ctx context.Context,
	cardAction *larkcard.CardAction) (interface{}, error) {
	messageHandler := NewCardHandler(m)
	return messageHandler(ctx, cardAction, m.TokenMappingID)
}

func (m MessageHandler) Handler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	return m.msgReceivedHandler(ctx, event)
}

func judgeMsgType(event *larkim.P2MessageReceiveV1) (string, error) {
	msgType := event.Event.Message.MessageType

	switch *msgType {
	case "text", "image", "audio", "post":
		return *msgType, nil
	default:
		return "", fmt.Errorf("unknown message type: %v", *msgType)
	}
}

func (m MessageHandler) msgReceivedHandler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	handlerType := judgeChatType(event)
	logger.Debug("handlerType", handlerType)
	if handlerType == "otherChat" {
		fmt.Println("unknown chat type")
		return nil
	}
	logger.Debug("收到消息：", larkcore.Prettify(event.Event.Message))

	msgType, err := judgeMsgType(event)
	if err != nil {
		fmt.Printf("error getting message type: %v\n", err)
		return nil
	}

	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	chatId := event.Event.Message.ChatId
	mention := event.Event.Message.Mentions

	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}
	msgInfo := MsgInfo{
		handlerType: handlerType,
		msgType:     msgType,
		msgId:       msgId,
		chatId:      chatId,
		qParsed:     strings.Trim(parseContent(*content, msgType), " "),
		fileKey:     parseFileKey(*content),
		imageKey:    parseImageKey(*content),
		imageKeys:   parsePostImageKeys(*content),
		sessionId:   sessionId,
		mention:     mention,
	}
	data := &ActionInfo{
		ctx:            &ctx,
		handler:        &m,
		info:           &msgInfo,
		TokenMappingID: m.TokenMappingID,
	}
	actions := []Action{
		&ProcessedUniqueAction{TokenMappingID: m.TokenMappingID}, //避免重复处理
		&ProcessMentionAction{TokenMappingID: m.TokenMappingID},  //判断机器人是否应该被调用
		&AudioAction{TokenMappingID: m.TokenMappingID},           //语音处理
		&ClearAction{TokenMappingID: m.TokenMappingID},           //清除消息处理
		// &VisionAction{TokenMappingID: m.TokenMappingID},          //图片推理处理
		// &PicAction{TokenMappingID: m.TokenMappingID},             //图片处理
		&AIModeAction{TokenMappingID: m.TokenMappingID},   //模式切换处理
		&RoleListAction{TokenMappingID: m.TokenMappingID}, //角色列表处理
		&HelpAction{TokenMappingID: m.TokenMappingID},     //帮助处理
		// &BalanceAction{TokenMappingID: m.TokenMappingID},         //余额处理
		&RolePlayAction{TokenMappingID: m.TokenMappingID},      //角色扮演处理
		&MessageAction{TokenMappingID: m.TokenMappingID},       //消息处理
		&EmptyAction{TokenMappingID: m.TokenMappingID},         //空消息处理
		&StreamMessageAction{TokenMappingID: m.TokenMappingID}, //流式消息处理
	}
	chain(data, actions...)
	return nil
}

// var _ MessageHandlerInterface = (*MessageHandler)(nil)

func NewMessageHandler(gpt *openai.ChatGPT,
	config initialization.Config, botName string, tokenMappingID int) MessageHandlerInterface {
	return &MessageHandler{
		sessionCache:   services.GetSessionCache(),
		msgCache:       services.GetMsgCache(),
		gpt:            gpt,
		config:         config,
		BotName:        botName,
		TokenMappingID: tokenMappingID,
	}
}

func (m MessageHandler) judgeIfMentionMe(mention []*larkim.
	MentionEvent) bool {
	if len(mention) != 1 {
		return false
	}
	// return *mention[0].Name == m.config.FeishuBotName
	return *mention[0].Name == m.BotName
}

func AzureModeCheck(a *ActionInfo) bool {
	if a.handler.config.AzureOn {
		//sendMsg(*a.ctx, "Azure Openai 接口下，暂不支持此功能", a.info.chatId)
		return false
	}
	return true
}
