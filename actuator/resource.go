package actuator

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/hashicorp/go-version"
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

	baseVesion, err := version.NewVersion("5.19")

	strs := strings.Split(string(uname.Release[:]), "-")
	curentVesion, err := version.NewVersion(strs[0])

	if err != nil {
		fmt.Println("parse version error ", string(uname.Release[:]))
		panic(err)
	}

	if baseVesion.GreaterThan(curentVesion) {
		panic("require Linux 5.19+ kernel")
	}
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
