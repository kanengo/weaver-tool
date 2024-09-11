package version

import (
	"context"
	"flag"
	"fmt"
	"github.com/ServiceWeaver/weaver/runtime/tool"
	"runtime"
	"runtime/debug"
)

func VersionCmd(toolName string) *tool.Command {
	return &tool.Command{
		Name:        "version",
		Flags:       flag.NewFlagSet("version", flag.ContinueOnError),
		Description: fmt.Sprintf("Show %q version", toolName),
		Help:        fmt.Sprintf("Usage:\n  %s version", toolName),
		Fn: func(context.Context, []string) error {
			v, err := SelfVersion()
			if err != nil {
				return err
			}
			fmt.Printf("%s %s %s/%s\n", toolName, v, runtime.GOOS, runtime.GOARCH)
			return nil
		},
	}
}

// SelfVersion returns the version of the running tool binary.
func SelfVersion() (string, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		// Should never happen.
		return "", fmt.Errorf("tool binary must be built from a module")
	}
	return info.Main.Version, nil
}
