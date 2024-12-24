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

	cmd := exec.Command("/proc/self/exe", "runContainer", name)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

// runContainer - Simulate container execution
func RunContainer(name string) {
	fmt.Printf("Running container process: %s\n", name)
	if err := syscall.Sethostname([]byte(name)); err != nil {
		fmt.Printf("Error setting hostname: %s\n", err)
		return
	}
	rootfsPath := fmt.Sprintf("/var/lib/containers/%s/rootfs", name)
	if err := syscall.Chroot(rootfsPath); err != nil {
		fmt.Printf("Error changing root filesystem: %s\n", err)
		return
	}
	if err := syscall.Chdir("/"); err != nil {
		fmt.Printf("Error changing directory: %s\n", err)
		return
	}
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Printf("Error mounting proc: %s\n", err)
		return
	}

	// Simulate container process
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running shell: %s\n", err)
	}
}

// ExecuteCommand - Execute basic commands inside container
func ExecuteCommand(containerName, command string) {
	cmd := exec.Command("/proc/self/exe", "runContainer", containerName, command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS,
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Executing in %s: %s\n", containerName, command)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %s\n", err)
	}
}
