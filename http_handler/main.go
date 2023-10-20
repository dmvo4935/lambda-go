package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MyEvent struct {
	BucketName string `json:"bucket"`
}

type MyResponse struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      string            `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
}

type S3Response struct {
	Bucket  string   `json:"bucket"`
	Objects []string `json:"objects"`
}

var svc *s3.Client

func init() {
	fmt.Println("Initializing")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}
	svc = s3.NewFromConfig(cfg)
}

// func (r *MyResponse) String() string {
// 	return fmt.Sprintf("{\"bucket\": \"%s\", \"contents\": %s}", r.BucketName, r.Contents)
// }

func HandleRequest(ctx context.Context, event *events.APIGatewayProxyRequest) (*MyResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}
	event_json, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Event: %s\n", string(event_json))

	var request *MyEvent
	err = json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return nil, err
	}

	contents, err := ListBucket(ctx, request.BucketName)
	if err != nil {
		return nil, err
	}

	response_data, err := json.Marshal(S3Response{Bucket: request.BucketName, Objects: contents})

	fmt.Printf(string(response_data))

	return &MyResponse{
		IsBase64Encoded: false,
		Headers:         make(map[string]string),
		StatusCode:      "200",
		Body:            string(response_data)}, err
}

func ListBucket(ctx context.Context, name string) ([]string, error) {
	fmt.Println("Trying to list")
	var results []string
	params := &s3.ListObjectsV2Input{Bucket: aws.String(name)}

	fmt.Printf("Reading from bucket: %+v\n", *params.Bucket)

	objects, err := svc.ListObjectsV2(ctx, params)
	fmt.Printf("%+v %v\n\n", &objects, err)

	if err != nil {
		return nil, err
	}

	for i := range objects.Contents {
		results = append(results, *objects.Contents[i].Key)
	}

	return results, err
}

func main() {
	lambda.Start(HandleRequest)

	// event := &MyEvent{BucketName: "sam-hello-world-4952"}

	// if contents, err := HandleRequest(context.TODO(), event); err != nil {
	// 	fmt.Printf("%v", err)
	// } else {
	// 	fmt.Println(*contents)
	// }
}
