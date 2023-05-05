package docker

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func RunCommandAndOutput(ctx context.Context, containerID string, commands []string) ExecResult {
	response, err := ExecuteCommandInContainer(ctx, containerID, commands)

	if err != nil {
		log.Fatal(err)
	}

	result, err := InspectCommandExecResponse(ctx, response.ID)

	if err != nil {
		log.Fatal(err)
	}

	return result
}

func ListDirectoriesInContainer(ctx context.Context, containerID string, path string) []string {
	commands := []string{"find", path, "-mindepth", "1", "-maxdepth", "1", "-type", "d"}

	result := RunCommandAndOutput(ctx, containerID, commands)

	return strings.Split(strings.TrimSpace(result.StdOut), "\n")
}

func ExecuteCommandInContainer(ctx context.Context, containerID string, commands []string) (types.IDResponse, error) {
	docker, err := client.NewClientWithOpts()
	if err != nil {
		return types.IDResponse{}, err
	}

	return docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          commands,
	})
}

func ContainerByName(name string) (types.Container, error) {
	containers := Containers()

	for _, container := range containers {
		for _, containerName := range container.Names {
			// Docker container names come in the format of '/name', so we need to prepend a '/'
			// see https://docs.docker.com/engine/api/v1.41/#tag/Container/operation/ContainerList for an example
			if fmt.Sprintf("/%s", name) == containerName {
				return container, nil
			}
		}
	}

	return types.Container{}, fmt.Errorf("container not found for name: %s", name)
}

func Containers() []types.Container {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	return containers
}

func InspectCommandExecResponse(ctx context.Context, id string) (ExecResult, error) {
	var execResult ExecResult
	docker, err := client.NewClientWithOpts()
	if err != nil {
		return execResult, err
	}

	resp, err := docker.ContainerExecAttach(ctx, id, types.ExecStartCheck{})
	if err != nil {
		return execResult, err
	}
	defer resp.Close()

	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return execResult, err
		}
		break

	case <-ctx.Done():
		return execResult, ctx.Err()
	}

	stdout, err := ioutil.ReadAll(&outBuf)
	if err != nil {
		return execResult, err
	}
	stderr, err := ioutil.ReadAll(&errBuf)
	if err != nil {
		return execResult, err
	}

	res, err := docker.ContainerExecInspect(ctx, id)
	if err != nil {
		return execResult, err
	}

	execResult.ExitCode = res.ExitCode
	execResult.StdOut = string(stdout)
	execResult.StdErr = string(stderr)
	return execResult, nil
}
