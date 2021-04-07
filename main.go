package main

import (
	// "fmt"
	// "os"
	"xo/cmd"
	// "xo/massdriver"
)

func main() {
	// deploy, err := massdriver.GetDeployment("foobar")
	// // massdriver.UploadArtifacts(map[string]{})
	// if err != nil {
	// 	fmt.Printf("oh no: %v", err)
	// 	os.Exit(1)
	// }
	//
	// fmt.Printf("I have a nice new deployment: %s\n", deploy)

	cmd.Execute()
}
