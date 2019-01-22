package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	//log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Test that we can post valid kcode data and retrieve expected blocks
func TestPostBlocks(t *testing.T) {
	fmt.Println("------------ TestPostBlocks ------------")
	//log.SetOutput(ioutil.Discard)
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"method": "blocks",
		},
		Body: `{"source":"<xml xmlns=\"http://www.w3.org/1999/xhtml\"><variables></variables><block type=\"events_onFlick\" id=\"001\" x=\"342\" y=\"266\"><field name=\"TYPE\">up</field><statement name=\"CALLBACK\"><block type=\"objects_setColor\" id=\"002\"><value name=\"TINT\"><shadow type=\"objects_get\" id=\"003\"><field name=\"ID\">all</field></shadow></value><value name=\"TO COLOR\"><shadow type=\"colour_picker\" id=\"004\"><field name=\"COLOUR\">#FF5723</field></shadow></value></block></statement></block></xml>","parts":[],"scene":"owlery"}`,
	}
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: `["events_onFlick","objects_setColor"]`,
	}

	response, err := router(request)
	fmt.Println(response)

	assert.Equal(t, response.Headers, expectedResponse.Headers)
	assert.Contains(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)
}

// Test that we can post valid kcode data and retrieve expected spells
func TestPostSpells(t *testing.T) {
	fmt.Println("------------ TestPostSpells ------------")
	//log.SetOutput(ioutil.Discard)
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"method": "spells",
		},
		Body: `{"source":"<xml xmlns=\"http://www.w3.org/1999/xhtml\"><variables></variables><block type=\"events_onFlick\" id=\"001\" x=\"342\" y=\"266\"><field name=\"TYPE\">up</field><statement name=\"CALLBACK\"><block type=\"objects_setColor\" id=\"002\"><value name=\"TINT\"><shadow type=\"objects_get\" id=\"003\"><field name=\"ID\">all</field></shadow></value><value name=\"TO COLOR\"><shadow type=\"colour_picker\" id=\"004\"><field name=\"COLOUR\">#FF5723</field></shadow></value></block></statement></block></xml>","parts":[],"scene":"owlery"}`,
	}
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: `[]`,
	}

	response, err := router(request)
	fmt.Println(response)

	assert.Equal(t, response.Headers, expectedResponse.Headers)
	assert.Contains(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)
}

// Test that we can post valid kcode data and retrieve expected validationg
func TestPostValidate(t *testing.T) {
	fmt.Println("------------ TestPostValidate ------------")
	//log.SetOutput(ioutil.Discard)
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"method": "validate",
		},
		Body: `{"source":"<xml xmlns=\"http://www.w3.org/1999/xhtml\"><variables></variables><block type=\"events_onFlick\" id=\"001\" x=\"342\" y=\"266\"><field name=\"TYPE\">up</field><statement name=\"CALLBACK\"><block type=\"objects_setColor\" id=\"002\"><value name=\"TINT\"><shadow type=\"objects_get\" id=\"003\"><field name=\"ID\">all</field></shadow></value><value name=\"TO COLOR\"><shadow type=\"colour_picker\" id=\"004\"><field name=\"COLOUR\">#FF5723</field></shadow></value></block></statement></block></xml>","parts":[],"scene":"owlery"}`,
	}
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: `isValid=true`,
	}

	response, err := router(request)
	fmt.Println(response)

	assert.Equal(t, response.Headers, expectedResponse.Headers)
	assert.Contains(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)
}

// Test that we can post valid kcode data and retrieve expected validationg
func TestInvalidPost(t *testing.T) {
	//log.SetOutput(ioutil.Discard)
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"method": "foobar",
		},
		Body: `{"source":"<xml xmlns=\"http://www.w3.org/1999/xhtml\"><variables></variables><block type=\"events_onFlick\" id=\"001\" x=\"342\" y=\"266\"><field name=\"TYPE\">up</field><statement name=\"CALLBACK\"><block type=\"objects_setColor\" id=\"002\"><value name=\"TINT\"><shadow type=\"objects_get\" id=\"003\"><field name=\"ID\">all</field></shadow></value><value name=\"TO COLOR\"><shadow type=\"colour_picker\" id=\"004\"><field name=\"COLOUR\">#FF5723</field></shadow></value></block></statement></block></xml>","parts":[],"scene":"owlery"}`,
	}
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 405,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: `Method Not Allowed`,
	}

	response, err := router(request)
	fmt.Println(response)

	assert.Equal(t, response.Headers, expectedResponse.Headers)
	assert.Contains(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)
}

func TestGet(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"method": "validate",
		},
		Body: "Something",
	}
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: "",
	}

	response, err := router(request)
	//fmt.Println(response)

	//assert.Equal(t, response.Headers, expectedResponse.Headers)
	assert.Contains(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)
}
