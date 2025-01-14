package cmd

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Verbose bool
var Debug bool

var rootCmd = &cobra.Command{
	Use:   "reqlab",
	Short: "CLI-Tool to test APIs",
	Long:  "CLI-Tool to test APIs",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	now := time.Now()
	logDir := "logs/"
	logFile := fmt.Sprint(now.Format("2006-01-02"), ".log")
	logFilePath := fmt.Sprint(logDir, logFile)

	_, err := os.Stat(logDir)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}

	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer f.Close()

	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)

	log.Info("Starting the application...")

	log.Debug("Loading Flags ...")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output. (default: false)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	log.Debug("Verbose: ", Verbose)

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	log.Debug("Debug: ", Debug)
}
