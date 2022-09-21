package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/davecgh/go-spew/spew"
)

type S3 struct {
	client *s3.Client
}

func NewS3() *S3 {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               "http://localhost:9000",
			SigningRegion:     "us-east-1",
			HostnameImmutable: true,
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("minio"), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		panic(err)
	}
	spew.Dump(cfg)
	cl := s3.NewFromConfig(cfg)
	return &S3{
		client: cl,
	}
}

func (client *S3) CreateBucket(name string) {
	_, err := client.client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket:                    aws.String(name),
		ACL:                       types.BucketCannedACLPrivate,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraintUsWest2},
	})

	if err != nil {
		panic("could not create bucket: " + err.Error())
	}
}

func main() {
	client := NewS3()
	client.CreateBucket("hello")
	spew.Dump(client)
}
