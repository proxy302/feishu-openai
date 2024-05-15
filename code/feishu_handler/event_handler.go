package feishu_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"start-feishubot/aihandlers"
	"start-feishubot/initialization"
	"start-feishubot/logger"
	"start-feishubot/models"
	"start-feishubot/services/openai"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

	sdkginext "github.com/larksuite/oapi-sdk-gin"
)

func EventHandler(c *gin.Context) {
	reqBody := bytes.Buffer{}
	io.Copy(&reqBody, c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody.Bytes()))
	logger.Info("reqBody:", reqBody.String())

	data := make(map[string]interface{})
	json.Unmarshal(reqBody.Bytes(), &data)
	if val, ok := data["challenge"]; ok {
		c.JSON(200, gin.H{
			"challenge": val,
		})
		return
	}

	if val1, ok := data["header"]; ok {
		header := val1.(map[string]interface{})
		if val2, ok := header["token"]; ok {
			tok := val2.(string)
			robotMapping, _ := models.GetRobotMappingByFeishuToken(tok)
			if robotMapping != nil {
				tokenMapping, _ := models.GetTokenMappingByID(robotMapping.TokenID)
				token, _ := models.GetTokenByID(tokenMapping.ExternalTokenID)
				model, _ := models.GetModelByID(tokenMapping.ModelID)
				gpt := openai.Ai302NewChatGPT(*initialization.GetConfig(), model.Name, token.Value)
				handlers := aihandlers.InitHandlers(gpt, *initialization.GetConfig(), tokenMapping.ID, robotMapping.FeishuBotName, robotMapping.FeishuAppID, robotMapping.FeishuAppSecret)

				eventHandler := dispatcher.NewEventDispatcher(
					robotMapping.FeishuVerificationToken,
					robotMapping.FeishuEncryptKey).
					OnP2MessageReceiveV1(handlers.Handler).
					OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
						logger.Debugf("收到请求 %v", event.RequestURI)
						return aihandlers.ReadHandler(ctx, event)
					})
				fun := sdkginext.NewEventHandlerFunc(eventHandler)
				fun(c)
				return
			}

		}
	}
	c.JSON(200, gin.H{"ret": 201})
}
