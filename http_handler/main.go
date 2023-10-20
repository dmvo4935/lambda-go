package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MyEvent struct {
	BucketName string `json:"bucket"`
}

type MyResponse struct {
	BucketName string   `json:"bucket"`
	Contents   []string `json:"contents"`
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

func (r *MyResponse) String() string {
	return fmt.Sprintf("{\"bucket\": \"%s\", \"contents\": %s}", r.BucketName, r.Contents)
}

func HandleRequest(ctx context.Context, event *MyEvent) (*MyResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}
	fmt.Printf("Event: %+v", event)
	contents, err := ListBucket(ctx, event.BucketName)

	return &MyResponse{BucketName: event.BucketName, Contents: contents}, err
}

func ListBucket(ctx context.Context, name string) ([]string, error) {
	fmt.Println("Trying to list")
	var results []string
	params := &s3.ListObjectsV2Input{Bucket: aws.String(name)}

	fmt.Printf("%+v\n", svc)
	fmt.Printf("%+v\n", *params.Bucket)

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
