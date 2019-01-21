package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	kcode "github.com/KanoComputing/go/kcode"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	// handlerName for logging purposes
	handlerName = "KCode Web Service"
	// sub method constants
	spellsMethod   = "spells"
	blocksMethod   = "blocks"
	validateMethod = "validate"
)

var (
	logger = log.WithField("handler", handlerName)
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

// Request is an alias to APIGatewayProxyRequest type
type Request = events.APIGatewayProxyRequest

// Response is an alias to the APIGatewayProxyResponse
type Response = events.APIGatewayProxyResponse

// Payload is the struct of expected payload body from request
type Payload struct {
	Method string `json:"method"`
	Body   string `json:"data"`
	Count  int    `json:"count"`
}

func makeResponse(status int, err error, result interface{}) (Response, error) {
	response := Response{
		StatusCode: status,
		Body:       http.StatusText(status),
	}
	// If there is any errors, include error string representation into the response Body
	if err != nil {
		response.Body = fmt.Sprintf("%s - %v", response.Body, err)
	} else {
		body, err := json.Marshal(result)
		if err != nil {
			response.Body = fmt.Sprintf("internal result response error  - %v", err)
		} else {
			response.Body = string(body)
		}
	}
	return response, nil
}

func validateSubMethod(method string) error {
	switch method {
	case spellsMethod, blocksMethod, validateMethod:
		return nil
	default:
		return fmt.Errorf("allowed methods are: speels, blocks or validate")
	}
}

// validatePayload validate the payload in string and return a parsed payload bytes alongside any error
func validatePayload(body string) ([]byte, error) {
	payload := Payload{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return nil, errors.New("invalid JSON payload")
	}
	// marshal payload back to JSON bytes for use in kcode.METHODS
	out, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("internal Payload struct error (%v)", err)
	}
	return out, nil
}

// kcodeHandler will parse input payload and return either spells or blocks
func kcodeHandler(input []byte, submethod string) []string {
	kc := kcode.ExtractKcode(input)
	spells, blocks := kcode.ProcessKcodeString(kc, submethod == spellsMethod, submethod == blocksMethod, true)
	switch submethod {
	case spellsMethod:
		return spells
	case blocksMethod:
		return blocks
	default:
		return []string{}
	}
}

// router is executed by AWS Lambda in the main function and returns an Amazon API Gateway response object to AWS Lambda.
func router(req Request) (Response, error) {
	funcLogger := logger.WithFields(log.Fields{
		"func":    "router",
		"mehtod":  req.HTTPMethod,
		"payload": req.Body,
		"isb64":   req.IsBase64Encoded,
	})
	funcLogger.Infof("request received")
	// defer function will be called before function exits
	defer funcLogger.Infof("response sent")

	if req.HTTPMethod == "POST" {
		submethod := req.QueryStringParameters["method"]
		if err := validateSubMethod(submethod); err != nil {
			return makeResponse(http.StatusMethodNotAllowed, err, nil)
		}
		payload, err := validatePayload(req.Body)
		if err != nil {
			return makeResponse(http.StatusBadRequest, err, nil)
		}
		result := kcodeHandler(payload, submethod)
		return makeResponse(http.StatusOK, nil, result)
	}
	// validate or GET method will currently return just a 200 OK reponse
	return testGetRequest(req)
}

func testGetRequest(req Request) (Response, error) {
	return makeResponse(200, nil, "ok")
}

func main() {
	lambda.Start(router)
}
