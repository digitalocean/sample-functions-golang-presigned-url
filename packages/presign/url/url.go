package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Request struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

var (
	key, secret, bucket, region string
	ErrNoFilename               = errors.New("no filename provided")
	ErrNoReq                    = errors.New("no request type provided")
)

const (
	RequestTypeGet = "GET"
	RequestTypePut = "PUT"
)

func init() {
	key = os.Getenv("SPACES_KEY")
	if key == "" {
		panic("no key provided")
	}
	secret = os.Getenv("SPACES_SECRET")
	if secret == "" {
		panic("no secret provided")
	}
	bucket = os.Getenv("BUCKET")
	if bucket == "" {
		panic("no bucket provided")
	}
	region = os.Getenv("REGION")
	if region == "" {
		panic("no region provided")
	}
}

func Main(in Request) (*Response, error) {
	if in.Filename == "" {
		//return &Response{StatusCode: http.StatusBadRequest}, ErrNoFilename
		in.Filename = "new-file.txt"
	}

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(fmt.Sprintf("%s.digitaloceanspaces.com:443", region)),
		Region:      aws.String(region),
	}

	sess := session.New(config)

	if in.Type == "" {
		in.Type = "PUT"
	}
	var url string
	var err error
	switch in.Type {
	case RequestTypeGet:
		url, err = downloadURL(sess, bucket, in.Filename)
		if err != nil {
			return &Response{StatusCode: http.StatusBadRequest}, err
		}
	case RequestTypePut:
		url, err = uploadURL(sess, bucket, in.Filename)
		if err != nil {
			return &Response{StatusCode: http.StatusBadRequest}, err
		}
	default:
		return &Response{StatusCode: http.StatusBadRequest}, ErrNoReq
	}

	return &Response{
		Body: url,
	}, nil
}

func uploadURL(sess *session.Session, bucket string, filename string) (string, error) {
	client := s3.New(sess)
	req, _ := client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		return "", err
	}
	return url, nil
}

func downloadURL(sess *session.Session, bucket string, filename string) (string, error) {
	client := s3.New(sess)
	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		return "", err
	}
	return url, nil
}

// To get a url:
// curl -X PUT -H 'Content-Type: application/json' {your-DO-app-url} -d '{"filename":"{filename}", "type":"GET or PUT"}'

// To Upload or Download the file:
// curl -X PUT -d 'The contents of the file.' "{url}"
