/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"xo/src/tfdef"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var provisionerCmd = &cobra.Command{
	Use:   "provisioner",
	Short: "Manage provisioners",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("provisioner called")
	// },
}

var provisionerCompileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile JSON Schema to provisioner variable definition file.",
	Long:  ``,
	RunE:  runProvisionerCompile,
}

func init() {
	rootCmd.AddCommand(provisionerCmd)
	provisionerCmd.AddCommand(provisionerCompileCmd)
	provisionerCompileCmd.Flags().StringP("schema", "s", "./schema.json", "Path to JSON Schema")
	provisionerCompileCmd.Flags().StringP("output", "o", "./variables.tf.json", "Output path. Use - for STDOUT")
}

func runProvisionerCompile(cmd *cobra.Command, args []string) error {
	var compiled string
	var err error

	provisioner := args[0]
	schema, _ := cmd.Flags().GetString("schema")
	outputPath, _ := cmd.Flags().GetString("output")

	log.Debug().
		Str("provisioner", provisioner).
		Str("schemaPath", schema).Msg("Compiling schema.")

	switch provisioner {
	case "terraform":
		compiled, err = tfdef.Compile(schema)
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("Unsupported argument %s the single argument 'terraform' is supported", provisioner)
		log.Error().Err(err).Msg("Compilation failed.")
		return err
	}

	return writeVariableFile(compiled, outputPath)
}

func writeVariableFile(compiled string, outPath string) error {
	switch outPath {
	case "-":
		fmt.Println(compiled)
		return nil
	default:
		data := []byte(compiled)
		err := ioutil.WriteFile(outPath, data, 0644)
		return err
	}
}
