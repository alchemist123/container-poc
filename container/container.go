package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// CreateContainer - Creates a container with isolated namespaces
func CreateContainer(name string) {
	fmt.Printf("Creating container: %s\n", name)

	// Fork a new process to simulate container
	cmd := exec.Command("/proc/self/exe", "child")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

// Execute basic commands inside container
func ExecuteCommand(containerName, command string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWPID}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Executing in %s: %s\n", containerName, command)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %s\n", err)
	}
}
