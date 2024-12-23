package main

import (
	"container-poc/container"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [create|status|exec] [container_name] [command]")
		return
	}

	switch os.Args[1] {
	case "create":
		containerName := os.Args[2]
		container.CreateContainer(containerName)
		container.ApplyCgroups()
	case "status":
		containerName := os.Args[2]
		container.ContainerStatus(containerName)
	case "exec":
		containerName := os.Args[2]
		command := os.Args[3]
		container.ExecuteCommand(containerName, command)
	default:
		fmt.Println("Unknown command")
	}
}
