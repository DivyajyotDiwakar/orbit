package db

import (
	"Orbit/configs"
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	r2client *s3.Client
	r2Once   sync.Once
)

func InitR2() {
	R2Obj := configs.LoadConfig().R2
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				R2Obj.Access_Key, R2Obj.Secret_Key, "",
			),
		),
	)

	if err != nil {
		panic("Failed to load R2 config: " + err.Error())
	}

	R2Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(R2Obj.Endpoint)
		o.UsePathStyle = true
	})

	r2client = R2Client
}

func GetR2Client() *s3.Client {
	r2Once.Do(func() {
		InitR2()
	})

	return r2client
}
