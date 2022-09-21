package main

import (
	"context"
	"log"
	"os"

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
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("minio"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               "http://localhost:9000",
			SigningRegion:     "us-east-1",
			HostnameImmutable: true,
		}, nil
	})

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

func (client *S3) Upload(bucket string, pathToFile string, key string) {
	stat, err := os.Stat(pathToFile)
	if err != nil {
		panic("Could not stat image " + err.Error())
	}

	file, err := os.Open(pathToFile)
	if err != nil {
		panic("Could not open local file " + err.Error())
	}

	_, err = client.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: stat.Size(),
	})

	file.Close()

	if err != nil {
		panic("Could not upload file: " + err.Error())
	}
}

func main() {
	client := NewS3()
	//	client.CreateBucket("theo")
	spew.Dump(client)
}
