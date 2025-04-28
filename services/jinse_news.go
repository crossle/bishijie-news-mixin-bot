package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	bot "github.com/MixinNetwork/bot-api-go-client/v3"
	"github.com/crossle/bishijie-news-mixin-bot/config"
	"github.com/crossle/bishijie-news-mixin-bot/models"
	"github.com/jasonlvhit/gocron"
)

type JinseNewsService struct{}

var jinseId int64

func sendJinseTopStoryToChannel(ctx context.Context, safeUser *bot.SafeUser) {
	stories, err := GetJinseStories()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := len(stories) - 1; i >= 0; i-- {
		story := stories[i]
		if story.CreatedAt > jinseId {
			log.Printf("Sending top story to channel...")
			jinseId = story.CreatedAt
			subscribers, err := models.FindSubscribers(ctx)
			if err != nil {
				log.Println("Error finding subscribers:", err)
				return
			}
			for _, subscriber := range subscribers {
				conversationId := bot.UniqueConversationId(config.MixinClientId, subscriber.UserId)
				data := base64.RawURLEncoding.EncodeToString([]byte(story.Content + " " + story.Link))
				mr := &bot.MessageRequest{
					ConversationId: conversationId,
					MessageId:      bot.UuidNewV4().String(),
					Category:       "PLAIN_TEXT",
					DataBase64:     data,
					RecipientId:    subscriber.UserId,
				}
				err = bot.PostMessages(ctx, []*bot.MessageRequest{mr}, safeUser)
				if err != nil {
					log.Println("bad send message", err)
				}
				log.Printf("Sent message to %s: %s", subscriber.UserId, story.Content)
			}
		} else {
			log.Printf("Same top jinse story ID: %d, no message sent.", jinseId)
		}
	}
}
func (service *JinseNewsService) Run(ctx context.Context) error {
	safeUser := bot.NewSafeUser(config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey)
	gocron.Every(5).Second().Do(sendJinseTopStoryToChannel, ctx, safeUser)
	<-gocron.Start()
	return nil
}
