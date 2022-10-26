package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"gopkg.in/yaml.v3"
)

type Repo struct {
	Url string
	Php string
}

func main() {
	data, err := ioutil.ReadFile("repositories.yml")

	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]Repo)

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	container := docker()

	var commands []string

	commands = append(commands, "apt-get")
	commands = append(commands, "install")
	for _, data := range m {
		commands = append(commands, fmt.Sprintf("php%s", data.Php))
		commands = append(commands, fmt.Sprintf("php%s-fpm", data.Php))
		commands = append(commands, fmt.Sprintf("php%s-cli", data.Php))
	}
	commands = append(commands, "-y")

	response, err := exec(ctx, container.ID, commands)

	if err != nil {
		log.Fatal(err)
	}

	result, err := inspectExecResp(ctx, response.ID)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(result.StdOut)
}

func docker() types.Container {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	return containers[0]
}

func exec(ctx context.Context, containerID string, command []string) (types.IDResponse, error) {
	docker, err := client.NewClientWithOpts()
	if err != nil {
		return types.IDResponse{}, err
	}

	config := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          command,
	}

	return docker.ContainerExecCreate(ctx, containerID, config)
}

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func inspectExecResp(ctx context.Context, id string) (ExecResult, error) {
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

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
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
