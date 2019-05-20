package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gorilla/mux"
)

// JSON enabled struct for the suggestion
type Message struct {
	ID         string  `json:"_id,omitempty"`
	Suggestion string  `json:"suggestion,omitempty"`
	Person     *Person `json:"Person,omitempty"`
}

// JSON Enabled Person struct
type Person struct {
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
}

var queueName string
var queueURL string

func SendSQSMessageEndpoint(w http.ResponseWriter, r *http.Request) {
	var message Message
	json.NewDecoder(r.Body).Decode(&message)

  // Print a few messages detailing the POST information
	fmt.Println("---------------")
	fmt.Println("--New Message--")
	fmt.Println("---------------")
	fmt.Println("Message ID: ", message.ID)
	fmt.Println("Author: ", message.Person.FirstName, message.Person.LastName)
	fmt.Println("Suggestion: ", message.Suggestion)
	fmt.Println("---------------")

	// call the the function to send the message to SQS
	SendSQSMessage(message)

  // Return a valid JSON as response
	json.NewEncoder(w).Encode(message)
	return
}

func SendSQSMessage(message Message) error {
	log.Info("Sending suggestion to SQS")
  //establish AWS session
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)

  // convert the suggestion to json
	messageBytes, _ := json.Marshal(message)
  messageStr := string(messageBytes)

  // set the parameters for .SendMessage
	params := &sqs.SendMessageInput{
		MessageBody: aws.String(messageStr),
		QueueUrl:    aws.String(queueURL),
	}

  // Send the message while handle errors
	resp, err := svc.SendMessage(params)

	if err != nil {
		return err
	}

  //Print the SendMessage result
	fmt.Println("")
	log.Info("Message Sent!")
	fmt.Println(resp)
	fmt.Println("")
	return nil
}

func main() {
	fmt.Println("Starting the application...")

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

	router := mux.NewRouter()

  // Create the handler for /api POST requests
	router.HandleFunc("/api", SendSQSMessageEndpoint).Methods("POST")
	log.Info("Initializing API on port 8123")
	log.Fatal(http.ListenAndServe(":8123", router))

}
