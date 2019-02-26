package slack

import (
	"encoding/json"
	//"fmt"
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	//"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/uthng/golog"
	anotherSlack "github.com/uthng/slack"
	//utils "github.com/uthng/goutils"
)

func TestHandlerSlackSimpleLog(t *testing.T) {
	output := anotherSlack.WebhookMessage{}
	outputDebug := anotherSlack.WebhookMessage{}
	outputInfo := anotherSlack.WebhookMessage{
		Channel:  "C4JLPQB7X",
		Username: "Golog",
		Attachments: []anotherSlack.Attachment{
			{
				MarkdownIn: []string{"text", "fields"},
				Color:      "good",
				Text:       "*INFO*: This is info log level info\n",
			},
		},
	}
	outputWarn := anotherSlack.WebhookMessage{
		Channel:  "C4JLPQB7X",
		Username: "Golog",
		Attachments: []anotherSlack.Attachment{
			{
				MarkdownIn: []string{"text", "fields"},
				Color:      "warning",
				Text:       "*WARN*: This is warn log level warning\n",
			},
		},
	}
	outputErrorSimple := anotherSlack.WebhookMessage{
		Channel:  "C4JLPQB7X",
		Username: "Golog",
		Attachments: []anotherSlack.Attachment{
			{
				MarkdownIn: []string{"text", "fields"},
				Color:      "danger",
				Text:       "*ERROR*: This is error log\n",
			},
		},
	}
	outputErrorSemiStructured := anotherSlack.WebhookMessage{
		Channel:  "C4JLPQB7X",
		Username: "Golog",
		Attachments: []anotherSlack.Attachment{
			{
				MarkdownIn: []string{"text", "fields"},
				Color:      "danger",
				Text:       "This is error log",
				Fields: []anotherSlack.AttachmentField{
					{
						Value: "*Field1*: value1",
					},
					{
						Value: "*Field2*: value2",
					},
					{
						Value: "*Level*: ERROR",
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&output)
	}))
	defer server.Close()

	serverAddr := server.Listener.Addr().String()

	logger := golog.NewLogger()

	logger.AddHandler(New(
		"http://"+serverAddr,
		"Golog",
		"",
		"",
		"C4JLPQB7X",
		"",
		golog.INFO))

	logger.Debug("This is debug log")
	assert.Equal(t, outputDebug, output)

	logger.Info("This is info log ", "level info\n")
	assert.Equal(t, outputInfo, output)

	logger.Warnln("This is warn log", "level warning")
	assert.Equal(t, outputWarn, output)

	logger.Errorf("This is %s log\n", "error")
	assert.Equal(t, outputErrorSimple, output)

	logger.Errorw("This is error log", "field1", "value1", "field2", "value2")
	assert.Equal(t, outputErrorSemiStructured, output)
}

func TestHandlerSlackStructuredLog(t *testing.T) {
	output := anotherSlack.WebhookMessage{}
	outputDebug := anotherSlack.WebhookMessage{}
	outputInfo := anotherSlack.WebhookMessage{
		Channel:  "C4JLPQB7X",
		Username: "Golog",
		Attachments: []anotherSlack.Attachment{
			{
				MarkdownIn: []string{"text", "fields"},
				Color:      "good",
				Title:      "Structured Log",
				Fields: []anotherSlack.AttachmentField{
					{
						Value: "*Msg*: This is info log",
					},
					{
						Value: "*Field1*: value1",
					},
					{
						Value: "*Field2*: value2",
					},
					{
						Value: "*Caller*: slack_test.go:225:TestHandlerSlackStructuredLog",
					},
					{
						Value: "*Level*: INFO",
					},
				},
			},
		},
	}
	outputWarn := anotherSlack.WebhookMessage{
		Channel:  "C4JLPQB7X",
		Username: "Golog",
		Attachments: []anotherSlack.Attachment{
			{
				MarkdownIn: []string{"text", "fields"},
				Color:      "warning",
				Title:      "Structured Log",
				Fields: []anotherSlack.AttachmentField{
					{
						Value: "*Msg*: This is warn log",
					},
					{
						Value: "*Field1*: value1",
					},
					{
						Value: "*Field2*: value2",
					},
					{
						Value: "*Caller*: slack_test.go:231:TestHandlerSlackStructuredLog",
					},
					{
						Value: "*Level*: WARN",
					},
				},
			},
		},
	}
	outputError := anotherSlack.WebhookMessage{
		Channel:  "C4JLPQB7X",
		Username: "Golog",
		Attachments: []anotherSlack.Attachment{
			{
				MarkdownIn: []string{"text", "fields"},
				Color:      "danger",
				Title:      "Structured Log",
				Fields: []anotherSlack.AttachmentField{
					{
						Value: "*Msg*: This is error log",
					},
					{
						Value: "*Field1*: value1",
					},
					{
						Value: "*Field2*: value2",
					},
					{
						Value: "*Caller*: slack_test.go:237:TestHandlerSlackStructuredLog",
					},
					{
						Value: "*Level*: ERROR",
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&output)
	}))
	defer server.Close()

	serverAddr := server.Listener.Addr().String()

	logger := golog.NewLogger()
	logger.SetTimeFormat("2006-01-02 15:04:05")

	logger.SetFlags(golog.FTIMESTAMP | golog.FCALLER | golog.FFULLSTRUCTUREDLOG)
	logger.AddHandler(New(
		"http://"+serverAddr,
		"Golog",
		"",
		"",
		"C4JLPQB7X",
		"Structured Log",
		golog.INFO))

	logger.Debugw("This is debug log", "field1", "value1", "field2", "value2")
	assert.Equal(t, outputDebug, output)

	ts := time.Now().Format("2006-01-02 15:04:05")
	logger.Infow("This is info log", "field1", "value1", "field2", "value2")
	// Insert ts in AttachmentField
	outputInfo = insertTsInOutput(outputInfo, ts)
	assert.Equal(t, outputInfo, output)

	ts = time.Now().Format("2006-01-02 15:04:05")
	logger.Warnw("This is warn log", "field1", "value1", "field2", "value2")
	// Insert ts in AttachmentField
	outputWarn = insertTsInOutput(outputWarn, ts)
	assert.Equal(t, outputWarn, output)

	ts = time.Now().Format("2006-01-02 15:04:05")
	logger.Errorw("This is error log", "field1", "value1", "field2", "value2")
	// Insert ts in AttachmentField
	outputError = insertTsInOutput(outputError, ts)
	assert.Equal(t, outputError, output)
}

func insertTsInOutput(output anotherSlack.WebhookMessage, ts string) anotherSlack.WebhookMessage {

	fields := output.Attachments[0].Fields

	field := anotherSlack.AttachmentField{
		Value: "*Ts*: " + ts,
	}

	fields = append(fields[:3], append([]anotherSlack.AttachmentField{field}, fields[3:]...)...)

	output.Attachments[0].Fields = fields

	return output
}
