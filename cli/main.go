package main

import (
	"github.com/openshift-ai/mpijob/cli/cmd"
	"github.com/openshift-ai/mpijob/cli/cmd/training"
)

func main() {
	// Add training command to root
	cmd.RootCmd.AddCommand(training.TrainingCmd)

	// Execute
	cmd.Execute()
}