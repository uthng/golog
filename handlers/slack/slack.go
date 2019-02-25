package slack

import (
	"fmt"
	"strings"
	"sync"
	//"time"

	//"github.com/spf13/cast"
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
	title     string

	mutex sync.Mutex
}

var slackWebhookURL = "https://hooks.slack.com/services/"

// Color Map following levels
var colors = map[int]string{
	log.FATAL: "#cc0000",
	log.ERROR: "danger",
	log.WARN:  "warning",
	log.INFO:  "good",
	log.DEBUG: "#7e7e7c",
}

// New creates a new slack handler
func New(token, username, iconEmoji, iconURL, channel, title string, verbose int) log.Handler {
	s := &handler{
		verbose:   verbose,
		token:     token,
		username:  username,
		iconEmoji: iconEmoji,
		iconURL:   iconURL,
		channel:   channel,
		title:     title,
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
			flag := l.GetFlags()
			if flag&log.FFULLSTRUCTUREDLOG != 0 {
				return h.printwWebhook(level, fields, true)
			}

			return h.printwWebhook(level, fields, false)
		}
	}

	return nil
}

func (h *handler) printWebhook(level int, fields log.Fields) error {
	var webhookMsg *slack.WebhookMessage
	var attachment *slack.Attachment
	var prefix string

	webhookMsg = initWebhookMessage(h, level)
	attachment = &webhookMsg.Attachments[0]

	for _, p := range fields.Prefix {
		val := p.Value

		if p.Key == "level" {
			val = "*" + val + "*:"
		}
		prefix = prefix + val + " "
	}

	var message string
	message = fmt.Sprint(fields.Log[0].Value)

	attachment.Text = prefix + message

	return slack.PostWebhook(slackWebhookURL+h.token, webhookMsg)
}

func (h *handler) printwWebhook(level int, fields log.Fields, full bool) error {
	var webhookMsg *slack.WebhookMessage
	var attachment *slack.Attachment

	webhookMsg = initWebhookMessage(h, level)
	attachment = &webhookMsg.Attachments[0]

	for _, p := range fields.Log {
		if p.Key == "msg" && !full {
			attachment.Text = p.Value
		} else {
			field := slack.AttachmentField{
				Value: "*" + strings.Title(p.Key) + "*: " + p.Value,
			}

			attachment.Fields = append(attachment.Fields, field)
		}
	}

	for _, p := range fields.Prefix {
		field := slack.AttachmentField{
			Value: "*" + strings.Title(p.Key) + "*: " + p.Value,
		}

		attachment.Fields = append(attachment.Fields, field)
	}

	return slack.PostWebhook(slackWebhookURL+h.token, webhookMsg)
}

func initWebhookMessage(h *handler, level int) *slack.WebhookMessage {
	webhookMsg := &slack.WebhookMessage{
		Username:  h.username,
		IconEmoji: h.iconEmoji,
		IconURL:   h.iconURL,
		Channel:   h.channel,
	}

	attachment := slack.Attachment{}

	if h.title != "" {
		attachment.Title = h.title
	}

	attachment.Color = colors[level]
	attachment.MarkdownIn = append(attachment.MarkdownIn, "text", "fields")

	webhookMsg.Attachments = append(webhookMsg.Attachments, attachment)

	return webhookMsg
}
