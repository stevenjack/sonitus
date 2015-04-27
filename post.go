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

type payloadMessage struct {
	Short bool   `json:"short"`
	Title string `json:"title"`
	Value string `json:"value"`
}

func main() {

	sns := os.Args[2]

	awsMain := &Message{}
	err := json.Unmarshal([]byte(sns), &awsMain)
	if err != nil {
		panic(err)
	}

	messageBody := awsMain.Sns.Message

	awsBody := &MessageBody{}
	err2 := json.Unmarshal([]byte(messageBody), &awsBody)
	if err2 != nil {
		panic(err2)
	}

	alarm(awsBody)

}

func alarm(awsBody *MessageBody) {

	var checkStatus []byte

	alarmState := awsBody.NewStateValue
	if alarmState == "ALARM" {
		checkStatus = slack(awsBody, alarmState, "#F35A00")
		send(checkStatus)
	} else if alarmState == "OK" {
		checkStatus = slack(awsBody, alarmState, "#2ecc71")
		send(checkStatus)
	} else {
		return
	}
}

func slack(awsBody *MessageBody, alarmState string, colour string) []byte {

	alarmMessageName := &payloadMessage{
		Short: true,
		Title: alarmState,
		Value: awsBody.AlarmName,
	}

	messageName, err := json.Marshal(alarmMessageName)

	if err != nil {
		panic(err)
	}

	alarmMessageLink := &payloadMessage{
		Short: true,
		Title: "Client",
		Value: fmt.Sprintf("<https://eu-west-1.console.aws.amazon.com/cloudwatch/home#alarm:alarmFilter=ANY;name= %s|AWS Console>", awsBody.AlarmName),
	}

	messageLink, err := json.Marshal(alarmMessageLink)

	if err != nil {
		panic(err)
	}

	alarmMessageDesc := &payloadMessage{
		Short: true,
		Title: "Description",
		Value: awsBody.NewStateReason,
	}

	messageDesc, err := json.Marshal(alarmMessageDesc)

	if err != nil {
		panic(err)
	}

	alarmMessageTime := &payloadMessage{
		Short: true,
		Title: "Alert Time",
		Value: awsBody.StateChangeTime,
	}

	messageTime, err := json.Marshal(alarmMessageTime)

	if err != nil {
		panic(err)
	}

	var jsonStr = []byte(fmt.Sprintf("{\"username\": \"Cloudwatch\",\"attachments\":[{\"fallback\":\"AWS Alert\"},{\"fields\":[ %s, %s, %s, %s],\"color\": \"%s\"}]}", messageName, messageLink, messageDesc, messageTime, colour))

	return jsonStr

}

func send(jsonStr []byte) {
	slackURL := os.Args[1]

	req, err := http.NewRequest("POST", slackURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
