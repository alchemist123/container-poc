package container

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "syscall"
)

// CreateContainer - Creates a container with isolated namespaces
func CreateContainer(name string) {
    fmt.Printf("Creating container: %s\n", name)
    
    rootfsPath := fmt.Sprintf("/var/lib/containers/%s/rootfs", name)
    if _, err := os.Stat(rootfsPath); os.IsNotExist(err) {
        fmt.Printf("Root filesystem not found, creating it...\n")
        if err := os.MkdirAll(rootfsPath, 0755); err != nil {
            fmt.Printf("Error creating root filesystem: %s\n", err)
            os.Exit(1)
        }
        // Populate root filesystem with minimal structure
        if err := setupRootfs(rootfsPath); err != nil {
            fmt.Printf("Error setting up root filesystem: %s\n", err)
            os.Exit(1)
        }
    }

    // Create a container status file to indicate it's running
    statusFilePath := fmt.Sprintf("/var/lib/containers/%s/container_status", name)
    if _, err := os.Stat(statusFilePath); os.IsNotExist(err) {
        statusContent := "Container is running" // Status message
        if err := ioutil.WriteFile(statusFilePath, []byte(statusContent), 0644); err != nil {
            fmt.Printf("Error creating status file: %s\n", err)
            os.Exit(1)
        }
    }

    // Continue with container execution (run container process)
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


// RunContainer - Simulate container execution
func RunContainer(name string) {
    fmt.Printf("Running container process: %s\n", name)

    // Debug: Print current working directory
    if cwd, err := os.Getwd(); err == nil {
        fmt.Printf("Current working directory before chroot: %s\n", cwd)
    }

    rootfsPath := fmt.Sprintf("/var/lib/containers/%s/rootfs", name)
    
    // Debug: Check rootfs contents
    fmt.Printf("Checking contents of rootfs at: %s\n", rootfsPath)
    if files, err := ioutil.ReadDir(rootfsPath); err == nil {
        fmt.Println("Rootfs contents:")
        for _, file := range files {
            fmt.Printf("- %s\n", file.Name())
        }
    } else {
        fmt.Printf("Error reading rootfs: %s\n", err)
        return
    }

    // Debug: Check /bin contents
    binPath := fmt.Sprintf("%s/bin", rootfsPath)
    fmt.Printf("\nChecking contents of /bin at: %s\n", binPath)
    if files, err := ioutil.ReadDir(binPath); err == nil {
        fmt.Println("Bin contents:")
        for _, file := range files {
            if file.Mode()&os.ModeSymlink != 0 {
                link, _ := os.Readlink(fmt.Sprintf("%s/%s", binPath, file.Name()))
                fmt.Printf("- %s -> %s\n", file.Name(), link)
            } else {
                fmt.Printf("- %s (regular file)\n", file.Name())
            }
        }
    } else {
        fmt.Printf("Error reading bin directory: %s\n", err)
        return
    }

    // Debug: Check busybox exists and is executable
    busyboxPath := fmt.Sprintf("%s/bin/busybox", rootfsPath)
    if info, err := os.Stat(busyboxPath); err == nil {
        fmt.Printf("\nBusybox permissions: %v\n", info.Mode())
    } else {
        fmt.Printf("Error checking busybox: %s\n", err)
        return
    }

    if err := syscall.Chroot(rootfsPath); err != nil {
        fmt.Printf("Error changing root filesystem: %s\n", err)
        return
    }

    if err := syscall.Chdir("/"); err != nil {
        fmt.Printf("Error changing directory: %s\n", err)
        return
    }

    // Debug: Print working directory after chroot
    if cwd, err := os.Getwd(); err == nil {
        fmt.Printf("Current working directory after chroot: %s\n", cwd)
    }

    // Mount proc
    if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
        fmt.Printf("Error mounting proc: %s\n", err)
        return
    }

    // Try to execute shell with full path and debug output
    cmd := exec.Command("/bin/sh")
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        fmt.Printf("Error running shell: %s\nError type: %T\n", err, err)
        if exitError, ok := err.(*exec.ExitError); ok {
            fmt.Printf("Exit status: %d\n", exitError.ExitCode())
        }
    }
}

// setupRootfs - Populates root filesystem with busybox
func setupRootfs(rootfsPath string) error {
    busyboxURL := "https://busybox.net/downloads/binaries/1.35.0-x86_64-linux-musl/busybox"
    busyboxPath := fmt.Sprintf("%s/bin/busybox", rootfsPath)

    // Create essential directories
    dirs := []string{"bin", "sbin", "usr/bin", "usr/sbin", "proc", "sys", "dev"}
    for _, dir := range dirs {
        dirPath := fmt.Sprintf("%s/%s", rootfsPath, dir)
        if err := os.MkdirAll(dirPath, 0755); err != nil {
            return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
        }
    }

    // Download busybox if not exists
    if _, err := os.Stat(busyboxPath); os.IsNotExist(err) {
        fmt.Println("Downloading busybox...")
        cmd := exec.Command("wget", "-O", busyboxPath, busyboxURL)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
            return fmt.Errorf("error downloading busybox: %v", err)
        }
    }

    // Setup executable permissions
    if err := os.Chmod(busyboxPath, 0755); err != nil {
        return fmt.Errorf("failed to chmod busybox: %v", err)
    }

    // Essential commands that should be available
    binaries := []string{
        "sh",
        "ls",
        "mkdir",
        "mount",
        "umount",
        "ps",
        "cat",
        "echo",
        "pwd",
    }

    // Create symlinks for all basic commands
    for _, binary := range binaries {
        linkPath := fmt.Sprintf("%s/bin/%s", rootfsPath, binary)
        if _, err := os.Stat(linkPath); err == nil {
            // Remove existing symlink if it exists
            os.Remove(linkPath)
        }
        if err := os.Symlink("busybox", linkPath); err != nil {
            return fmt.Errorf("error creating symlink for %s: %v", binary, err)
        }
    }

    fmt.Println("Root filesystem setup complete.")
    return nil
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
func Status(containerName string) {
	fmt.Printf("Checking status of container: %s\n", containerName)

	rootfsPath := fmt.Sprintf("/var/lib/containers/%s/rootfs", containerName)
	if _, err := os.Stat(rootfsPath); os.IsNotExist(err) {
		fmt.Printf("Container not found: %s\n", containerName)
		return
	}

	// Checking for process information of the container
	// Here we assume that each container is running as a separate PID namespace (or you can adjust based on your setup)
	statusFile := fmt.Sprintf("/var/lib/containers/%s/container_status", containerName)
	if _, err := os.Stat(statusFile); os.IsNotExist(err) {
		fmt.Printf("Container status file not found: %s\n", containerName)
		return
	}

	// Read container's status (process information)
	statusCmd := exec.Command("cat", statusFile)
	statusCmd.Stdout = os.Stdout
	statusCmd.Stderr = os.Stderr
	if err := statusCmd.Run(); err != nil {
		fmt.Printf("Error retrieving status for container: %s\n", err)
		return
	}

	// Debug: Additional info like running processes (you can tailor this as per your container's setup)
	fmt.Println("\nContainer processes:")
	psCmd := exec.Command("ps", "-ef") // Example to list processes
	psCmd.Stdout = os.Stdout
	psCmd.Stderr = os.Stderr
	if err := psCmd.Run(); err != nil {
		fmt.Printf("Error listing processes: %s\n", err)
		return
	}
}