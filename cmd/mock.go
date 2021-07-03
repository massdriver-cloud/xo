package cmd

import (
	"fmt"
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
}

func RunMock(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetString("port")

	fmt.Println("Starting a mock Massdriver server on port localhost:" + port + "...")

	return massdriver.RunMockServer(port)
}
