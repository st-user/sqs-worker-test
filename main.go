package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx := context.TODO()
	// Load environment variables named SQS_QUEUE_URL
	queueURL := os.Getenv("SQS_QUEUE_URL")
	pollingIntervalStr := os.Getenv("POLLING_INTERVAL")

	pollingInterval, err := strconv.Atoi(pollingIntervalStr)
	if err != nil {
		fmt.Println("Failed to convert POLLING_INTERVAL to int so use the default interval:", err)
		pollingInterval = 1
	}
	fmt.Printf("SQS_QUEUE_URL: %s\n", queueURL)
	fmt.Printf("POLLING_INTERVAL: %d\n", pollingInterval)

	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Failed to load AWS configuration:", err)
		os.Exit(1)
	}

	// Create an SQS client using the AWS SDK configuration
	svc := sqs.NewFromConfig(cfg)

	// Define the maximum number of tasks to pull from the queue
	var maxTasks int32 = 10

	// Define the duration for which the receive message request waits for tasks
	var waitTimeSeconds int32 = 20

	for {

		fmt.Println("Waiting for tasks...")

		resp, err := svc.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: maxTasks,
			WaitTimeSeconds:     waitTimeSeconds,
		})

		// Pull tasks from the SQS queue
		if err != nil {
			fmt.Println("Failed to pull tasks from SQS queue:", err)
			// os.Exit(1)
			continue
		}

		// Process the tasks pulled from the SQS queue
		for _, message := range resp.Messages {
			fmt.Println("Task:", *message.Body)

			// Delete the task from the SQS queue
			_, err = svc.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				fmt.Println("Failed to delete task from SQS queue:", err)
				os.Exit(1)
			}
		}

		// sleep for pollingInterval seconds
		time.Sleep(time.Duration(pollingInterval) * time.Second)
	}

}
