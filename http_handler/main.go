package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
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
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}
	svc = s3.NewFromConfig(cfg)
}

func (r *MyResponse) String() string {
	return fmt.Sprintf("{\"bucket\": \"%s\", \"contents\": %s}", r.BucketName, r.Contents)
}

func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}
	contents, err := ListBucket(ctx, event.BucketName)
	response := fmt.Sprintf("%s", MyResponse{BucketName: event.BucketName, Contents: contents})
	return &response, err
}

func ListBucket(ctx context.Context, name string) ([]string, error) {
	var results []string
	params := s3.ListObjectsV2Input{Bucket: &name}
	objects, err := svc.ListObjectsV2(ctx, &params)

	for i := range objects.Contents {
		results = append(results, *objects.Contents[i].Key)
	}
	return results, err
}

func main() {
	lambda.Start(HandleRequest)
}
