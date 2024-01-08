package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func getQueryURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}
	return result, err
}
func sendMessage(sess *session.Session, queueURL *string) error {
	svc := sqs.New(sess)

	_, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("Harry Potter"),
			},
			"Author": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("J K Rowling"),
			},
			"WeeksOn": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String("Information about your favorite book"),
		QueueUrl:    queueURL,
	})

	if err != nil {
		return err
	}
	return nil

}
func main() {
	queue := flag.String("q", "golang-prac", "The name of the queue")
	flag.Parse()

	if *queue == "" {
		fmt.Println("The name of the queue is missing")
		return
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	result, err := getQueryURL(sess, queue)
	if err != nil {
		log.Fatal("Getting error while trying for URL")
		log.Fatal(err)
	}
	queueURL := result.QueueUrl

	err = sendMessage(sess, queueURL)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln("The message is sent to the queue")

}
