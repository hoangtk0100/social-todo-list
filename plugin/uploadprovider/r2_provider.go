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
	"github.com/hoangtk0100/social-todo-list/common"
)

type r2Provider struct {
	bucketName string
	region     string
	accessKey  string
	secretKey  string
	endPoint   string
	domain     string
	client     *s3.Client
}

func NewR2Provider(bucketName, region, accessKey, secretKey, endPoint, domain string) *r2Provider {
	provider := &r2Provider{
		bucketName: bucketName,
		region:     region,
		accessKey:  accessKey,
		secretKey:  secretKey,
		endPoint:   endPoint,
		domain:     domain,
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: provider.endPoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(provider.accessKey, provider.secretKey, "")),
	)

	if err != nil {
		log.Fatalln(err)
	}

	provider.client = s3.NewFromConfig(cfg)

	return provider
}

func (provider *r2Provider) SaveUploadedFile(ctx context.Context, data []byte, dst string, contentType string) (*common.Image, error) {
	fileBytes := bytes.NewReader(data)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(provider.bucketName),
		Key:         aws.String(dst),
		Body:        fileBytes,
		ContentType: aws.String(contentType),
	}

	_, err := provider.client.PutObject(context.TODO(), params)
	if err != nil {
		return nil, err
	}

	img := &common.Image{
		Url:       fmt.Sprintf("%s/%s", provider.domain, dst),
		CloudName: "r2",
	}

	return img, nil
}

func (provider *r2Provider) RemoveUploadedFile(ctx context.Context, dst string) error {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(provider.bucketName),
		Key:    aws.String(dst),
	}
	_, err := provider.client.DeleteObject(context.TODO(), params)
	if err != nil {
		return err
	}

	return nil
}
