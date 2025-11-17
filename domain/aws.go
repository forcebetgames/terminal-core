package domain

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type AWSConfig struct {
	SQS *sqs.Client
}

func NewAWSConfig(AWSAccessKey string, secretKey string, region string) *AWSConfig {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("sa-east-1"), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AWSAccessKey, secretKey, "")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	sqs := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.Region = region
	})
	if sqs == nil {
		panic("Unatable to instanciate aws SQS service")
	}

	return &AWSConfig{
		SQS: sqs,
	}
}
