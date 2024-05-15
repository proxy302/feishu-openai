package feishu_handler

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"start-feishubot/aihandlers"
	"start-feishubot/initialization"
	"start-feishubot/logger"
	"start-feishubot/models"
	"start-feishubot/services/openai"

	"github.com/gin-gonic/gin"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

	sdkginext "github.com/larksuite/oapi-sdk-gin"
)

func CardHandler(c *gin.Context) {
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
				aihandlers.InitHandlers(gpt, *initialization.GetConfig(), tokenMapping.ID, robotMapping.FeishuBotName, robotMapping.FeishuAppID, robotMapping.FeishuAppSecret)

				cardHandler := larkcard.NewCardActionHandler(
					robotMapping.FeishuVerificationToken, robotMapping.FeishuEncryptKey,
					aihandlers.CardHandler(tokenMapping.ID))

				fun := sdkginext.NewCardActionHandlerFunc(cardHandler)
				fun(c)
				return
			}

		}
	}
	c.JSON(200, gin.H{"ret": 201})
}
