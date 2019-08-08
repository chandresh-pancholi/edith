package main

import (
	"fmt"
	"github.com/chandresh-pancholi/edith/api"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func main() {

	var command = &cobra.Command{
		Use: 	"edith",
		Short: 	"Platform for event",
		Long:	"edith is a platform which can be used as Async communication, Data pipeline and auditing",
		RunE: func(cmd *cobra.Command, args []string) error {

			// TODO CPU and Memory profiler implementation using pprof.
			// TODO links to follow https://austburn.me/blog/go-profile.html, https://golang.org/pkg/runtime/pprof/

			viper.SetConfigName("local")
			viper.AddConfigPath(".")

			err := viper.ReadInConfig() // Find and read the config file
			if err != nil { // Handle errors reading the config file
				panic(fmt.Errorf("Fatal error config file: %s \n", err))
			}

			server := api.NewServer()

			server.Run()

			return nil
		},
	}

	if error := command.Execute(); error != nil {
		log.Errorf("error in executing main command. %v" , error)
		os.Exit(1)
	}
}
