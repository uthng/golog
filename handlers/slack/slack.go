package slack

import (
	"fmt"
	"sync"
	"time"

	"github.com/nlopes/slack"

	log "github.com/uthng/golog"
)

// Handler represents slack structure
type handler struct {
	verbose int

	token    string
	channel  string
	msgTitle string

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

// PrintMsg formats messages to post to slack channel
// according to logger informations
func (h *handler) PrintMsg(p int, l *log.Logger, level int, caller string, f string, v ...interface{}) error {

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.verbose >= level {
		//color := colors[level]

		switch p {
		case log.PRINT, log.PRINTF, log.PRINTLN:
			h.postWebhook(l, level, caller, f, v...)
		case log.PRINTW:
			//printw(l, level, prefix, ct, f, v...)
		default:
		}
	}

	return nil
}

func (h *handler) buildPrefixFields(l *log.Logger, level int, caller string) []interface{} {
	var kv []interface{}

	if l.flag&FTIMESTAMP != 0 {
		ts = time.Now().Format(l.timeFormat)
		kv = append(kv, "ts="+ts)
	}

	if l.flag&FCALLER != 0 {
		kv = append(kv, "caller="+caller)
	}

	return kv
}

//func (h *handler) buildLogFields(l *log.Logger, level int, v... interface{}) []interface{} {
//var pairs []interface{}

//kv := v

//if len(kv)%2 != 0 {
//kv = append(kv, "missing")
//}

//// if no key/value fields, return line after print message
//for i := 0; i < len(kv); i += 2 {
//// cast 1st elem = key to string
//k := cast.ToString(kv[i])
//if k == "" {
//k = "missing"
//}
//k = ct.SprintFunc()(k)

//// cast 2nd elem = value
//v := kv[i+1]
//pair := ""
//kind := reflect.ValueOf(v).Kind()

//if kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Struct || kind == reflect.Ptr {
//pair = fmt.Sprintf("%s=%+v", ct.SprintFunc()(k), v)
//} else {
//s := cast.ToString(v)
//pair = fmt.Sprintf("%s=%s", ct.SprintFunc()(k), quoteString(s))
//}
////pair = fmt.Sprintf("%s=%+v", k, v)
//pairs = append(pairs, pair)
//if i != len(kv)-2 {
//format += "%s "
//} else {
//format += "%s\n"
//}
//}
//}

func (h *handler) postWebhook(level int, kv ...interface{}) error {
	webhookMsg := &slack.WebhookMessage{}
	attachment := slack.Attachment{}

	if h.msgTitle != "" {
		attachment.Title = h.msgTitle
	}

	attachment.Color = colors[level]
	attachment.MarkdownIn = append(attachment.MarkdownIn, "text", "fields")

	return slack.PostWebhook(slackWebhookURL+h.token, webhookMsg)
}
