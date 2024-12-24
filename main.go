package main

import (
	"container-poc/container"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "create":
		handleCreate()
	case "status":
		handleStatus()
	case "exec":
		handleExec()
	case "runContainer":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ./mycontainer runContainer [container_name]")
			os.Exit(1)
		}
		containerName := os.Args[2]
		container.RunContainer(containerName)
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run main.go [create|status|exec] [container_name] [command]")
}

func handleCreate() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go create [container_name]")
		return
	}
	containerName := os.Args[2]
	container.CreateContainer(containerName)
	if err := container.ApplyCgroups(); err != nil {
		fmt.Printf("Error applying cgroups: %s\n", err)
		os.Exit(1)
	}
}

func handleStatus() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go status [container_name]")
		return
	}
	containerName := os.Args[2]
	container.ContainerStatus(containerName)
}

func handleExec() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go exec [container_name] [command]")
		return
	}
	containerName := os.Args[2]
	command := os.Args[3]
	container.ExecuteCommand(containerName, command)
}
