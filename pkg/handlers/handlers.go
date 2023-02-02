package handlers

import (
	"net/http"
	"src/go-serverless/pkg/handlers/user"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Define Error Methods
var ErrorMethodNotAllowed = "Method Not Allowed"
var ErrorMethodNotAcceptable = "Method Not Acceptable"

type ErrorBody struct{
	ErrorMsg *string `json:"error,omitempty`
}


// GetUser Handler, takes in a request tableName and dynaClient
// 
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
		if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
			return apiResponse(http.StatusNotAcceptable, ErrorMethodNotAcceptable)
		}
		result, err := user.CreateUser(tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error){
		result, err := user.UpdateUser(req, tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(req, tableName, dynaClient)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, nil)

}

func UnhandleMethod()(*events.APIGatewayProxyResponse, error){
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
