package repositories

import (
	"context"
	"fmt"

	"github.com/danielcarney96/maintenance/docker"
)

func BranchAndCommit(ctx context.Context, containerID string, repoDir string, branchName string, commitMessage string) {
	// Create branch
	commands := []string{"git", "-C", repoDir, "checkout", "-b", branchName}
	result := docker.RunCommandAndOutput(ctx, containerID, commands)
	fmt.Printf(result.StdOut)

	// Stage changes
	commands = []string{"git", "-C", repoDir, "add", "."}
	result = docker.RunCommandAndOutput(ctx, containerID, commands)
	fmt.Printf(result.StdOut)

	// Commit changes
	commands = []string{"git", "-C", repoDir, "commit", "-m", commitMessage}
	result = docker.RunCommandAndOutput(ctx, containerID, commands)
	fmt.Printf(result.StdOut)
}

func PushAndPR(ctx context.Context, containerID string, repoDir string, branchName string, prTitle string, prBody string) {
	// Push branch
	commands := []string{"git", "-C", repoDir, "push", "origin", branchName}
	result := docker.RunCommandAndOutput(ctx, containerID, commands)
	fmt.Printf(result.StdOut)

	// Create PR
	commands = []string{"sh", "-c", fmt.Sprintf(`cd %q && gh pr create --title "%s" --body "%s"`, repoDir, prTitle, prBody)}
	result = docker.RunCommandAndOutput(ctx, containerID, commands)
	fmt.Printf(result.StdOut)
}
