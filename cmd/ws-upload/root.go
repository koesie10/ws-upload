package main

import (
	"fmt"

	"github.com/koesie10/ws-upload/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

var rootCmd = &cobra.Command{
	Use: "ws-upload",
	Version: fmt.Sprintf("%s (%s at %s)",
		version.Version,
		version.Commit,
		version.BuildDate,
	),
}
