package cmd

import (
	"strings"

	"github.com/alecthomas/kong"
)

func enforceEnabledCommands(kctx *kong.Context, enabled string) error {
	enabled = strings.TrimSpace(enabled)
	if enabled == "" {
		return nil
	}
	allow := parseEnabledCommands(enabled)
	if len(allow) == 0 {
		return nil
	}
	if allow["*"] || allow["all"] {
		return nil
	}
	cmd := strings.Fields(kctx.Command())
	if len(cmd) == 0 {
		return nil
	}
	top := strings.ToLower(cmd[0])
	if !allow[top] {
		return usagef("command %q is not enabled (set --enable-commands to allow it)", top)
	}
	return nil
}

func parseEnabledCommands(value string) map[string]bool {
	out := map[string]bool{}
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(strings.ToLower(part))
		if part == "" {
			continue
		}
		out[part] = true
	}
	return out
}

func enforceDisabledCommands(kctx *kong.Context, disabled string) error {
	return enforceDisabledCommandsForCommand(kctx.Command(), disabled)
}

func enforceDisabledCommandsForCommand(command, disabled string) error {
	disabled = strings.TrimSpace(disabled)
	if disabled == "" {
		return nil
	}

	denyList := parseEnabledCommands(disabled)
	if len(denyList) == 0 {
		return nil
	}

	cmdParts := strings.Fields(command)
	if len(cmdParts) == 0 {
		return nil
	}

	// Check if any prefix of the command path is disabled
	// e.g., ["gmail", "send"] checks "gmail.send", then "gmail"
	for i := len(cmdParts); i >= 1; i-- {
		prefix := strings.ToLower(strings.Join(cmdParts[:i], "."))
		if denyList[prefix] {
			return usagef("command %q is disabled (blocked by --disable-commands)", strings.Join(cmdParts[:i], " "))
		}
	}

	return nil
}
