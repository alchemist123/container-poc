package container

import (
	"fmt"
	"os"
)

// ContainerStatus - Checks if container exists
func ContainerStatus(containerName string) {
	if _, err := os.Stat(fmt.Sprintf("/mycontainer/%s", containerName)); err == nil {
		fmt.Printf("Container %s is running\n", containerName)
	} else {
		fmt.Printf("Container %s not found\n", containerName)
	}
}
