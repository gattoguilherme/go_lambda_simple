package db

import (
	"go-hello-lambda/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

const AWS_REGION = "sa-east-1"
const TABLE_NAME = "go-serverless"

var db = dynamodb.New(session.New(),
	aws.NewConfig().WithRegion(AWS_REGION))

func GetUsers() ([]domain.User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(TABLE_NAME),
	}
	result, err := db.Scan(input)

	if err != nil {
		return []domain.User{}, err
	}

	if len(result.Items) == 0 {
		return []domain.User{}, nil
	}

	var users []domain.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return []domain.User{}, err
	}
	return users, nil
}

func CreateUser(user domain.User) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(uuid.String()),
			},
			"name": {
				S: aws.String(user.Name),
			},
		},
	}
	_, err = db.PutItem(input)
	return err
}
