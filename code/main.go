package main

import (
	"start-feishubot/feishu_handler"
	"start-feishubot/gredis"
	"start-feishubot/initialization"
	"start-feishubot/logger"
	"start-feishubot/models"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

func main() {
	initialization.InitRoleList()
	pflag.Parse()
	config := initialization.GetConfig()
	// initialization.LoadLarkClient(*config)
	// gpt := openai.NewChatGPT(*config)
	// handlers.InitHandlers(gpt, *config)

	models.Setup()
	gredis.Setup()

	//eventHandler := dispatcher.NewEventDispatcher(
	//	config.FeishuAppVerificationToken, config.FeishuAppEncryptKey).
	//	OnP2MessageReceiveV1(handlers.Handler).
	//	OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
	//		logger.Debugf("收到请求 %v", event.RequestURI)
	//		return handlers.ReadHandler(ctx, event)
	//	})

	// cardHandler := larkcard.NewCardActionHandler(
	// 	config.FeishuAppVerificationToken, config.FeishuAppEncryptKey,
	// 	handlers.CardHandler())

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/webhook/event", feishu_handler.EventHandler)
	// r.POST("/webhook/card", feishu_handler.CardHandler)
	r.POST("/webhook/update", feishu_handler.UpdateHandler)
	// 	r.POST("/webhook/event",
	// 		sdkginext.NewEventHandlerFunc(eventHandler))
	// 	r.POST("/webhook/card",
	// 		sdkginext.NewCardActionHandlerFunc(
	// 			cardHandler))

	if err := initialization.StartServer(*config, r); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}
