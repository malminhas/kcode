package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"reflect"

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
	// some error values
	errInvalidJSONPayload = errors.New("invalid JSON payload")
	errInvalidSubmethod   = errors.New("allowed methods are: spells, blocks or validate")
	errInvalidHTTPRequest = errors.New("allowed HTTP requests are: GET, POST")
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

// Request is an alias to APIGatewayProxyRequest type
type Request = events.APIGatewayProxyRequest

// Response is an alias to the APIGatewayProxyResponse
type Response = events.APIGatewayProxyResponse

func makeResponse(status int, err error, result interface{}) (Response, error) {
	response := Response{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: http.StatusText(status),
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
		return errInvalidSubmethod
	}
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// validatePayload validate the payload in string and return a parsed payload bytes alongside any error
func validatePayload(body string) ([]byte, error) {
	//fmt.Println("---------- body ---------------")
	//fmt.Println(body)
	//fmt.Println("---------- body ---------------")
	//fmt.Println(typeof(body))

	payload := kcode.KCode{}
	//payload := kcode.KCode{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return nil, errInvalidJSONPayload
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
	verbose := false
	kc, _ := kcode.ExtractKcode(input)
	// NOTE: if both spells and blocks are just []string,
	// we can simply append/concat them regardless and return it, like:
	// `return append(spells, blocks...)`
	switch submethod {
	case spellsMethod:
		spells, _ := kcode.ProcessKcodeString(kc, true, false, verbose)
		return spells
	case blocksMethod:
		_, blocks := kcode.ProcessKcodeString(kc, false, true, verbose)
		return blocks
	case validateMethod:
		expectedSpells, foundSpells, expectedBlocks, foundBlocks, valid := kcode.ValidateBlocksAndSpellsString(kc, verbose)
		r := []string{}
		r = append(r, fmt.Sprintf("expectedSpells=%d", expectedSpells))
		r = append(r, fmt.Sprintf("foundSpells=%d", foundSpells))
		r = append(r, fmt.Sprintf("expectedBlocks=%d", expectedBlocks))
		r = append(r, fmt.Sprintf("foundBlocks=%d", foundBlocks))
		r = append(r, fmt.Sprintf("isValid=%t", valid))
		return r
	default:
		return []string{}
	}
}

// router is executed by AWS Lambda in the main function and returns an Amazon API Gateway response object to AWS Lambda.
func router(req Request) (Response, error) {
	funcLogger := logger.WithFields(log.Fields{
		"func":    "router",
		"method":  req.HTTPMethod,
		"payload": req.Body,
		"isb64":   req.IsBase64Encoded,
	})
	funcLogger.Infof("request received")
	// defer function will be called before function exits
	//defer funcLogger.Infof("response sent")

	switch req.HTTPMethod {
	case "POST":
		// get sub-method on POST from "method" querystring parameter
		submethod := req.QueryStringParameters["method"]
		if err := validateSubMethod(submethod); err != nil {
			return makeResponse(http.StatusMethodNotAllowed, err, nil)
		}
		// parse and turn string payload to bytes payload and validate it
		payload, err := validatePayload(req.Body)
		if err != nil {
			return makeResponse(http.StatusBadRequest, err, nil)
		}
		// result is of either "spells" or "blocks" or "validate" in []string values
		result := kcodeHandler(payload, submethod)
		return makeResponse(http.StatusOK, nil, result)
	case "GET":
		// GET method will currently return just a 200 OK response
		return testGetRequest(req)
	default:
		return makeResponse(http.StatusBadRequest, errInvalidHTTPRequest, "HTTP Request not supported")
	}
}

func testGetRequest(req Request) (Response, error) {
	// FIXME: this is a dummy, please fix me
	return makeResponse(200, nil, "ok")
}

func main() {
	lambda.Start(router)
}
