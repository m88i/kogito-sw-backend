package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleSWDeployment(t *testing.T) {
	workflowSource, err := ioutil.ReadFile("testdata/greetingevent.sw.json")
	assert.NoError(t, err)
	assert.NotEmpty(t, workflowSource)
	server := httptest.NewServer(defaultHandler(handleSWDeployment))
	defer server.Close()

	res, err := http.Post(server.URL, "application/json", bytes.NewReader(workflowSource))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var reply successDeployment
	err = json.NewDecoder(res.Body).Decode(&reply)
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, "Event Based Greeting Workflow", reply.WorkflowName)
}
