package main
import (
	//"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
	"github.com/nlopes/slack"
	"log"
	"os"
	)
	
	
	type GrabBot struct {
		bot reddit.Bot
		slackClient *slack.Client
		slackChannel string
	}

	func NewGrabBot(configFileName string) *GrabBot {
		returnBot := new(GrabBot)
		newBot, err := reddit.NewBotFromAgentFile(configFileName, 0)
		if err != nil {
			log.Println("Failed to create bot handle: ", err)
			os.Exit(12)
		}
		returnBot.bot = newBot
		return returnBot
	}

	func (r *GrabBot) sendSlackMessage(slackChannel string ,message string) {

		channelID, timestamp, err := r.slackClient.PostMessage(slackChannel, slack.MsgOptionText(message, false))
		if err != nil {
			log.Printf("%s\n", err)
			return
		}
		log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	}

	func (r *GrabBot) sendSlackImage(slackChannel string ,message string,text string,pretext string, imageURL string) {

		attachment := slack.Attachment{
		Pretext: pretext,
		Text:    text,
		ImageURL : imageURL,
		}

		channelID, timestamp, err := r.slackClient.PostMessage(slackChannel, slack.MsgOptionText(message, false),slack.MsgOptionAttachments(attachment))
		if err != nil {
			log.Printf("%s\n", err)
			return
		}
		log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	}


	func (r *GrabBot) harvestSubreddit(subRedditName string ) reddit.Harvest {
	harvest, err := r.bot.Listing(subRedditName, "")
 	if err != nil {
       log.Println("Failed to fetch " + subRedditName, err)
       return harvest
	}
	return harvest
	}

	func (r *GrabBot) harvestNewestInSubreddit(subRedditName string ,amount int) {
		harvest := r.harvestSubreddit(subRedditName)
		log.Println(harvest.Posts)
		log.Println(len(harvest.Posts))
		for _, post := range harvest.Posts[:amount] {
			//GrabBot.sendSlackImage(r.slackChannel,post.Subreddit,post.Title,string(post.Ups),post.URL)
			log.Printf(`[%d] %s posted "%s"\n`, post.CreatedUTC, post.Author, post.Title)
		}
	}

//-----------
//EVENTS
//-----------


	func (r *GrabBot) Post(post *reddit.Post) error {
		//r.sendSlackImage(r.slackChannel,post.Subreddit,post.Title,string(post.Ups),post.URL)
		log.Printf(`%s posted "%s"\n`, post.Author, post.Title)
		return nil
	}

func main() {

	GrabBot := NewGrabBot("agent.config")


	GrabBot.harvestNewestInSubreddit("/r/pics/",12)

	//cfg := graw.Config{Subreddits: []string{"pics"}}
	//handler := GrabBot
	//handler := &GrabBot{bot: bot}
	/*
    if _, wait, err := graw.Run(GrabBot, GrabBot.bot, cfg); err != nil {
            log.Println("Failed to start graw run: ", err)
    } else {
            log.Println("graw run failed: ", wait())
	}
	*/


}