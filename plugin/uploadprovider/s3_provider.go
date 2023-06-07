package uploadprovider

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hoangtk0100/social-todo-list/common"
)

type s3Provider struct {
	bucketName string
	region     string
	accessKey  string
	secretKey  string
	endPoint   string
	domain     string
	client     *s3.Client
}

func NewS3Provider(bucketName, region, accessKey, secretKey, endPoint, domain string) *s3Provider {
	provider := &s3Provider{
		bucketName: bucketName,
		region:     region,
		accessKey:  accessKey,
		secretKey:  secretKey,
		endPoint:   endPoint,
		domain:     domain,
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(provider.region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(provider.accessKey, provider.secretKey, "")),
	)

	if err != nil {
		log.Fatalln(err)
	}

	provider.client = s3.NewFromConfig(cfg)

	return provider
}

func (provider *s3Provider) SaveUploadedFile(ctx context.Context, data []byte, dst string, contentType string) (*common.Image, error) {
	fileBytes := bytes.NewReader(data)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(provider.bucketName),
		Key:         aws.String(dst),
		Body:        fileBytes,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACL(types.BucketCannedACLPrivate),
	}

	_, err := provider.client.PutObject(context.TODO(), params)
	if err != nil {
		return nil, err
	}

	img := &common.Image{
		Url:       fmt.Sprintf("%s/%s", provider.domain, dst),
		CloudName: "s3",
	}

	return img, nil
}
