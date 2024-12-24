# Containerization Autopsy - Open Source Project

Welcome to the Containerization Autopsy project! 🕵️‍♂️

This open-source project is the culmination of an autopsy I performed on containerization technologies. It’s a hands-on containerization toolset that simulates creating, managing, and running containers with isolated namespaces, leveraging basic system commands, and integrating with cgroups for resource management.

I’ve made it open-source so everyone can learn, explore, and contribute. Feel free to break it, modify it, or simply enjoy the chaos it might bring. Welcome aboard!

Before getting started, ensure you have the following tools installed:

- **Go** (v1.16 or later) - [Install Go](https://golang.org/dl/)
- **Docker** (for containerized environments) - [Install Docker](https://www.docker.com/get-started)
- **wget** (for downloading necessary binaries) - Install via your package manager if not already installed.

Additionally, ensure you have appropriate permissions to create directories and execute system commands.

## Project Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/alchemist123/container-poc.git
   cd container-poc
2. Install necessary dependencies:
    ```bash
    go mod tidy
Commands
This project implements a simple container management tool with the following commands:

create
Creates a new container with a specified name.

Usage:
go run main.go create [container_name]

status
Checks if a container is running.

Usage:
go run main.go status [container_name]

exec
Executes a command inside a running container.

Usage:
go run main.go exec [container_name] [command]

runContainer
Simulates running a container by using namespaces and cgroups to isolate processes. This is the core function that "runs" the container after it's created.

Usage:
go run main.go runContainer [container_name]

Features
Namespace Isolation: Creates isolated environments for containers with PID, UTS, and network namespaces.
Cgroup Support: Uses Linux cgroups to limit CPU and memory usage for containers.
Root Filesystem Setup: Automatically sets up a minimal root filesystem using Busybox.
Container Status: Checks the status of a container and prints basic information.
Command Execution: Execute commands within the container's isolated environment

#Features
Namespace Isolation: Creates isolated environments for containers with PID, UTS, and network namespaces.
Cgroup Support: Uses Linux cgroups to limit CPU and memory usage for containers.
Root Filesystem Setup: Automatically sets up a minimal root filesystem using Busybox.
Container Status: Checks the status of a container and prints basic information.
Command Execution: Execute commands within the container's isolated environment.
License
This project is licensed under the MIT License - see the LICENSE file for details.

Contribution Guidelines
We welcome contributions! To contribute to this project, please follow these steps:

Fork the repository.
Create a new branch for your feature or bug fix.
Make your changes and commit them.
Submit a pull request with a detailed description of your changes.
Please ensure that your code follows the existing coding style and that any new functionality is well-tested.

Open Source Rules
Respect Others: Be kind and respectful in discussions and contributions.
Contributions: Contributions are always welcome, but please make sure to issue a PR for review before merging.
Fork It!: Feel free to fork the repository, create new features, and share your improvements!