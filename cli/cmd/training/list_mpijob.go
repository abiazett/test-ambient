package training

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	listOutputFormat string
	listLabelSelector string
	listAllNamespaces bool
	listStatus string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List training jobs",
	Long:  `List training jobs in OpenShift AI.`,
}

// listMPIJobCmd represents the list mpijob command
var listMPIJobCmd = &cobra.Command{
	Use:   "mpijob [flags]",
	Short: "List MPIJobs",
	Long: `List MPIJobs in OpenShift AI.

Examples:
  # List all MPIJobs
  odh training list mpijob

  # List all MPIJobs with status=Running
  odh training list mpijob --status=Running

  # List all MPIJobs with custom label selector
  odh training list mpijob -l app=my-app

  # List MPIJobs across all namespaces
  odh training list mpijob --all-namespaces

  # List MPIJobs with custom output format
  odh training list mpijob -o wide
  odh training list mpijob -o json
  odh training list mpijob -o yaml
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		namespace, _ := cmd.Flags().GetString("namespace")

		// Use all namespaces if specified
		if listAllNamespaces {
			namespace = ""
			fmt.Println("Listing MPIJobs across all namespaces")
		} else {
			fmt.Printf("Listing MPIJobs in namespace %s\n", namespace)
		}

		// Load kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			klog.Fatalf("Error building kubeconfig: %s", err.Error())
		}

		// For now, just print some sample output
		printSampleMPIJobList(namespace, listOutputFormat, listLabelSelector, listStatus)
	},
}

// Helper function to print a sample job list
func printSampleMPIJobList(namespace, outputFormat, labelSelector, status string) {
	// Create a table with headers based on output format
	if outputFormat == "wide" {
		fmt.Printf("%-20s %-12s %-8s %-8s %-12s %-8s %-8s %-20s\n",
			"NAME", "STATUS", "WORKERS", "GPU", "DURATION", "AGE", "NAMESPACE", "IMAGE")
	} else if outputFormat == "name" {
		fmt.Println("NAME")
	} else {
		// Default table format
		fmt.Printf("%-30s %-12s %-8s %-12s %-8s\n",
			"NAME", "STATUS", "WORKERS", "DURATION", "AGE")
	}

	// Sample data
	fmt.Println("No MPIJobs found.")
}

func init() {
	listCmd.AddCommand(listMPIJobCmd)

	// Add flags
	listMPIJobCmd.Flags().StringVarP(&listOutputFormat, "output", "o", "", "Output format (name|wide|json|yaml)")
	listMPIJobCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector")
	listMPIJobCmd.Flags().BoolVarP(&listAllNamespaces, "all-namespaces", "A", false, "List jobs across all namespaces")
	listMPIJobCmd.Flags().StringVar(&listStatus, "status", "", "Filter by status (Pending|Running|Succeeded|Failed)")
}