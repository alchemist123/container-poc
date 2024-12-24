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
			Shares: uint64Ptr(512), // 50% CPU
		},
		Memory: &specs.LinuxMemory{
			Limit: int64Ptr(128 * 1024 * 1024), // 128MB Memory
		},
	})
	if err != nil {
		return err
	}

	return control.Add(cgroups.Process{Pid: os.Getpid()})
}

func uint64Ptr(i uint64) *uint64 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}
