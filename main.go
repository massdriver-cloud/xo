package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"xo/massdriver"
)

func main() {
	// cmd.Execute()

	// The URL below sucks. We need to fix something in elixir so its not so redundant...
	// http://localhost:4000/rpc/deployment/twirp/mdtwirp.Deployments/Get
	client := massdriver.NewDeploymentsProtobufClient("http://localhost:4000/rpc/deployment", &http.Client{})

	dep, err := client.Get(context.Background(), &massdriver.GetDeploymentRequest{Id: "1"})

	if err != nil {
		fmt.Printf("oh no: %v", err)
		os.Exit(1)
	}

	fmt.Printf("I have a nice new deployment: %+v", dep)
}
