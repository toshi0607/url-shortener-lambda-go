package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestNew() DB {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(Region),
		Endpoint: aws.String("http://localhost:8000")}),
	)

	return DB{Instance: dynamodb.New(sess)}
}

func (d DB) CreateLinkTable() error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("shorten_resource"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("shorten_resource"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(LinkTableName),
	}

	_, err := d.Instance.CreateTable(input)
	if err != nil {
		return err
	}
	return nil
}

func (d DB) DeleteLinkTable() error {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(LinkTableName),
	}
	_, err := d.Instance.DeleteTable(input)
	if err != nil {
		return err
	}
	return nil
}
