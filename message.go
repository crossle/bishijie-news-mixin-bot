package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client"

	"github.com/crossle/bishijie-news-mixin-bot/config"
	"github.com/crossle/bishijie-news-mixin-bot/durable"
	"github.com/crossle/bishijie-news-mixin-bot/models"
	"github.com/crossle/bishijie-news-mixin-bot/session"
)

func StartBlaze(db *sql.DB) error {
	logger, err := durable.NewLoggerClient("", true)
	if err != nil {
		return err
	}
	defer logger.Close()
	ctx, cancel := newBlazeContext(db, logger)
	defer cancel()

	for {
		blazeClient := bot.NewBlazeClient(config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey)
		channelMessage := ResponseMessage{blaze: blazeClient}
		if err := blazeClient.Loop(ctx, channelMessage); err != nil {
			session.Logger(ctx).Error(err)
		}
		session.Logger(ctx).Info("connection loop end")
		time.Sleep(time.Second)
	}
	return nil
}

func newBlazeContext(db *sql.DB, client *durable.LoggerClient) (context.Context, context.CancelFunc) {
	ctx := session.WithLogger(context.Background(), durable.BuildLogger(client, "blaze", nil))
	ctx = session.WithDatabase(ctx, db)
	return context.WithCancel(ctx)
}

type ResponseMessage struct {
	blaze *bot.BlazeClient
}

func (r ResponseMessage) OnMessage(ctx context.Context, msg bot.MessageView, uid string) error {
	if msg.Category != bot.MessageCategorySystemAccountSnapshot && msg.Category != bot.MessageCategorySystemConversation && msg.ConversationId == bot.UniqueConversationId(config.MixinClientId, msg.UserId) {
		if msg.Category == "PLAIN_TEXT" {
			data, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				return bot.BlazeServerError(ctx, err)
			}
			if "/start" == string(data) {
				_, err = models.CreateSubscriber(ctx, msg.UserId)
				if err == nil {
					if err := r.blaze.SendPlainText(ctx, msg, "订阅成功"); err != nil {
						return bot.BlazeServerError(ctx, err)
					}
				}
			} else if "/stop" == string(data) {
				err = models.RemoveSubscriber(ctx, msg.UserId)
				if err == nil {
					if err := r.blaze.SendPlainText(ctx, msg, "已取消订阅"); err != nil {
						return bot.BlazeServerError(ctx, err)
					}
				}
			} else {
				content := `请更新 Mixin Messenger 到最新版本 0.14.1+`
				if err := r.blaze.SendAppButton(ctx, msg.ConversationId, msg.UserId, "点我订阅", "input:/start", "#AA4848"); err != nil {
					return bot.BlazeServerError(ctx, err)
				}
				if err := r.blaze.SendAppButton(ctx, msg.ConversationId, msg.UserId, "点我取消订阅", "input:/stop", "#AA4848"); err != nil {
					return bot.BlazeServerError(ctx, err)
				}
				if err := r.blaze.SendPlainText(ctx, msg, content); err != nil {
					return bot.BlazeServerError(ctx, err)
				}
			}
		} else {
			content := `请更新 Mixin Messenger 到最新版本 0.14.1+`
			if err := r.blaze.SendAppButton(ctx, msg.ConversationId, msg.UserId, "点我订阅", "input:/start", "#AA4848"); err != nil {
				return bot.BlazeServerError(ctx, err)
			}
			if err := r.blaze.SendAppButton(ctx, msg.ConversationId, msg.UserId, "点我取消订阅", "input:/stop", "#AA4848"); err != nil {
				return bot.BlazeServerError(ctx, err)
			}
			if err := r.blaze.SendPlainText(ctx, msg, content); err != nil {
				return bot.BlazeServerError(ctx, err)
			}

		}
	}
	return nil
}
