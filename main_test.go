package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestPostBlocks(t *testing.T) {

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"method": "validate",
		},
		Body: `{"source": "<xml xmlns="http://www.w3.org/1999/xhtml"><variables></variables><block type="events_onFlick" x="342" y="266"><field name="TYPE">up</field><statement name="CALLBACK"><block type="objects_setColor"><value name="TINT"><shadow type="objects_get"><field name="ID">all</field></shadow></value><value name="TO COLOR"><shadow type="colour_picker"><field name="COLOUR">#FF5723</field></shadow></value></block></statement></block></xml>","parts":[],"scene":"owlery"}`,
	}
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: "events_onFlick",
	}

	response, err := router(request)
	fmt.Println(response)

	//assert.Equal(t, response.Headers, expectedResponse.Headers)
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
	fmt.Println(response)

	//assert.Equal(t, response.Headers, expectedResponse.Headers)
	assert.Contains(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)
}
