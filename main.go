package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"regexp"
)

type Client struct {
	Class     string `json:"class"`
	Workspace struct {
		ID int `json:"id"`
	} `json:"workspace"`
}

type ActiveWorkspace struct {
	ID int `json:"id"`
}

func getEmacsWorkspace() (int, bool) {
	cmd := exec.Command("hyprctl", "clients", "-j")
	output, err := cmd.Output()
	if err != nil {
		return 0, false
	}

	var clients []Client
	if err := json.Unmarshal(output, &clients); err != nil {
		return 0, false
	}

	for _, client := range clients {
		fmt.Println(client)
		if strings.Contains(client.Class, "emacs") {
			fmt.Println(client.Workspace.ID)
			return client.Workspace.ID, true
		}
	}
	return 0, false
}

func getCurrentWorkspace() (int, error) {
	cmd := exec.Command("hyprctl", "activeworkspace", "-j")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var workspace ActiveWorkspace
	if err := json.Unmarshal(output, &workspace); err != nil {
		return 0, err
	}
	return workspace.ID, nil
}

func executeEmacsCommand(command string) error {
	cmd := exec.Command("emacsclient", "-n", "-e", command)
	return cmd.Run()
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "elauncher",
		Short: "Emacs Launcher for Hyprland",
		Args:  cobra.MinArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			emacsCommand := strings.Join(args[0:], " ")

			targetWS, found := getEmacsWorkspace()
			if !found {
				if err := executeEmacsCommand(emacsCommand); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: emacsclient failed: %v\n", err)
					os.Exit(1)
				}
				return
			}

			currentWS, err := getCurrentWorkspace()
			if err != nil {
				fmt.Println("Been here")
				executeEmacsCommand(emacsCommand)
				return
			}

			if currentWS == targetWS {
				fmt.Println("Current workspace matches target workspace")
				executeEmacsCommand(emacsCommand)
				return
			}

			switchCmd := exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", targetWS))
			if err := switchCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: workspace switch failed: %v\n", err)
				executeEmacsCommand(emacsCommand)
				return
			}

			maxWait := 500 * time.Millisecond
			pollInterval := 2 * time.Millisecond // Start very fast
			deadline := time.Now().Add(maxWait)

			for time.Now().Before(deadline) {
				current, err := getCurrentWorkspace()
				if err == nil && current == targetWS {
					executeEmacsCommand(emacsCommand)
					return
				}

				time.Sleep(pollInterval)

				pollInterval *= 2
				if pollInterval > 10*time.Millisecond {
					pollInterval = 10 * time.Millisecond
				}
			}

			executeEmacsCommand(emacsCommand)

		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
