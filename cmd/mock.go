package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"xo/src/massdriver"

	"github.com/spf13/cobra"
)

var mockLong = `
	Starts a mock Massdriver server for testing twirp functionality.

	Twirp is the protocol we're using for communication between xo and Massdriver.
	To local testing you can use this mock server to test xo against.
	`
var mockExamples = `
	# Start a mock Massdriver server
	xo mock
	`

var mockCmd = &cobra.Command{
	Use:                   "mock",
	Short:                 "Start a mock Massdriver server for testing",
	Long:                  mockLong,
	Example:               mockExamples,
	RunE:                  RunMock,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(mockCmd)
	mockCmd.Flags().StringP("port", "p", "8080", "Port to run the server on (default 8080)")
	mockCmd.Flags().StringP("params-file", "i", "", "JSON file representing data to be returned as params (default is empty JSON object)")
	mockCmd.Flags().StringP("connections-file", "c", "", "JSON file representing data to be returned as connections (default is empty JSON object)")
}

func RunMock(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetString("port")
	paramsPath, _ := cmd.Flags().GetString("params-file")
	connectionsPath, _ := cmd.Flags().GetString("connections-file")

	paramsBytes := []byte("{}")
	connectionsBytes := []byte("{}")

	if paramsPath != "" {
		bytes, err := readFile(paramsPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		paramsBytes = bytes
	}
	if connectionsPath != "" {
		bytes, err := readFile(connectionsPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		connectionsBytes = bytes
	}

	fmt.Println("Starting a mock Massdriver server on port localhost:" + port + "...")

	return massdriver.RunMockServer(port, string(paramsBytes), string(connectionsBytes))
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	return bytes, nil
}
