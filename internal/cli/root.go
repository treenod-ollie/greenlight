package cli

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	appVersion string
	verbose    bool
)

// ErrThreshold is returned by scan commands when --exit-code is set and the
// findings cross the failure threshold. main() maps it to a non-zero exit status
// without printing anything — the report is already the user-facing output.
var ErrThreshold = errors.New("findings exceeded --exit-code threshold")

var purple = color.New(color.FgHiMagenta)
var dim = color.New(color.Faint)

var rootCmd = &cobra.Command{
	Use:   "greenlight",
	Short: "Pre-submission compliance scanner for the App Store and Google Play",
	Long: fmt.Sprintf(`%s

Greenlight scans your app against Apple's App Store Review Guidelines and
Google Play's Developer Program Policies before you submit, catching
rejection risks so you ship with confidence.

This is an OFFLINE-ONLY build: every command reads local files and nothing
else. There is no network code left in the binary — no account, no upload,
no telemetry. Source and build artifacts never leave the machine.

Get started:
  greenlight preflight .          Run ALL checks — one command, fully offline
  greenlight preflight . --ipa X  Include IPA binary analysis
  greenlight playscan .           Google Play policies + target API deadline
  greenlight guidelines search    Browse Apple's review guidelines`,
		purple.Sprint("greenlight — know before you submit.")),
}

func SetVersion(v string) {
	appVersion = v
	// Enables `greenlight --version` (cobra auto-adds the flag when set).
	rootCmd.Version = resolveVersion()
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Errors are printed by main() (so --exit-code's sentinel can exit quietly);
	// usage on a runtime error is just noise.
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	// Match the `version` subcommand's wording ("greenlight <v>") for `--version`.
	rootCmd.SetVersionTemplate("greenlight {{.Version}}\n")

	rootCmd.AddCommand(guidelinesCmd)
	rootCmd.AddCommand(versionCmd)
}
