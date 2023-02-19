package user

import (
	"encoding/json"
	"errors"
	"src/go-serverless/pkg/validators"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorFailedToFetchRecord = "failed to fetch record"
	ErrorFailedToUnMarshallRecord = "failed to UnMarshall record"
	ErrorInvalidEmail = "Invalid Email"
	ErrorCouldNotDynamoPutItem = "could not dynamo put item"
	ErrorUserAreadyExists = "This user already exists"
	ErrorCouldNotMarshallItem = "could not marshall item"
	ErrorUserDoesNotExist = "User does not exist"
	ErrorCouldNotDeleteItem = "Could not delete item"
)

//user struct
type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password string `json:"password"`
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

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {
	var u User

	err := json.Unmarshal([]byte(req.Body), &u); 
	if err != nil {
		return nil,  errors.New(ErrorFailedToUnMarshallRecord )
	}

	if !validators.IsEmailValid(u.Email){
		return nil, errors.New(ErrorInvalidEmail)
	}

	//check if user exists
	currentUser,_ := Fetchuser(u.Email, tableName, dynaClient)
	if currentUser != nil {
		return nil, errors.New(ErrorUserAreadyExists)
	}

	//Hash Password
	hashed, _ := HashPassword(u.Password)

	u.Password = string(hashed)
	
	// add item 
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshallItem)
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: av,
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := Fetchuser(u.Email, tableName, dynaClient)
	if currentUser == nil {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	av, err :=  dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshallItem)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: av,
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {

	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)

	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil 
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}