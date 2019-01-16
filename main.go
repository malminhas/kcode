package main

import (
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	index, err := ioutil.ReadFile("public/index.html")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	/*
	From: https://docs.aws.amazon.com/lambda/latest/dg/with-on-demand-https-create-package.html#with-apigateway-example-deployment-pkg-go
	func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    		fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestId)
    		fmt.Printf("Body size = %d.\n", len(request.Body))
    		fmt.Println("Headers:")
		for key, value := range request.Headers {
        		fmt.Printf("    %s: %s\n", key, value)
    		}
    	return events.APIGatewayProxyResponse { Body: request.Body, StatusCode: 200 }, nil
	}
	*/
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(index),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil

}

func main() {
	lambda.Start(Handler)
}
