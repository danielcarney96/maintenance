package requirement

import (
	"fmt"

	"github.com/danielcarney96/maintenance/config"
)

type Adapter struct {
	InstallCommands      []string
	SwitchVersionCommand string
}

func MakePhpAdapter(req config.Requirement) Adapter {
	return Adapter{
		InstallCommands: []string{
			"apt-get",
			"install",
			fmt.Sprintf("php%s", req.Version),
			fmt.Sprintf("php%s-cli", req.Version),
			"-y",
		},
		SwitchVersionCommand: fmt.Sprintf("update-alternatives --config php -%s", req.Version),
	}
}

func MakeNodeAdapter(req config.Requirement) Adapter {
	return Adapter{
		InstallCommands: []string{
			"bash",
			"-c",
			fmt.Sprintf("PS1=x; . ~/.bashrc; nvm install v%s", req.Version),
		},
		SwitchVersionCommand: fmt.Sprintf(`bash -c "PS1=x; . ~/.bashrc; nvm use v%s"`, req.Version),
	}
}
