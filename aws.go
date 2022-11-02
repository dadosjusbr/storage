package storage

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type S3Client struct {
	s3       *s3.S3
	newrelic *newrelic.Application
	bucket   string
}

func NewS3Client(region string, bucket string) (*S3Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating aws session: %w", err)
	}
	s3Client := s3.New(sess)
	return &S3Client{s3: s3Client, bucket: bucket}, nil
}

func (s S3Client) UploadFile(srcPath string, dstFolder string) (*Backup, error) {
	txn := s.newrelic.StartTransaction("aws.UploadFile")
	defer txn.End()
	ctx := newrelic.NewContext(aws.BackgroundContext(), txn)
	uploader := s3manager.NewUploaderWithClient(s.s3)

	f, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file at %s: %v", srcPath, err)
	}
	defer f.Close()

	_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Body:   f,
		Key:    aws.String(dstFolder),
	})
	if err != nil {
		return nil, fmt.Errorf("Error trying to upload file in S3 with key (%s): %v", dstFolder, err)
	}

	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(dstFolder),
	}
	headObjectOutput, err := s.s3.HeadObject(headObjectInput)
	if err != nil {
		return nil, fmt.Errorf("Error getting file metadata from (%s): %q", dstFolder, err)
	}
	backup := &Backup{
		Size: *headObjectOutput.ContentLength,
		Hash: strings.ReplaceAll(*headObjectOutput.ETag, "\"", ""),
		URL:  fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, dstFolder),
	}

	return backup, nil
}

func (s S3Client) Backup(Files []string, dstFolder string) ([]Backup, error) {
	if len(Files) == 0 {
		return []Backup{}, nil
	}
	var backups []Backup
	txn := s.newrelic.StartTransaction("aws.Backup")
	defer txn.End()
	ctx := newrelic.NewContext(aws.BackgroundContext(), txn)
	uploader := s3manager.NewUploaderWithClient(s.s3)

	for _, file := range Files {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("Error opening file at %s: %v", file, err)
		}
		defer f.Close()

		_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
			Bucket: aws.String(s.bucket),
			Body:   f,
			Key:    aws.String(dstFolder),
		})
		if err != nil {
			return nil, fmt.Errorf("Error trying to upload file in S3 with key (%s): %v", dstFolder, err)
		}

		headObjectInput := &s3.HeadObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(dstFolder),
		}

		headObjectOutput, err := s.s3.HeadObject(headObjectInput)
		if err != nil {
			log.Fatalf("Error getting file metadata from (%s): %q", dstFolder, err)
		}

		backup := &Backup{
			Size: *headObjectOutput.ContentLength,
			Hash: strings.ReplaceAll(*headObjectOutput.ETag, "\"", ""),
			URL:  fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, dstFolder),
		}
		backups = append(backups, *backup)
	}
	return backups, nil
}
