package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

var token string
var cmdRegex string
var daemon bool

var rootCmd = &cobra.Command{
	Use:   "httpcmd",
	Short: "Run local commands via http",
	Long:  `Run local commands via http.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			_ = cmd.Help()
		} else {
			server(daemon)
		}
	},
}

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "access token")
	rootCmd.PersistentFlags().StringVarP(&cmdRegex, "regex", "r", "", "command regex")
	rootCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "run as daemon")
}

func initLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
