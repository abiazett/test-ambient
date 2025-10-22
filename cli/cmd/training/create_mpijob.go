package training

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

var (
	createFileName string
	createFromStdin bool
	createDryRun bool
	createWait bool
	createWorkers int
	createGPU int
	createCPU string
	createMemory string
	createImage string
	createCommand []string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a training job",
	Long: `Create a training job in OpenShift AI.`,
}

// createMPIJobCmd represents the create mpijob command
var createMPIJobCmd = &cobra.Command{
	Use:   "mpijob [NAME] [flags]",
	Short: "Create an MPIJob",
	Long: `Create an MPIJob in OpenShift AI.

Examples:
  # Create an MPIJob from a file
  odh training create mpijob --from-file=mpijob.yaml

  # Create an MPIJob from stdin
  cat mpijob.yaml | odh training create mpijob --from-stdin

  # Create an MPIJob with inline parameters
  odh training create mpijob my-job --workers=4 --gpu=2 --cpu=4 --memory=16Gi --image=myregistry.com/myimage:latest --command="python" --command="/train.py"

  # Create an MPIJob with dry run
  odh training create mpijob --from-file=mpijob.yaml --dry-run

  # Create an MPIJob and wait for completion
  odh training create mpijob --from-file=mpijob.yaml --wait
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		namespace, _ := cmd.Flags().GetString("namespace")

		// Load kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			klog.Fatalf("Error building kubeconfig: %s", err.Error())
		}

		// For now, just print that we would create an MPIJob
		fmt.Printf("Creating MPIJob in namespace %s\n", namespace)

		// Handle file input
		if createFileName != "" {
			content, err := ioutil.ReadFile(createFileName)
			if err != nil {
				klog.Fatalf("Error reading file: %s", err.Error())
			}
			fmt.Printf("Would create MPIJob from file %s\n", createFileName)
			return
		}

		// Handle stdin input
		if createFromStdin {
			content, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				klog.Fatalf("Error reading from stdin: %s", err.Error())
			}
			fmt.Println("Would create MPIJob from stdin")
			return
		}

		// Handle inline parameters
		if len(args) > 0 {
			jobName := args[0]
			fmt.Printf("Would create MPIJob %s with:\n", jobName)
			fmt.Printf("  Workers: %d\n", createWorkers)
			fmt.Printf("  GPU per worker: %d\n", createGPU)
			fmt.Printf("  CPU per worker: %s\n", createCPU)
			fmt.Printf("  Memory per worker: %s\n", createMemory)
			fmt.Printf("  Image: %s\n", createImage)
			fmt.Printf("  Command: %v\n", createCommand)
			return
		}

		fmt.Println("Error: either --from-file, --from-stdin, or inline parameters must be specified")
	},
}

func init() {
	createCmd.AddCommand(createMPIJobCmd)

	// File and stdin flags
	createMPIJobCmd.Flags().StringVar(&createFileName, "from-file", "", "Create from file")
	createMPIJobCmd.Flags().BoolVar(&createFromStdin, "from-stdin", false, "Create from stdin")

	// Inline parameter flags
	createMPIJobCmd.Flags().IntVar(&createWorkers, "workers", 1, "Number of worker replicas")
	createMPIJobCmd.Flags().IntVar(&createGPU, "gpu", 0, "Number of GPUs per worker")
	createMPIJobCmd.Flags().StringVar(&createCPU, "cpu", "1", "CPU cores per worker")
	createMPIJobCmd.Flags().StringVar(&createMemory, "memory", "1Gi", "Memory per worker")
	createMPIJobCmd.Flags().StringVar(&createImage, "image", "", "Container image to use")
	createMPIJobCmd.Flags().StringArrayVar(&createCommand, "command", []string{}, "Command to run (specify multiple times for array)")

	// Other flags
	createMPIJobCmd.Flags().BoolVar(&createDryRun, "dry-run", false, "Only print the object that would be created")
	createMPIJobCmd.Flags().BoolVar(&createWait, "wait", false, "Wait for job completion")
}