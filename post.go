package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Message struct {
	EventSource          string
	EventVersion         string
	EventSubscriptionArn string

	Sns struct {
		Type     string
		TopicArn string
		Subject  string
		Message  string
	}
	Timestamp        string
	SignatureVersion int
	Signature        string
	SigningCertUrl   string
	UnsubscribeUrl   string
}

type MessageBody struct {
	AWSAccountId     string `json:"AWSAccountId"`
	AlarmDescription string `json:"AlarmDescription"`
	AlarmName        string `json:"AlarmName"`
	NewStateReason   string `json:"NewStateReason"`
	NewStateValue    string `json:"NewStateValue"`
	OldStateValue    string `json:"OldStateValue"`
	Region           string `json:"Region"`
	StateChangeTime  string `json:"StateChangeTime"`
	Trigger          struct {
		ComparisonOperator string `json:"ComparisonOperator"`
		Dimensions         []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"Dimensions"`
		EvaluationPeriods float64     `json:"EvaluationPeriods"`
		MetricName        string      `json:"MetricName"`
		Namespace         string      `json:"Namespace"`
		Period            float64     `json:"Period"`
		Statistic         string      `json:"Statistic"`
		Threshold         float64     `json:"Threshold"`
		Unit              interface{} `json:"Unit"`
	} `json:"Trigger"`
}

func main() {

	sns := os.Args[2]

	awsMain := &Message{}
	err := json.Unmarshal([]byte(sns), &awsMain)
	check(err)

	messageBody := awsMain.Sns.Message

	awsBody := &MessageBody{}
	err2 := json.Unmarshal([]byte(messageBody), &awsBody)
	check(err2)

	slackMessage := slack(awsBody)
	send(slackMessage)

}

func slack(awsBody *MessageBody) []byte {
	messageName := fmt.Sprintf("{\"title\":\"Alarm\",\"value\":\" %s\",\"short\":true}", awsBody.AlarmName)
	messageLink := fmt.Sprintf("{\"title\":\"Client\",\"value\":\" <https://eu-west-1.console.aws.amazon.com/cloudwatch/home#alarm:alarmFilter=ANY;name= %s|AWS Console>\",\"short\":true}", awsBody.AlarmName)
	messageDesc := fmt.Sprintf("{\"title\":\"Description\",\"value\":\" %s\",\"short\":true}", awsBody.NewStateReason)
	messageTime := fmt.Sprintf("{\"title\":\"Alert Time\",\"value\":\" %s\",\"short\":true}", awsBody.StateChangeTime)

	var jsonStr = []byte(fmt.Sprintf("{\"attachments\":[{\"fallback\":\"AWS Alert\"},{\"fields\":[ %s, %s, %s, %s],\"color\":\"#F35A00\"}]}", messageName, messageLink, messageDesc, messageTime))

	return jsonStr

}

func send(jsonStr []byte) {
	slackUrl := os.Args[1]

	req, err := http.NewRequest("POST", slackUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
