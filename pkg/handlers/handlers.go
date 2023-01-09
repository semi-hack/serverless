package handlers

import (
	"net/http"
	"src/go-serverless/pkg/handlers/user"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ErrorMethodNotAllowed = "Method Not Allowed"

type ErrorBody struct{
	ErrorMsg *string `json:"error,omitempty`
}


func GetUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error) {
		email := req.QueryStringParameters["email"]
		if len(email) > 0 {
			result, err := user.Fetchuser(email, tableName, dynaClient)
			if err != nil {
				return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
			}
			return apiResponse(http.StatusOK, result)
		}

		result, err := user.FetchUsers(tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
}


func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error){

}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error
){

}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI){

}

func UnhandleMethod()(*events.APIGatewayProxyResponse, error){
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}