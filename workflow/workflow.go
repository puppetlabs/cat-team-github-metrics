package main

import (
	"context"
	"fmt"
	"os"

	collar "github.com/chelnak/collar/pkg/modules"
	"github.com/chelnak/relay-workflow-builder/pkg/workflow"
)

const (
	moduleOwner  = "puppetlabs"
	imageName    = "ghcr.io/chelnak/cat-team-github-metrics:latest"
	scheduleCron = "0 0 * * *"
	scheduleType = "schedule"
)

func main() {
	modules, err := getSupportedModules()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a new workflow
	w := workflow.NewWorkflow("A workflow for collecting GitHub metrics.")

	// Add a schedule trigger
	w.AddTrigger(
		workflow.Trigger{
			Name: "schedule",
			Source: map[string]string{
				"type":     scheduleType,
				"schedule": scheduleCron,
			},
		},
	)

	// Add an export step for each module
	for _, module := range *modules {
		w.AddStep(
			workflow.Step{
				Name:  module.Name,
				Image: imageName,
				Spec: map[string]string{
					"connection":           "${connections.gcp.'content-and-tooling-lab'}",
					"repo_owner":           moduleOwner,
					"repo_name":            module.Name,
					"github_token":         "${secrets.GITHUB_TOKEN}",
					"big_query_project_id": "${secrets.BIG_QUERY_PROJECT_ID}",
					"command":              "export",
				},
			},
		)
	}

	// Get the names of all of the current steps
	var dependsOn []string
	for _, step := range w.GetSteps() {
		dependsOn = append(dependsOn, step.Name)
	}

	// Add a stamp step. This will update the stamp table with the time of the last successful run.
	// Note that DependsOn is set here and uses the names of all of the steps in the workflow prior to this step.
	// The stamp step also holds extra spec requirements. This will be fixed in a future refactor.
	w.AddStep(
		workflow.Step{
			Name:      "Successful run timestamp",
			Image:     imageName,
			DependsOn: dependsOn,
			Spec: map[string]string{
				"connection":           "${connections.gcp.'content-and-tooling-lab'}",
				"repo_owner":           "not_used",
				"repo_name":            "not_used",
				"github_token":         "not_used",
				"big_query_project_id": "${secrets.BIG_QUERY_PROJECT_ID}",
				"command":              "stamp",
			},
		},
	)

	// Finally, write the workflow to a file
	err = w.Write(nil)
	if err != nil {
		if err != workflow.ErrValidation {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}

func getSupportedModules() (*[]collar.Module, error) {
	client := collar.NewModuleClient(nil, "")
	ctx := context.Background()
	return client.GetSupportedModules(ctx)
}
