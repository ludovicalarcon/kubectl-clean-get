package cmd

import (
	"bufio"
	"bytes"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Input string
var Output string

var rootCmd = &cobra.Command{
	Use:   "kubectl-clean-get",
	Short: "kubectl-clean-get removes fields that are instance specific",
	Long: `kubectl-clean-get removes fields that are instance specific.

Read from stdin and output on stdout:
> kubectl get TYPE NAME | kubectl-clean-get -f -

Read from file and outpout on stdout:
> kubectl-clean-get -f FILE

Output in a file:
> kubectl-clean-get -f FILE -o OUTPUT
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestInput, err := readInput(Input)
		if err != nil {
			return err
		}

		var data map[string]interface{}
		err = yaml.Unmarshal(manifestInput, &data)
		if err != nil {
			return err
		}

		_, err  = cleanManifestInput(data)
		if err != nil {
			return err
		}

		return nil
	},
}

func readInput(input string) ([]byte, error) {
	if input == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		input := ""
		for scanner.Scan() {
			tmp := scanner.Text()
			if tmp != "EOF" {
				input += (tmp + "\n")
			} else {
				break
			}
		}
		return []byte(input), scanner.Err()
	} else {
		return os.ReadFile(input)
	}
}

func cleanManifestInput(data map[string]interface{}) ([]byte, error) {
	if data["uid"] != nil {
		delete(data, "uid")
	}

	yaml, err := EncodeYaml(data)

	return yaml, err
}

func EncodeYaml(data map[string]interface{}) ([]byte, error) {
	var raw bytes.Buffer
	encoder := yaml.NewEncoder(&raw)
	encoder.SetIndent(2)
	err := encoder.Encode(&data)

	return raw.Bytes(), err
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&Input, "input-file", "f", "", "Input yaml file, use '-' for stdin (required)")
	rootCmd.Flags().StringVarP(&Output, "output-file", "o", "", "Output yaml file")
	rootCmd.MarkFlagRequired("input-file")
}
