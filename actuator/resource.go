package actuator

import (
	"github.com/containerd/cgroups/v3/cgroup2"
)

const DefaultMountPoint = "/sys/fs/cgroup"

type CgroupWrap struct {
	Path   string
	Cgroup *cgroup2.Manager
}

func NewCgroupWrap(_path string) (*CgroupWrap, error) {
	path := "/pojj/" + _path

	c, err := cgroup2.NewManager(DefaultMountPoint, path, &cgroup2.Resources{})
	if err != nil {
		return nil, err
	}

	instance := CgroupWrap{
		Path:   path,
		Cgroup: c,
	}

	return &instance, nil
}
