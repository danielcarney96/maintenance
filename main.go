package main

import (
	"context"
	"fmt"
	"log"

	"github.com/danielcarney96/maintenance/config"
	"github.com/danielcarney96/maintenance/docker"
	"github.com/danielcarney96/maintenance/requirement"
)

func main() {
	contents := config.ReadRepositoriesFromFile("repositories.yml")

	ctx := context.Background()

	container, err := docker.ContainerByName("maintenance")

	if err != nil {
		log.Fatal(err)
	}

	var commands []string

	for _, data := range contents {
		for _, req := range data.Requirements {
			if req.Key == "php" {
				commands = append(commands, requirement.PhpAdapter(req)...)

				response, err := docker.ExecuteCommandInContainer(ctx, container.ID, commands)

				if err != nil {
					log.Fatal(err)
				}

				result, err := docker.InspectCommandExecResponse(ctx, response.ID)

				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf(result.StdOut)

				commands = nil
			}
		}
	}
}
