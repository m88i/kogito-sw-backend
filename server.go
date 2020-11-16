package main

import (
	"encoding/json"
	"github.com/serverlessworkflow/sdk-go/parser"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	envKeyServerPort  = "SERVER_PORT"
	defaultServerPort = "9000"
)

type defaultHandler func(http.ResponseWriter, *http.Request) *serverError

// ServeHTTP http.ServeHTTP implementation that can handle errors
func (fn defaultHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := fn(writer, request); err != nil {
		handleServerError(err, writer)
	}
}

func getServerPort() string {
	serverPort := os.Getenv(envKeyServerPort)
	if len(serverPort) == 0 {
		serverPort = defaultServerPort
	}
	return serverPort
}

func startServingRequests() {
	http.HandleFunc("/", defaultHandler(handleSWDeployment).ServeHTTP)
	http.HandleFunc("/live", defaultHandler(handleLiveness).ServeHTTP)
	http.HandleFunc("/ready", defaultHandler(handleReady).ServeHTTP)
	serverPort := getServerPort()
	log.Printf("Listening at port %s", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}

type successDeployment struct {
	WorkflowName string `json:"workflowName,omitempty"`
}

func handleSWDeployment(writer http.ResponseWriter, request *http.Request) *serverError {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return newServerError("Failed to read request body", err)
	}
	if len(body) == 0 {
		return newServerError("Body contains no data", nil)
	}
	workflow, err := parser.FromJSONSource(body)
	if err != nil {
		return newServerError("Failed to parse request body into a workflow definition", err)
	}
	if err := deploy(*workflow, body); err != nil {
		return newServerError("Failed to deploy Kogito service", err)
	}
	if err := json.NewEncoder(writer).Encode(&successDeployment{WorkflowName: workflow.Name}); err != nil {
		return newServerError("Failed to encode success response", err)
	}
	return nil
}

func handleLiveness(writer http.ResponseWriter, request *http.Request) *serverError {
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte("ok")); err != nil {
		return newServerError("Failed to write in the response stream", err)
	}
	return nil
}

func handleReady(writer http.ResponseWriter, request *http.Request) *serverError {
	// verify namespace env var
	if len(os.Getenv(envKeyNamespace)) == 0 {
		return newServerError("Environment variable "+envKeyNamespace+" not defined", nil)
	}
	// verify kubectl
	if err := execCmd("oc", "version"); err != nil {
		return newServerError("Failed to run kubectl", err)
	}
	// verify kogito
	if err := execCmd("kogito", "--version"); err != nil {
		return newServerError("Failed to run kogito CLI", err)
	}
	writer.WriteHeader(http.StatusOK)
	return nil
}
