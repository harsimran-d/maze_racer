package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No arguments provided")
		return
	} else if len(args) < 3 {
		fmt.Println("Also provide the command to execute")
		return
	}

	var cmd *exec.Cmd
	var dir string

	switch args[1] {
	case "flutter":
		fmt.Println("Running in flutter_app")
		dir = "flutter_app"
		cmd = exec.Command("flutter", args[2:]...)
	case "go":
		fmt.Println("Running in go_backend")
		dir = "go_backend"
		cmd = exec.Command("go", args[2:]...)
	default:
		fmt.Printf("Unknown environment: %s. Please use 'flutter' or 'go'.\n", args[1])
		return
	}

	// Validate the directory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", dir)
		return
	}
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
