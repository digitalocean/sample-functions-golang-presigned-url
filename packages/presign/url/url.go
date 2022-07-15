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

// Request takes in the user's input for the filename they want and if the type is a GET or PUT.
type Request struct {
	// Filename is the name of the file that will be uploaded or downloaded.
	Filename string `json:"filename"`
	// Type is a presigned request type to "GET" or "PUT" an object.
	Type string `json:"type"`
	// Duration is the duration in which the presigned url will last.
	Duration string `json:"duration"`
}

// Response returns back the http code, type of data, and the presigned url to the user.
type Response struct {
	// StatusCode is the http code that will be returned back to the user.
	StatusCode int `json:"statusCode,omitempty"`
	// Headers is the information about the type of data being returned back.
	Headers map[string]string `json:"headers,omitempty"`
	// Body will contain the presigned url to upload or download files.
	Body string `json:"body,omitempty"`
}

var (
	key, secret, bucket, region string
	// ErrNoFilename will return an error if no filename is provided by the user.
	ErrNoFilename = errors.New("no filename provided")
	// ErrNoFilename will return an error if no request type is provided by the user.
	ErrNoReq = errors.New("no request type provided")
	// ErrNoDuration will return an error if no duration is provided by the user.
	ErrNoDuration = errors.New("no duration provided")
)

const (
	// RequestTypeGet is the presigned request type to download a file.
	RequestTypeGet = "GET"
	// RequestTypePUT is the presigned request type to upload a file.
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

// Main configures a client using the key, secret, and region provided and returns a presigned
// url to upload a file or download a file from a DigitalOcean Space.
func Main(in Request) (*Response, error) {
	if in.Filename == "" {
		return &Response{StatusCode: http.StatusBadRequest}, ErrNoFilename
	}

	duration, err := time.ParseDuration(in.Duration)
	if err != nil {
		return &Response{StatusCode: http.StatusBadRequest}, ErrNoDuration
	}

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(fmt.Sprintf("%s.digitaloceanspaces.com:443", region)),
		Region:      aws.String(region),
	}

	sess := session.New(config)

	var url string
	var err error
	switch in.Type {
	case RequestTypeGet:
		url, err = downloadURL(sess, bucket, in.Filename, duration)
		if err != nil {
			return &Response{StatusCode: http.StatusBadRequest}, err
		}
	case RequestTypePut:
		url, err = uploadURL(sess, bucket, in.Filename, duration)
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

func uploadURL(sess *session.Session, bucket string, filename string, duration tiime.Duration) (string, error) {
	client := s3.New(sess)
	req, _ := client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	url, err := req.Presign(duration)
	if err != nil {
		return "", err
	}
	return url, nil
}

func downloadURL(sess *session.Session, bucket string, filename string, duration tiime.Duration) (string, error) {
	client := s3.New(sess)
	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	url, err := req.Presign(duration)
	if err != nil {
		return "", err
	}
	return url, nil
}
