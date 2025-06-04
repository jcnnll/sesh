package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jcnnll/sesh/internal/config"
)

func startSesh() {
	args := os.Args

	if len(args) == 1 {
		// No extra args, run interactive session picker
		launchInteractive()
		return
	}

	switch args[1] {
	case "add":
		path := args[2]
		if !isValidDir(path) {
			fmt.Fprintf(os.Stderr, "Invalid path: %s\n", path)
			os.Exit(1)
		}
		if err := config.AddPath(path); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add path: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Path added:", path)

	case "remove":
		path := args[2]
		if !isValidDir(path) {
			fmt.Fprintf(os.Stderr, "Invalid path: %s\n", path)
			os.Exit(1)
		}
		if err := config.RemovePath(path); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove path: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Path removed:", path)

	case "help":
		printHelp()

	default:
		// If the first argument is a valid path, launch a tmux session
		if isValidDir(args[1]) {
			if err := launchTmux(args[1]); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to launch session: %v\n", err)
				os.Exit(1)
			}
			return
		}
		fmt.Println("Unknown command. Run `sesh help` for usage.")
		os.Exit(1)
	}
}

func launchInteractive() {
	projectPath := getProjectPath()

	if projectPath == "" {
		os.Exit(0)
	}

	if err := launchTmux(projectPath); err != nil {
		fmt.Printf("Failed to launch tmux session: (Path: %s)%v\n", projectPath, err)
		os.Exit(1)
	}
}

func getProjectPath() string {
	paths, err := config.GetPaths()
	if err != nil {
		fmt.Printf("failed to load paths: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		path := os.Args[1]
		_, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Invalid project path: %v\n", err)
			os.Exit(1)
		}
		return os.Args[1]
	}

	projectPath, err := selectProject(paths)
	if err != nil {
		fmt.Printf("Failed to get project: %v\n", err)
		os.Exit(1)
	}
	return projectPath
}

func selectProject(paths []string) (string, error) {

	args := append(paths, "-mindepth", "1", "-maxdepth", "1", "-type", "d", "!", "-name", ".*")
	cmdFind := exec.Command("find", args...)
	cmdFzf := exec.Command("fzf")

	stdoutPipe, err := cmdFind.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmdFzf.Stdin = stdoutPipe

	outPipe, err := cmdFzf.StdoutPipe()
	if err != nil {
		return "", err
	}

	//start processes
	if err := cmdFind.Start(); err != nil {
		return "", err
	}
	if err := cmdFzf.Start(); err != nil {
		return "", err
	}

	selected, err := io.ReadAll(outPipe)
	if err != nil {
		return "", err
	}

	// wait for process to finish
	if err := cmdFind.Wait(); err != nil {
		return "", err
	}
	if err := cmdFzf.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(string(selected)), nil
}

func launchTmux(projectPath string) error {
	editor, err := config.GetEditor()
	if err != nil {
		return fmt.Errorf("failed to get editor: %w", err)
	}

	sessionName := strings.ReplaceAll(filepath.Base(projectPath), ".", "_")

	// if session does not exist creat it detached
	if err := exec.Command("tmux", "has-session", "-t", sessionName).Run(); err != nil {
		// Editor window
		createEditor := exec.Command(
			"tmux", "new-session", "-ds", sessionName,
			"-c", projectPath,
			"-n", "editor",
			"sh", "-c", fmt.Sprintf("%s; exec $SHELL", editor),
		)
		if err := createEditor.Run(); err != nil {
			return fmt.Errorf("failed to creste tmux session: %w", err)
		}
		// Terminal window
		createTerminal := exec.Command("tmux", "new-window", "-t", sessionName, "-c", projectPath, "-n", "terminal")
		if err := createTerminal.Run(); err != nil {
			return fmt.Errorf("failed to create terminal window: %w", err)
		}
		// Focus editor window
		_ = exec.Command("tmux", "select-window", "-t", sessionName+":editor").Run()
		_ = exec.Command("tmux", "set-option", "-t", sessionName, "set-titles", "on").Run()
		_ = exec.Command("tmux", "set-option", "-t", sessionName, "set-titles-string", "#W").Run()
	}

	var cmd *exec.Cmd
	if os.Getenv("TMUX") != "" {
		cmd = exec.Command("tmux", "switch-client", "-t", sessionName)
	} else {
		cmd = exec.Command("tmux", "attach-session", "-t", sessionName)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func printHelp() {
	fmt.Println(`Usage:
____________________________________________________________

sesh                Launch fzf to pick a session path
sesh add PATH       Add a directory to session paths
sesh remove PATH    Remove a directory from session paths
sesh help           Show this help message
____________________________________________________________
		`)
}
