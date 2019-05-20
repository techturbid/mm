package main

import (
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// JSON enabled struct for the suggestion
type Message struct {
	ID         string  `json:"_id,omitempty" bson:"_id,omitempty"`
	Suggestion string  `json:"suggestion,omitempty" bson:"suggestion,omitempty"`
	Person     *Person `json:"Person,omitempty" bson:"Person,omitempty"`
}

// JSON Enabled Person struct
type Person struct {
	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var queueName string
var queueURL string

func main() {
	// Set the queue name to search
	flag.StringVar(&queueName, "n", "mm-sqs", "Queue name")
	flag.Parse()

	// Check if queue exists, otherwise asks to run terraform first
	if len(queueName) == 0 {
		flag.PrintDefaults()
		log.Fatal("Queue name required, run terraform first")
	}

	// Establish AWS Connection
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)

	// Terraform already created the queue with an state file, when attempting to create an existing Queue, AWS SDK for go will populate res variable with some cool stuff
	res, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})

	// Error handling for creating the queue
	if err != nil {
		log.Fatal("Could not create queue:", err)
		return
	}

	// Get the queueURL from the res variable as explained before
	queueURL = aws.StringValue(res.QueueUrl)

	// Print the Queue URL on the console
	log.Info("Queue url: ", queueURL)

	// Receive message
	receive_params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(3),  // 一次最多取幾個 message
		VisibilityTimeout:   aws.Int64(30), // 如果這個 message 沒刪除，下次再被取出來的時間
		WaitTimeSeconds:     aws.Int64(20), // long polling 方式取，會建立一條長連線並且等在那邊，直到 SQS 收到新 message 回傳給這條連線才中斷
	}
	receive_resp, err := svc.ReceiveMessage(receive_params)
	if err != nil {
		log.Println(err)
	}

	//var message Message
	//json.Unmarshal(receive_resp.Messages, message)
	// convert the suggestion to json
	//messageBytes, _ := json.Marshal(receive_resp.Messages.Body)
	//messageStr := string(messageBytes)

	fmt.Printf("[Receive message] \n%v \n\n", receive_resp.String())

	//fmt.Printf("[Receive message] \n%v \n\n", message)

	// Delete message
	for _, message := range receive_resp.Messages {
		delete_params := &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),  // Required
			ReceiptHandle: message.ReceiptHandle, // Required

		}
		_, err := svc.DeleteMessage(delete_params) // No response returned when successed.
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *message.MessageId)
	}
}
