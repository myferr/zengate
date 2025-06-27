package main

import (
	"fmt"
	"os"
	"os/exec"
)

func RunDockerCompose() error {
	// Check if git is installed
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git is not installed or not in PATH")
	}

	// Check if docker-compose or docker compose is available
	dcCmd := "docker-compose"
	if _, err := exec.LookPath(dcCmd); err != nil {
		// fallback to docker compose (docker CLI v20+)
		dcCmd = "docker"
	}

	// Clone repo
	if _, err := os.Stat("zengate"); os.IsNotExist(err) {
		fmt.Println("Cloning https://github.com/myferr/zengate ...")
		cmd := exec.Command("git", "clone", "https://github.com/myferr/zengate")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repo: %w", err)
		}
	} else {
		fmt.Println("zengate directory exists, skipping clone.")
	}

	// Change directory and run docker-compose up --build
	fmt.Println("Running docker compose up --build...")
	var runCmd *exec.Cmd
	if dcCmd == "docker-compose" {
		runCmd = exec.Command("docker-compose", "up", "--build")
	} else {
		runCmd = exec.Command("docker", "compose", "up", "--build")
	}
	runCmd.Dir = "zengate"
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	err = runCmd.Run()
	if err != nil {
		return fmt.Errorf("docker compose failed: %w", err)
	}

	return nil
}
