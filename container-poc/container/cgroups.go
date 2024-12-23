package container

import (
	"os"

	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// ApplyCgroups - Apply memory and CPU limits
func ApplyCgroups() error {
	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath("/mycontainer"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: cgroups.NewUint64(512), // 50% CPU
		},
		Memory: &specs.LinuxMemory{
			Limit: cgroups.NewInt64(128 * 1024 * 1024), // 128MB Memory
		},
	})
	if err != nil {
		return err
	}

	process := cgroups.Process{
		Pid: os.Getpid(),
	}
	return control.Add(cgroups.Process{Pid: os.Getpid()})
}
