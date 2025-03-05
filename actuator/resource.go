package actuator

import (
	"fmt"
	"os/exec"

	"github.com/containerd/cgroups/v3/cgroup2"

	"golang.org/x/sys/unix"
)

const _DefaultMountPoint = "/sys/fs/cgroup"

var Prefix = "/pojj/"

type CgroupWrap struct {
	Path   string
	Cgroup *cgroup2.Manager
}

func init() {
	out, _ := exec.Command("bash", "-c", "stat -fc %T /sys/fs/cgroup/").Output()

	if string(out) != "cgroup2fs\n" {
		panic("cgroup2 is not mounted")
	}
	var uname unix.Utsname
	unix.Uname(&uname)

	fmt.Println(string(uname.Release[:]))
	panic("cgroup2 is not mounted")
}

func NewCgroupWrap(_path string) (*CgroupWrap, error) {
	path := Prefix + _path

	c, err := cgroup2.NewManager(_DefaultMountPoint, path, &cgroup2.Resources{})
	if err != nil {
		return nil, err
	}

	instance := CgroupWrap{
		Path:   path,
		Cgroup: c,
	}

	return &instance, nil
}
