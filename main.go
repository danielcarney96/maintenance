package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/danielcarney96/maintenance/config"
	"github.com/danielcarney96/maintenance/docker"
	"github.com/danielcarney96/maintenance/repositories"
	"github.com/danielcarney96/maintenance/requirement"
)

func main() {
	contents := config.ReadRepositoriesFromFile("repositories.yml")

	ctx := context.Background()

	container, err := docker.ContainerByName("maintenance")

	if err != nil {
		log.Fatal(err)
	}

	for _, data := range contents {
		commands := []string{"git", "-C", "repositories", "clone", data.Url}
		result := docker.RunCommandAndOutput(ctx, container.ID, commands)
		fmt.Printf(result.StdOut)

		for _, req := range data.Requirements {
			var adapter requirement.Adapter

			switch req.Key {
			case "php":
				adapter = requirement.MakePhpAdapter(req)
			case "node":
				adapter = requirement.MakeNodeAdapter(req)
			default:
				continue
			}

			result = docker.RunCommandAndOutput(ctx, container.ID, adapter.InstallCommands)
			fmt.Printf(result.StdOut)
		}
	}

	directories := docker.ListDirectoriesInContainer(ctx, container.ID, "/repositories")

	for _, dir := range directories {
		commands := []string{"sh", "-c", fmt.Sprintf("cd %q && composer update", dir)}
		result := docker.RunCommandAndOutput(ctx, container.ID, commands)
		fmt.Printf(result.StdOut)

		branchName := fmt.Sprintf("maintenance/%s-%d", strings.ToLower(time.Now().Month().String()), time.Now().Year())

		repositories.BranchAndCommit(ctx, container.ID, dir, branchName, "update composer libraries")
		repositories.PushAndPR(ctx, container.ID, dir, branchName, fmt.Sprintf("Maintenance/%s %d", time.Now().Month().String(), time.Now().Year()), "test")
	}
}
