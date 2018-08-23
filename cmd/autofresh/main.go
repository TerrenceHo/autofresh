package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/TerrenceHo/autofresh"
	"github.com/TerrenceHo/autofresh/config"
	"github.com/spf13/cobra"
)

const (
	GOARCH string = runtime.GOARCH
	GOOS   string = runtime.GOOS
)

var (
	Version   string
	GitHash   string
	BuildTime string
	GoVersion string = runtime.Version()
)

var mainCmd = &cobra.Command{
	Use:   "autofresh",
	Short: "Autofresh: live reloading server",
	Run:   serve,
}

func init() {
	flags := mainCmd.Flags()

	flags.BoolP("version", "v", false, "print version of autofresh")
	flags.Bool("hidebanner", false, "hide autofresh banner")
	flags.StringP("watchman", "w", "watchman", "path to watchman executable")
	flags.StringP("build", "b", "", "set build command")
	flags.StringP("run", "r", "", "set process start command/run built program")
	flags.StringSliceP("suffixes", "s", []string{}, "array of file suffixes")
}

func main() {
	mainCmd.Execute()
}

func serve(cmd *cobra.Command, args []string) {
	conf, errs := config.LoadConfig(cmd)
	if conf.Version {
		fmt.Printf("autofresh %s %s %s/%s\n", Version, GoVersion, GOOS, GOARCH)
		fmt.Printf("git hash: %s\n", GitHash)
		fmt.Printf("built at: %s\n", BuildTime)
		os.Exit(0)
	}
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}

	autofresh.Start(conf)
}
