Sonitus
==========

Sontius is a AWS Lambda project to send alert messages to various chat services such as Slack.  This is achieved by using Go to process the SNS message and send a formatted message on another service.  The main application is written in Go as I wanted to learn some more Go and not JavaScript.

## Usage
Edit the index.js file to include your Slack URL, then Zip up the binary and index.js and upload to Lambda.  From there, subscribe your Alarms to a topic and add a Lambda subscription to that topic.

## Building
As this is designed for AWS, we will be targeting the Linux platform.  There is a pre build binary available, or you can build yourself, I recommend [gox](https://github.com/mitchellh/gox).

## Development
To replicate Lambda, there is a debug folder with a example Alarm message from SNS.  Build a local binary and put the path to that binary into index.js instead of `/var/task/lambda`.  Simply run `node start.js` and it will send a message.
