package s3_test

import (
	"testing"

	s3 "github.com/harou24/aws_s3.git"
)

func TestCreateBucket(t *testing.T) {
	cl := s3.NewS3()
	cl.CreateBucket("hello")
	cl.ListBuckets()
}
