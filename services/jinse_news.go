package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/crossle/bishijie-news-mixin-bot/config"
	"github.com/crossle/bishijie-news-mixin-bot/models"
	"github.com/jasonlvhit/gocron"
)

type JinseNewsService struct{}

var jinseId int64

func sendJinseTopStoryToChannel(ctx context.Context) {
	stories, err := GetJinseStories()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := len(stories) - 1; i >= 0; i-- {
		story := stories[i]
		if story.ID > jinseId {
			log.Printf("Sending top story to channel...")
			jinseId = story.ID
			subscribers, _ := models.FindSubscribers(ctx)
			for _, subscriber := range subscribers {
				conversationId := bot.UniqueConversationId(config.MixinClientId, subscriber.UserId)
				data := base64.StdEncoding.EncodeToString([]byte(story.Content + " " + story.Link))
				bot.PostMessage(ctx, conversationId, subscriber.UserId, bot.UuidNewV4().String(), "PLAIN_TEXT", data, config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey)
			}
		} else {
			log.Printf("Same top story ID: %d, no message sent.", jinseId)
		}
	}
}
func (service *JinseNewsService) Run(ctx context.Context) error {
	gocron.Every(5).Minutes().Do(sendJinseTopStoryToChannel, ctx)
	<-gocron.Start()
	return nil
}
