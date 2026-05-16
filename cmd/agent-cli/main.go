package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type agentInfo struct {
	Name        string
	Description string
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{Use: "agent-kit"}
	rootCmd.AddCommand(newCreateCmd())
	return rootCmd
}

func newCreateCmd() *cobra.Command {
	var info agentInfo

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Register a new ERC-8004 agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			collectedInfo, err := collectAgentInfo(cmd, info)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "\n--- Agent Registration Details ---")
			fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\n", collectedInfo.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", collectedInfo.Description)
			return nil
		},
	}

	createCmd.Flags().SetNormalizeFunc(func(_ *pflag.FlagSet, name string) pflag.NormalizedName {
		if name == "desc" {
			name = "description"
		}
		return pflag.NormalizedName(name)
	})
	createCmd.Flags().StringVarP(&info.Name, "name", "n", "", "Name of the agent")
	createCmd.Flags().StringVarP(&info.Description, "description", "d", "", "Description of the agent")

	return createCmd
}

func collectAgentInfo(cmd *cobra.Command, info agentInfo) (agentInfo, error) {
	reader := bufio.NewReader(cmd.InOrStdin())
	writer := cmd.OutOrStdout()

	name, err := promptRequiredField(reader, writer, "Agent Name", info.Name)
	if err != nil {
		return agentInfo{}, err
	}

	description, err := promptRequiredField(reader, writer, "Agent Description", info.Description)
	if err != nil {
		return agentInfo{}, err
	}

	return agentInfo{
		Name:        name,
		Description: description,
	}, nil
}

func promptRequiredField(reader *bufio.Reader, writer io.Writer, label string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value != "" {
		return value, nil
	}

	for {
		fmt.Fprintf(writer, "Enter %s: ", label)
		text, err := reader.ReadString('\n')
		isEOF := errors.Is(err, io.EOF)
		if err != nil && !isEOF {
			return "", err
		}
		text = strings.TrimSpace(text)
		if isEOF {
			if text == "" {
				return "", fmt.Errorf("unexpected EOF: %s is required", strings.ToLower(label))
			}
			return text, nil
		}
		if text != "" {
			return text, nil
		}
	}
}
