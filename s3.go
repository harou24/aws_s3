package s3

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3 struct {
	client *s3.Client
}

func NewS3() *S3 {
	cfg, err := config.LoadDefaultConfig(context.TODO())
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

func (client *S3) UploadFile(bucket string, pathToFile string, key string) {
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

func (client *S3) UploadObj(bucket string, obj []byte, key string) {
	_, err := client.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewBuffer(obj),
	})

	if err != nil {
		panic("Could not upload object " + err.Error())
	}
}

func (client *S3) DeleteObject(bucket string, key string) {
	_, err := client.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		panic("Could not delete object: " + err.Error())
	}
}

func (client *S3) ListBuckets() {
	list, err := client.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	if err != nil {
		panic("Could not list buckets " + err.Error())
	}

	for _, bucket := range list.Buckets {
		fmt.Println("Bucket: ", *bucket.Name)
	}
}

func (client *S3) DeleteEmptyBucket(name string) {
	_, err := client.client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		panic("Could not delete bucket " + err.Error())
	}
}

func (client *S3) ListObjects(bucket string, prefix string) {
	list, err := client.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		panic("Could not list objects " + err.Error())
	}
	for _, obj := range list.Contents {
		fmt.Println("Object: ", obj)
	}
}
