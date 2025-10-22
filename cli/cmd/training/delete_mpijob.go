package training

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	deleteForce bool
	deleteWait  bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a training job",
	Long:  `Delete a training job from OpenShift AI.`,
}

// deleteMPIJobCmd represents the delete mpijob command
var deleteMPIJobCmd = &cobra.Command{
	Use:   "mpijob NAME [flags]",
	Short: "Delete an MPIJob",
	Long: `Delete an MPIJob from OpenShift AI.

Examples:
  # Delete an MPIJob
  odh training delete mpijob my-job

  # Force delete an MPIJob
  odh training delete mpijob my-job --force

  # Delete an MPIJob and wait for cleanup
  odh training delete mpijob my-job --wait
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		namespace, _ := cmd.Flags().GetString("namespace")

		// Load kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			klog.Fatalf("Error building kubeconfig: %s", err.Error())
		}

		jobName := args[0]

		// For now, just print that we would delete the MPIJob
		fmt.Printf("Deleting MPIJob %s in namespace %s\n", jobName, namespace)

		if deleteForce {
			fmt.Println("Force delete enabled")
		}

		if deleteWait {
			fmt.Println("Waiting for job cleanup to complete")
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteMPIJobCmd)

	// Add flags
	deleteMPIJobCmd.Flags().BoolVar(&deleteForce, "force", false, "Force delete the job")
	deleteMPIJobCmd.Flags().BoolVar(&deleteWait, "wait", false, "Wait for job cleanup to complete")
}