Slack handler
======
This handler uses Slack Incomming Webhook to post a log message to a Slack channel.

### Usage
To use this handler in Golog, a new Slack handler must be created and added to Golog handler list as below:
```golang
import (
  "github.com/uthng/golog/handlers/slack"
)

logger := golog.NewLogger()
logger.AddHandler(slack.New(
  "https://hooks.slack.com/services/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",
  "Golog",
  "",
  "",
  "C4JLPQB7X",
  "Structured Log",
  golog.INFO))
```
`New` function takes 7 arguments:
- **WebhookURL:** URL of Incomming Webhook
- **Username:** Custom username
- **IconEmoji:** Custom emoji
- **IconURL:** Custom emoji URL
- **Channel:** Channel ID
- **Title:** Attachment message title
- **Verbose:** Level from which log will be sent to Slack

### Message formatting
Slack messages are formatted in 3 ways according to how log message is printed in Golog:

#### Simple log
When log messages are printed with functions such as: Debug, Info, Warn, Error, Debugf, Infof, Warnf, Errorf, Debugln, Infoln, Warnln and Errorln, they are formatted as a simple Slack message as the following example:

![Slack simple message without prefix](./images/slack_msg_simple_without_prefix.png)

![Slack simple message with prefix](./images/slack_msg_simple_with_prefix.png)

#### Semi structured log
When using Debugw, Infow, Warnw and Errorw, log messages are formatted in Slack as the following:
- Message: is used as attachment text
- Key/value fields: will be attachment fields

![Slack semi structured log without prefix](./images/slack_msg_semi_structured_without_prefix.png)

![Slack semi structured with prefix](./images/slack_msg_semi_structured_with_prefix.png)

#### Full structured log
When using Debugw, Infow, Warnw and Errorw with the flag FFULLSTRUCTUREDLOG, each key/value field of structured log is formatted as an attachment field

![Slack full structured log without prefix](./images/slack_msg_full_structured_without_prefix.png)

![Slack full structured with prefix](./images/slack_msg_full_structured_with_prefix.png)

