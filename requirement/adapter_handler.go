package requirement

import (
	"fmt"
	"strings"

	"github.com/danielcarney96/maintenance/config"
)

type Adapter struct {
	InstallCommand       string
	SwitchVersionCommand string
}

func MakePhpAdapter(req config.Requirement) Adapter {
	return Adapter{
		InstallCommand:       fmt.Sprintf("apt-get install php%s php%s-fpm php%s-cli -y", req.Version, req.Version, req.Version),
		SwitchVersionCommand: fmt.Sprintf("update-alternatives --config php -%s", req.Version),
	}
}

func MakeNodeAdapter(req config.Requirement) Adapter {
	return Adapter{
		InstallCommand:       "",
		SwitchVersionCommand: "",
	}
}

func ToCommandArray(command string) []string {
	return strings.Split(command, " ")
}
