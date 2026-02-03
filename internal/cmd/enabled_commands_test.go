package cmd

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestParseEnabledCommands(t *testing.T) {
	allow := parseEnabledCommands("calendar, tasks ,Gmail")
	if !allow["calendar"] || !allow["tasks"] || !allow["gmail"] {
		t.Fatalf("unexpected allow map: %#v", allow)
	}
}

type mockKongContext struct {
	command string
}

func (m *mockKongContext) Command() string {
	return m.command
}

func TestEnforceDisabledCommands(t *testing.T) {
	tests := []struct {
		name     string
		disabled string
		command  string
		wantErr  bool
	}{
		{
			name:     "empty disabled allows all",
			disabled: "",
			command:  "gmail send",
			wantErr:  false,
		},
		{
			name:     "whitespace-only disabled allows all",
			disabled: "   ",
			command:  "gmail send",
			wantErr:  false,
		},
		{
			name:     "exact match blocks command",
			disabled: "gmail.send",
			command:  "gmail send",
			wantErr:  true,
		},
		{
			name:     "parent blocks all children",
			disabled: "gmail",
			command:  "gmail send",
			wantErr:  true,
		},
		{
			name:     "parent blocks nested children",
			disabled: "gmail",
			command:  "gmail messages list",
			wantErr:  true,
		},
		{
			name:     "specific subcommand does not block sibling",
			disabled: "gmail.send",
			command:  "gmail search",
			wantErr:  false,
		},
		{
			name:     "specific subcommand does not block different parent",
			disabled: "gmail.messages",
			command:  "gmail send",
			wantErr:  false,
		},
		{
			name:     "case insensitive matching",
			disabled: "Gmail.SEND",
			command:  "gmail send",
			wantErr:  true,
		},
		{
			name:     "multiple disabled commands",
			disabled: "gmail.send,calendar.delete",
			command:  "gmail send",
			wantErr:  true,
		},
		{
			name:     "multiple disabled commands - second match",
			disabled: "gmail.send,calendar.delete",
			command:  "calendar delete",
			wantErr:  true,
		},
		{
			name:     "multiple disabled commands - no match",
			disabled: "gmail.send,calendar.delete",
			command:  "gmail search",
			wantErr:  false,
		},
		{
			name:     "empty command allowed",
			disabled: "gmail.send",
			command:  "",
			wantErr:  false,
		},
		{
			name:     "three-level command blocked by middle",
			disabled: "gmail.messages",
			command:  "gmail messages list",
			wantErr:  true,
		},
		{
			name:     "three-level command blocked by exact match",
			disabled: "gmail.messages.list",
			command:  "gmail messages list",
			wantErr:  true,
		},
		{
			name:     "three-level command - sibling allowed",
			disabled: "gmail.messages.list",
			command:  "gmail messages get",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := kong.New(&struct{}{})
			if err != nil {
				t.Fatalf("failed to create parser: %v", err)
			}
			kctx, err := parser.Parse([]string{})
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			// Override the command for testing
			kctx.Model.Name = tt.command

			err = enforceDisabledCommands(kctx, tt.disabled)
			if (err != nil) != tt.wantErr {
				t.Errorf("enforceDisabledCommands() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
