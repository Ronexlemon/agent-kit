package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCreateUsesFlagValues(t *testing.T) {
	cmd := newRootCmd()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{
		"create",
		"--name", "Atlas",
		"--description", "Handles agent workflows",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := output.String()
	if !strings.Contains(got, "Name: Atlas") {
		t.Fatalf("expected output to contain agent name, got %q", got)
	}
	if !strings.Contains(got, "Description: Handles agent workflows") {
		t.Fatalf("expected output to contain agent description, got %q", got)
	}
}

func TestCreatePromptsUntilRequiredFieldsAreProvided(t *testing.T) {
	cmd := newRootCmd()
	var output bytes.Buffer
	cmd.SetIn(strings.NewReader("\nAtlas\n\nHandles agent workflows\n"))
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{"create"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := output.String()
	if strings.Count(got, "Enter Agent Name: ") < 2 {
		t.Fatalf("expected agent name prompt at least twice, got %q", got)
	}
	if strings.Count(got, "Enter Agent Description: ") < 2 {
		t.Fatalf("expected agent description prompt at least twice, got %q", got)
	}
	if !strings.Contains(got, "Name: Atlas") || !strings.Contains(got, "Description: Handles agent workflows") {
		t.Fatalf("expected collected agent details in output, got %q", got)
	}
}
