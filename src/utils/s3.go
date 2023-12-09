package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Instance struct {
	s3BucketName string
	s3BucketKey  string
	s3BucketURL  string
	s3Region     string
	accessKeyId  string
	secretKey    string
	ACL          string
	client       *s3.S3
}

var awsS3 *s3Instance

// CreateS3Session creates a single s3 session once that can be reused across the application
func CreateS3Session(bucketName, bucketKey, bucketURL, s3Region, accessKey, secretKey, ACL string) {
	var once sync.Once
	if awsS3 == nil {
		once.Do(
			func() {
				awsS3 = &s3Instance{
					s3BucketName: bucketName,
					s3BucketKey:  bucketKey,
					s3BucketURL:  bucketURL,
					s3Region:     s3Region,
					accessKeyId:  accessKey,
					secretKey:    secretKey,
					ACL:          ACL,
				}
				s3Session := session.Must(session.NewSession(&aws.Config{
					Credentials: credentials.NewStaticCredentials(
						accessKey,
						secretKey,
						"",
					),
					Region: aws.String(s3Region),
				}))
				awsS3.client = s3.New(s3Session)
			},
		)
	}
}

// UploadToS3 uploads a file to the s3 bucket with the passed file name (file key) and content type
func UploadToS3(file *bytes.Reader, fileKey, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	bucketKey := awsS3.s3BucketKey + "/" + fileKey

	_, err := awsS3.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(awsS3.s3BucketName),
		Key:         aws.String(bucketKey),
		ACL:         aws.String(awsS3.ACL),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		if uploadError, ok := err.(awserr.Error); ok && uploadError.Code() == request.CanceledErrorCode {
			return errors.New("upload canceled due to a timeout")
		}
		return fmt.Errorf("failed to upload the object to s3 - %s", err.Error())
	}

	return nil
}

// DeleteFromS3 deletes a file from the s3 bucket with the passed file name (as a full URL from the DB)
func DeleteFromS3(fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	file := strings.Split(fileName, awsS3.s3BucketURL)[1]
	key := awsS3.s3BucketKey + file

	_, err := awsS3.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(awsS3.s3BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		if deleteError, ok := err.(awserr.Error); ok && deleteError.Code() == request.CanceledErrorCode {
			return errors.New("upload canceled due to a timeout")
		}
		return fmt.Errorf("failed to delete the object from s3 - %s", err.Error())
	}

	return nil
}

// GetTheFullS3BucketURL retrieves the base full s3 bucket URL
func GetTheFullS3BucketURL() string {
	return awsS3.s3BucketURL
}
