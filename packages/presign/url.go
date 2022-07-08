// export SPACES_KEY=XXXXXXXX && export SPACES_SECRET=XXXXXXXXXXXXX
package presign

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
	key, secret, bucket, region, url string
	ErrNoFilename                    = errors.New("no filename provided")
	ErrNoReq                         = errors.New("no request type provided")
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
	region = os.Getenv("Region")
	if region == "" {
		panic("no region provided")
	}
}

func Main(in Request) (*Response, error) {
	if in.Filename == "" {
		return &Response{StatusCode: http.StatusBadRequest}, ErrNoFilename
	}
	reg, err := checkRegion(region)
	if err != nil {
		fmt.Println(err)
		return &Response{StatusCode: http.StatusBadRequest}, err
	}

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("https://" + reg + ".digitaloceanspaces.com"),
		Region:      aws.String("us-east-1"),
	}

	sess := session.New(config)

	if in.Type == "PUT" {
		url, err = uploadURL(sess, bucket, in.Filename)
		if err != nil {
			return &Response{StatusCode: http.StatusBadRequest}, err
		}
	} else if in.Type == "GET" {
		url, err = downloadURL(sess, bucket, in.Filename)
		if err != nil {
			return &Response{StatusCode: http.StatusBadRequest}, err
		}
	} else {
		return &Response{StatusCode: http.StatusBadRequest}, ErrNoReq
	}

	return &Response{
		Body: fmt.Sprintf("The presigned URL: %s", url),
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
		return "", errors.New("invalid region given")
	}
	return region, nil
}

// Once you get the url outputed: run this command in terminal
//curl -X PUT \
//-H "Content-Type: text" \
//-d "The contents of the file." \
// enter presigned url here in "" : "https://slack.nyc3.digitaloceanspaces.com/"
