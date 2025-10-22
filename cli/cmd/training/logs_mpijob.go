package training

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	logsFollow    bool
	logsTail      int
	logsWorker    string
	logsContainer string
	logsAggregate bool
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get logs from a training job",
	Long:  `Get logs from a training job in OpenShift AI.`,
}

// logsMPIJobCmd represents the logs mpijob command
var logsMPIJobCmd = &cobra.Command{
	Use:   "mpijob NAME [flags]",
	Short: "Get logs from an MPIJob",
	Long: `Get logs from an MPIJob in OpenShift AI.

Examples:
  # Get logs from an MPIJob launcher
  odh training logs mpijob my-job

  # Get logs from a specific worker
  odh training logs mpijob my-job --worker=0

  # Get logs from all workers
  odh training logs mpijob my-job --aggregate

  # Follow logs in real-time
  odh training logs mpijob my-job -f

  # Show only the last N lines
  odh training logs mpijob my-job --tail=100
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

		// For now, just print that we would get logs from the MPIJob
		fmt.Printf("Getting logs for MPIJob %s in namespace %s\n", jobName, namespace)

		if logsWorker != "" {
			fmt.Printf("Worker: %s\n", logsWorker)
		} else {
			fmt.Println("Worker: launcher")
		}

		if logsContainer != "" {
			fmt.Printf("Container: %s\n", logsContainer)
		}

		if logsAggregate {
			fmt.Println("Aggregating logs from all workers")
		}

		if logsFollow {
			fmt.Println("Following logs in real-time")
		}

		if logsTail > 0 {
			fmt.Printf("Showing last %d lines\n", logsTail)
		}

		// Print sample log output
		printSampleMPIJobLogs(jobName, logsWorker)
	},
}

// Helper function to print sample logs
func printSampleMPIJobLogs(name, worker string) {
	fmt.Println("Sample log output:")
	fmt.Println("2025-10-21T15:04:05.000Z INFO Starting MPI job")
	fmt.Println("2025-10-21T15:04:06.000Z INFO Initializing worker nodes")
	fmt.Println("2025-10-21T15:04:07.000Z INFO MPI rank assigned")
	fmt.Println("2025-10-21T15:04:08.000Z INFO Beginning training")
	fmt.Println("2025-10-21T15:04:09.000Z INFO Epoch 1/10, Loss: 0.452")
	fmt.Println("2025-10-21T15:04:10.000Z INFO Epoch 2/10, Loss: 0.387")
}

func init() {
	logsCmd.AddCommand(logsMPIJobCmd)

	// Add flags
	logsMPIJobCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow logs")
	logsMPIJobCmd.Flags().IntVar(&logsTail, "tail", 0, "Lines of recent logs to display")
	logsMPIJobCmd.Flags().StringVar(&logsWorker, "worker", "", "Worker index or 'launcher'")
	logsMPIJobCmd.Flags().StringVar(&logsContainer, "container", "", "Container name")
	logsMPIJobCmd.Flags().BoolVar(&logsAggregate, "aggregate", false, "Aggregate logs from all workers")
}