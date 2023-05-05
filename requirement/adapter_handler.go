package requirement

import (
	"fmt"

	"github.com/danielcarney96/maintenance/config"
)

type Adapter struct {
	InstallCommand       string
	SwitchVersionCommand string
}

func MakePhpAdapter(req config.Requirement) Adapter {
	return Adapter{
		InstallCommand:       fmt.Sprintf("apt-get install php%s php%s-cli -y", req.Version, req.Version),
		SwitchVersionCommand: fmt.Sprintf("update-alternatives --config php-%s", req.Version),
	}
}

func MakeNodeAdapter(req config.Requirement) Adapter {
	return Adapter{
		InstallCommand:       fmt.Sprintf("bash -c PS1=x; . ~/.bashrc; nvm install v%s", req.Version),
		SwitchVersionCommand: fmt.Sprintf("bash -c PS1=x; . ~/.bashrc; nvm use v%s", req.Version),
	}
}
