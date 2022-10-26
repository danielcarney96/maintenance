package requirement

import (
	"fmt"
	"strings"

	"github.com/danielcarney96/maintenance/config"
)

func PhpAdapter(req config.Requirement) []string {
	installCommand := fmt.Sprintf("apt-get install php%s php%s-fpm php%s-cli -y", req.Version, req.Version, req.Version)

	return ToCommandArray(installCommand)
}

func ToCommandArray(command string) []string {
	return strings.Split(command, " ")
}
