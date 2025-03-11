package actuator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Pigeon-Developer/cgroups/cgroup2"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-version"
	"golang.org/x/sys/unix"
)

const (
	_DefaultMountPoint = "/sys/fs/cgroup"
	Prefix             = "pojj-"
)

var (
	IsSystemd = false
)

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

	baseVesion, _ := version.NewVersion("5.19")

	strs := strings.Split(string(uname.Release[:]), "-")
	curentVesion, err := version.NewVersion(strs[0])

	if err != nil {
		fmt.Println("parse version error ", string(uname.Release[:]))
		panic(err)
	}

	if baseVesion.GreaterThan(curentVesion) {
		panic("require Linux 5.19+ kernel")
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("create docker client error")
		panic(err)
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		fmt.Println("fetch docker info error")
		panic(err)
	}

	if info.CgroupDriver == "systemd" {
		IsSystemd = true
	}
}

func NewCgroupWrap(_path string) (*CgroupWrap, error) {
	if IsSystemd {
		_path = strings.ReplaceAll(_path, "-", "")
		path := Prefix + _path + ".slice"
		c, err := cgroup2.NewSystemd("/", path, -1, &cgroup2.Resources{})
		if err != nil {
			return nil, err
		}

		instance := CgroupWrap{
			Path:   path,
			Cgroup: c,
		}

		return &instance, nil
	}
	path := "/" + Prefix + _path

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

func (c CgroupWrap) Delete() error {
	if IsSystemd {
		return c.Cgroup.DeleteSystemd()

	}
	return c.Cgroup.Delete()
}
