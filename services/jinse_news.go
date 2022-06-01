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

type JinseStats struct {
	prevStoryId int
}

func (self JinseStats) getPrevTopStoryId() int {
	return self.prevStoryId
}

func (self *JinseStats) updatePrevTopStoryId(id int) {
	self.prevStoryId = id
}

func getTopJinseStory() LiveItem {
	stories, _ := GetJinseStories()
	return stories[len(stories)-1]
}

func sendJinseTopStoryToChannel(ctx context.Context, stats *JinseStats) {
	prevStoryId := stats.getPrevTopStoryId()
	stories, err := GetJinseStories()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := len(stories) - 1; i >= 0; i-- {
		story := stories[i]
		if story.ID > prevStoryId {
			log.Printf("Sending top story to channel...")
			stats.updatePrevTopStoryId(story.ID)
			subscribers, _ := models.FindSubscribers(ctx)
			for _, subscriber := range subscribers {
				conversationId := bot.UniqueConversationId(config.MixinClientId, subscriber.UserId)
				data := base64.StdEncoding.EncodeToString([]byte(story.Content + " " + story.Link))
				bot.PostMessage(ctx, conversationId, subscriber.UserId, bot.UuidNewV4().String(), "PLAIN_TEXT", data, config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey)
			}
		} else {
			log.Printf("Same top story ID: %d, no message sent.", prevStoryId)
		}
	}
}
func (service *JinseNewsService) Run(ctx context.Context) error {
	topStory := getTopJinseStory()
	stats := &JinseStats{topStory.ID}
	gocron.Every(5).Minutes().Do(sendJinseTopStoryToChannel, ctx, stats)
	<-gocron.Start()
	return nil
}
