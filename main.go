package main

import (
    "encoding/json"
	"fmt"
    "net/http"
	kcode "github.com/KanoComputing/go/kcode"
    
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type resp struct {
    Method  string `json:"method"`
    Body  string `json:"data"`
    Count int `json:"count"`
}

func dumpHTTPRequest(req events.APIGatewayProxyRequest) {
	method := req.HTTPMethod
	params := req.QueryStringParameters
	hdrs := req.Headers
	body := req.Body
	b64 := req.IsBase64Encoded
	fmt.Printf("::dumpHTTPRequest::\n\t%s\n\tHeaders=%s\n\tQueryStringParameters=%s\n\tBody=%s (%d bytes)\n\tIsBase64:%t\n",
	    method,hdrs,params,body,len(body),b64)
}

// router is executed by AWS Lambda in the main function and returns an Amazon API Gateway response object to AWS Lambda.
func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dumpHTTPRequest(req)
    method := req.QueryStringParameters["method"]
    switch req.HTTPMethod {
	case "GET":
		return testGetRequest(req)
    case "POST":
		switch method {
		case "spells":
			return findBlocksRequest(req)
		case "blocks":
			return findBlocksRequest(req)
		case "validate":
			return validateRequest(req)
		default:
			return clientError(http.StatusBadRequest)
		}
    default:
        return clientError(http.StatusMethodNotAllowed)
    }
}

func getError() (err error) {
   err = nil
   return
}

// Looking for blocks in .kcode input
func findBlocksRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// https://docs.aws.amazon.com/lambda/latest/dg/with-on-demand-https-create-package.html#with-apigateway-example-deployment-pkg-go

    body := req.Body
	// You might want to add CORS header here:
	// resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	//s := fmt.Sprintf("{\"method\":\"%s\",\"body\": \"Hello World\",\"count\": 1}", method)
	
	err := getError()
	
    if err != nil {
        return serverError(err)
    }
    //if resp == nil {
    //    return clientError(http.StatusNotFound)
    //}

	verbose := false
	filename := ""
	kc := kcode.ExtractKcode([]byte(body))
	_, blocks := kcode.ProcessKcodeString(filename, kc, false, true, verbose)
	fmt.Println(blocks)
	resp.Body = blocks[0]

    js, err := json.Marshal(resp)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

// Looking for spells in .kcode input
func findSpellsRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    body := req.Body
	// You might want to add CORS header here:
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	
	resp.Body = body
	
    js, err := json.Marshal(resp)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

// Looking for spells in .kcode input
func validateRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    body := req.Body
	// You might want to add CORS header here:
	// resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"

	resp.Body = body

    js, err := json.Marshal(resp)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

func testGetRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    body := req.Body
	// You might want to add CORS header here:
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	
	/*
	verbose := false
	filename := ""
	b := []byte(body)
	_, blocks := kcode.ProcessKcodeString(filename, b, false, true, verbose)
	fmt.Println(blocks)
	resp.Body = blocks[0]
	*/

	resp.Body = body
	
    js, err := json.Marshal(resp)
    if err != nil {
        return serverError(err)
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
    //var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
	//errorLogger.Println(err.Error())
	fmt.Println(err.Error())
	
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body:       http.StatusText(status),
    }, nil
}

func main() {
	lambda.Start(router)
}
