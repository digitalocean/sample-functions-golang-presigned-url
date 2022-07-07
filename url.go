// export SPACES_KEY=XXXXXXXX && export SPACES_SECRET=XXXXXXXXXXXXX
package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Request struct {
	Filename string `json:"filename"`
	Region   string `json:"region"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

func Main(in Request) (*Response, error) {
	key := os.Getenv("SPACES_KEY")
	secret := os.Getenv("SPACES_SECRET")
	bucket := os.Getenv("BUCKET")

	if in.Filename == "" {
		fmt.Println("No Filename Entered")
	}
	if in.Region == "" {
		fmt.Println("No Region Entered")
	}

	region, err := checkRegion(in.Region)
	if err != nil {
		fmt.Println(err)
	}

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("https://" + region + ".digitaloceanspaces.com"),
		Region:      aws.String("us-east-1"),
	}

	sess := session.New(config)
	url, err := UploadURL(sess, bucket, in.Filename)
	if err != nil {
		fmt.Println("Error retrieving URL: ", err)
	}

	return &Response{
		Body: fmt.Sprintf("The presigned URL: %s", url),
	}, nil
}

func UploadURL(sess *session.Session, bucket string, filename string) (string, error) {
	client := s3.New(sess)
	req, _ := client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	url, err := req.Presign(10 * time.Minute)
	if err != nil {
		return "", err
	}
	return url, nil
}

func checkRegion(region string) (string, error) {
	if region == "San Francisco" || region == "san francisco" || region == "sfo3" {
		region = "sfo3"
	} else if region == "Frankfurt" || region == "frankfurt" || region == "fra1" {
		region = "fra1"
	} else if region == "Amsterdam" || region == "amsterdam" || region == "ams3" {
		region = "ams3"
	} else if region == "New York" || region == "new york" || region == "nyc3" {
		region = "nyc3"
	} else if region == "Singapore" || region == "singapore" || region == "sgp1" {
		region = "sgp1"
	} else {
		return "", errors.New("Invalid Region Given.")
	}
	return region, nil
}

// Once you get the url outputed: run this command in terminal
//curl -X PUT \
//-H "Content-Type: text" \
//-d "The contents of the file." \
// enter presigned url here in "" : "https://slack.nyc3.digitaloceanspaces.com/"
