package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

var (
	LinkTableName = os.Getenv("LINK_TABLE")
	Region        = os.Getenv("REGION")
)

type DB struct {
	Instance *dynamodb.DynamoDB
}

type Link struct {
	ShortenResource string `json:"shorten_resource"`
	OriginalURL     string `json:"original_url"`
}

func New() DB {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(Region)}),
	)

	return DB{Instance: dynamodb.New(sess)}
}

func (d DB) GetItem(i interface{}) (string, error) {
	item, err := d.Instance.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(LinkTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"shorten_resource": {
				S: aws.String(i.(string)),
			},
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to get item")
	}
	if item.Item == nil {
		return "", nil
	}

	link := Link{}
	err = dynamodbattribute.UnmarshalMap(item.Item, &link)
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal item")
	}

	return link.OriginalURL, nil
}

func (d DB) PutItem(i interface{}) (interface{}, error) {
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(LinkTableName),
	}
	item, err := d.Instance.PutItem(input)
	if err != nil {
		return nil, err
	}

	return item, nil
}
