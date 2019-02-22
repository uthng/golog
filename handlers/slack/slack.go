package slack

import (
	"fmt"
	"sync"
	//"time"

	"github.com/spf13/cast"
	"github.com/uthng/slack"

	log "github.com/uthng/golog"
)

// Handler represents slack structure
type handler struct {
	verbose int

	token     string
	username  string
	iconEmoji string
	iconURL   string
	channel   string
	msgTitle  string

	mutex sync.Mutex
}

var slackWebhookURL = "https://hooks.slack.com/services/"

// Color Map following levels
var colors = map[int]string{
	log.FATAL: "#fffff",
	log.ERROR: "#fffff",
	log.WARN:  "#fffff",
	log.INFO:  "#fffff",
	log.DEBUG: "#fffff",
}

// New creates a new slack handler
func New(token, username, iconEmoji, iconURL, channel, msgTitle string, verbose int) log.Handler {
	s := &handler{
		verbose:   verbose,
		token:     token,
		username:  username,
		iconEmoji: iconEmoji,
		iconURL:   iconURL,
		channel:   channel,
		msgTitle:  msgTitle,
	}

	return s
}

// PrintMsg formats messages to post to slack channel
// according to logger informations
func (h *handler) PrintMsg(p int, l *log.Logger, level int, fields log.Fields) error {

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.verbose >= level {
		//color := colors[level]
		switch p {
		case log.PRINT, log.PRINTF, log.PRINTLN:
			return h.printWebhook(level, fields)
		case log.PRINTW:
			break
		}
	}

	return nil
}

func (h *handler) printWebhook(level int, fields log.Fields) error {
	webhookMsg := &slack.WebhookMessage{
		Username:  h.username,
		IconEmoji: h.iconEmoji,
		IconURL:   h.iconURL,
		Channel:   h.channel,
	}

	attachment := slack.Attachment{}

	if h.msgTitle != "" {
		attachment.Title = h.msgTitle
	}

	attachment.Color = colors[level]
	attachment.MarkdownIn = append(attachment.MarkdownIn, "text", "fields")

	var prefix string
	for _, p := range fields.Prefix {
		prefix = prefix + cast.ToString(p.Value) + " "
	}

	var message string
	message = fmt.Sprint(fields.Log[0].Value)

	attachment.Text = prefix + message

	webhookMsg.Attachments = append(webhookMsg.Attachments, attachment)

	return slack.PostWebhook(slackWebhookURL+h.token, webhookMsg)
}
