package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/davecgh/go-spew/spew"
)

type S3 struct {
	client *s3.Client
}

func NewS3() *S3 {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	//spew.Dump(cfg)
	cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "http://localhost:9000",
			SigningRegion: "us-east-1",
		}, nil
	})
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	cl := s3.NewFromConfig(cfg)
	return &S3{
		client: cl,
	}
}

func (*S3) CreateBucket(name string) {

}

func main() {
	client := NewS3()
	spew.Dump(client)
}
