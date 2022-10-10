package token

import (
	"Weather-Bot-Discord/mylogger"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"os"
	"strings"
)

func GetToken() (string, error) {
	cloud := os.Getenv("CLOUD")
	if strings.ToLower(cloud) == "aws" {
		secretId := os.Getenv("AWS_SECRET_ID")
		region := os.Getenv("AWS_REGION")
		if secretId == "" || region == "" {
			return "", errors.New("specified aws as cloud but aws secret id or aws region were not given")
		}

		token, err := getTokenFromAws(secretId, region)
		if err == nil {
			mylogger.Println("Obtained bot token from AWS")
		}
		return token, err
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return "", errors.New("DISCORD_TOKEN environment variable is empty")
	}
	mylogger.Println("Obtained bot token from env variable")
	return token, nil
}

func getTokenFromAws(secretId, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", err
	}

	client := secretsmanager.NewFromConfig(cfg)

	secretValueOutput, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	})
	if err != nil {
		return "", err
	}

	if secretValueOutput.SecretString == nil {
		return "", errors.New("aws secret value is empty")
	}

	return *secretValueOutput.SecretString, nil
}
