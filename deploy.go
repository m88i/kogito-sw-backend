package main

import (
	"github.com/serverlessworkflow/sdk-go/model"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

const (
	kogitoInstallKnativeBrokerCmdF = "install infra kogito-knative --apiVersion eventing.knative.dev/v1 --kind Broker --resource-name default"

	defaultDataPath            = "/home/kogito/data"
	regularProjectTemplatePath = "/projects_template/regular-sw"
	eventsProjectTemplatePath  = "/projects_template/events-sw"
	resourcesPath              = "/src/main/resources"
	knativeResourcesPath       = "/knative"

	workflowFileSuffix = ".sw.json"

	// exposed via Kubernetes pod
	envKeyNamespace = "DEPLOY_NAMESPACE"
	envKeyDataPath  = "PROJECTS_DATA_PATH"
)

func deploy(workflow model.Workflow, workflowFile []byte) error {
	namespace := os.Getenv(envKeyNamespace)
	projectTemplate := regularProjectTemplatePath
	if err := execCmd("oc", "project", namespace); err != nil {
		log.Error("Failed to set namespace", err)
		return err
	}
	if len(workflow.Events) > 0 {
		projectTemplate = eventsProjectTemplatePath
		if err := deployEventsInfra(namespace); err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(getWorkflowFilePath(workflow, projectTemplate), workflowFile, 0644); err != nil {
		log.Error("Failed to write Workflow file in project's path", err)
		return err
	}
	if err := execCmd("kogito", "deploy", workflow.ID, getDataPath()+projectTemplate); err != nil {
		log.Error("Failed to deploy Kogito service in the target cluster", err)
		return err
	}

	return nil
}

func deployEventsInfra(namespace string) error {
	// apply the knative broker
	if err := execCmd("oc", "apply", "-f", getDataPath()+knativeResourcesPath, "-n", namespace); err != nil {
		log.Error("Failed create Knative Eventing broker", err)
		return err
	}
	// installs kogito infra that uses the same broker
	if err := execCmd("kogito", strings.Split(kogitoInstallKnativeBrokerCmdF, " ")...); err != nil {
		log.Error("Failed create Kogito Infra abstraction for Knative Broker", err)
		return err
	}
	return nil
}

func getWorkflowFilePath(workflow model.Workflow, projectPath string) string {
	return getDataPath() + projectPath + resourcesPath + "/" + workflow.ID + workflowFileSuffix
}

func getDataPath() string {
	dataPath := os.Getenv(envKeyDataPath)
	if len(dataPath) == 0 {
		dataPath = defaultDataPath
	}
	return dataPath
}
