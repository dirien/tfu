package tfu

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
)

func init() {
	tfuCmd.AddCommand(versionCmd)
}

var tfuCmd = &cobra.Command{
	Use:   "tfu",
	Short: "Update Terraform provider versions",
	Run:   runTfu,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the clients version information.",
	Run:   parseBaseCommand,
}

func getVersion() string {
	if len(Version) != 0 {
		return Version
	}
	return "dev"
}

func parseBaseCommand(_ *cobra.Command, _ []string) {
	printLogo()

	fmt.Println("Version:", getVersion())
	fmt.Println("Git Commit:", GitCommit)
	os.Exit(0)
}

func Execute(version, gitCommit string) error {
	Version = version
	GitCommit = gitCommit

	if err := tfuCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func runTfu(cmd *cobra.Command, args []string) {
	printLogo()
	err := cmd.Help()
	if err != nil {
		os.Exit(0)
	}
}

func printLogo() {
	minectlLogo := color.WhiteString(tfuFigletStr)
	fmt.Println(minectlLogo)
}

const tfuFigletStr = `
 _____ _____ _   _ 
|_   _|  ___| | | |
  | | | |_  | | | |
  | | |  _| | |_| |
  |_| |_|    \___/
`
