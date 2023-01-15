package user

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorFailedToUnMarshallRecord = "failed to UnMarshall record"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}


func Fetchuser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnMarshallRecord)
	}
	return item, nil
}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnMarshallRecord)
	}
	return item, nil

}

func CreateUser(data User, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {
	var u User
	return &u, nil
}
