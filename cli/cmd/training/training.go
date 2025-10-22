package training

import (
	"github.com/spf13/cobra"
)

// TrainingCmd represents the training command
var TrainingCmd = &cobra.Command{
	Use:   "training",
	Short: "Manage training jobs",
	Long: `Manage training jobs in OpenShift AI.
This command allows you to create, list, describe, and delete training jobs.`,
}

func init() {
	// Add subcommands
	TrainingCmd.AddCommand(createCmd)
	TrainingCmd.AddCommand(deleteCmd)
	TrainingCmd.AddCommand(listCmd)
	TrainingCmd.AddCommand(describeCmd)
	TrainingCmd.AddCommand(logsCmd)
}