package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

func main() {
	ctx := context.Background()

	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	projectID := os.Getenv("PROJECT_ID")
	locationID := os.Getenv("REGION")
	queueID := os.Getenv("QUEUE_ID")
	url := os.Getenv("ENDPOINT_URL")

	parent := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	requestBody := []byte(`{"message":"Hello, Cloud Tasks!"}`)
	serviceAccountEmail := os.Getenv("SA_EMAIL")

	req := &taskspb.CreateTaskRequest{
		Parent: parent,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
					Headers:    map[string]string{"Content-Type": "application/json"},
					Body:       requestBody,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: serviceAccountEmail,
							Audience:            url,
						},
					},
				},
			},
		},
	}

	createTask, err := client.CreateTask(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	fmt.Printf("Created task: %v\n", createTask.GetName())
}
