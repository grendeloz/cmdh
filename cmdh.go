// This file centralises initialisation and logging for a style of
// cobra commander applications as used by grendeloz. It may suit or
// interest nobody else.

package cmdh

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	FlagConfigFile string
	FlagLogFile    string
	FlagLogLevel   string
	FlagVerbose    bool
	runParams      RunParameters
)

// Initialise adds global (Persistent) flags to the cobra root command
// and sets a version string. To use cmdh, you should call cmdh.Initialise
// from an init()  - probably in cmd/root.go. For example:
//
//   func init() {
//       cobra.OnInitialize(initConfig)
//       cmdh.Initialise(rootCmd, "myapp", "v0.1.0-dev")
//
func Initialise(rootCmd *cobra.Command, tool, version string) {
	runParams = NewRunParameters()
	runParams.Tool = tool
	runParams.Version = version

	// Persistent flags, global for the application.
	rootCmd.PersistentFlags().StringVar(&FlagConfigFile, "config",
		"", "config file")
	rootCmd.PersistentFlags().StringVar(&FlagLogFile, "logfile",
		"", "log file (defaults to STDERR if no file specified)")
	rootCmd.PersistentFlags().StringVar(&FlagLogLevel, "loglevel",
		"INFO", "log level")
	rootCmd.PersistentFlags().BoolVar(&FlagVerbose, "verbose",
		false, "turn on verbose messaging")
}

// Initialise and start logging. Note that this can not happen until after
// cobra flags have been parsed, assuming that we are allowing users to
// set values for logfile and loglevel. You would usually call
// StartLogging and FinishLogging from the Run func of the
// cobra.Command. For example:
//
//   var bamIndexCmd = &cobra.Command{
//       Use:   "index",
//       Short: "Tests on BAM and BAI files",
//       Long:  `Test read from BAM/BAI files.`,
//       Run: func(cmd *cobra.Command, args []string) {
//           cmdh.StartLogging()
//           bamIndexCmdRun(cmd, args)
//           cmdh.FinishLogging()
//       },
//   }
func StartLogging() {
	// Use our custom formatter
	formatter := LogFormat{}
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(&formatter)

	// Should fail if user-supplied logfile already exists
	if FlagLogFile != "" {
		file, err := os.OpenFile(FlagLogFile,
			os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
		if err == nil {
			log.SetOutput(file)
		} else {
			// Using fmt and os.Exit - logging is not established yet.
			fmt.Println("unable to log to file", FlagLogFile, ":", err)
			os.Exit(1)
		}
	}

	// cobra.PersistentFlags() handles the defaulting so FlagLogLevel
	// will be set to INFO if no level was supplied by the user.
	switch strings.ToUpper(FlagLogLevel) {
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	default:
		// This can only happen if the user sets a loglevel and it's not
		// one of the expected values.
		log.Fatalf("%v is not a recognised loglevel", FlagLogLevel)
	}

	// Log key execution parameters
	log.Info("Tool: ", runParams.Tool, ` `, runParams.Version)
	log.Info("Cmdline: ", runParams.Args)
	log.Info("Host: ", runParams.HostName)
	log.Infof("User: %d (%s)", runParams.UserId, runParams.UserName)
	log.Infof("Group: %d (%s)", runParams.GroupId, runParams.GroupName)

	// Read config file (default or user-supplied)
	initConfig()
	log.Infof("Config file: %v", viper.ConfigFileUsed())

	//return true
}

// FinishLogging logs elapsed time.
func FinishLogging() {
	end := time.Now()
	elapsed := end.Sub(runParams.StartTime)
	log.Info("Elapsed time: ", elapsed)
}

// The LogFormat struct and Format function below are based on info from:
// stackoverflow questions/48971780/change-format-of-log-output-logrus

// LogFormat is a custom format for log messages (via logrus)
type LogFormat struct {
	TimestampFormat string
}

// Format method (on LogFormat) implements our custom logrus log format
func (f *LogFormat) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(entry.Time.Format(f.TimestampFormat))
	b.WriteString(" [")
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString("]")

	if entry.Message != "" {
		b.WriteString(" - ")
		b.WriteString(entry.Message)
	}

	if len(entry.Data) > 0 {
		b.WriteString(" || ")
	}
	for key, value := range entry.Data {
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteByte('{')
		fmt.Fprint(b, value)
		b.WriteString("}, ")
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// initConfig reads in config file and ENV variables if set. It is
// called from StartLogging() so users do not need to call it themselves.
func initConfig() {
	if FlagConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(FlagConfigFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
}

// Tool returns the name of the application. This relies on appropriate
// values being supplied to Initialise.
func Tool() string {
	return runParams.Tool
}

// Version returns the version of the application. This relies on
// appropriate values being supplied to Initialise.
func Version() string {
	return runParams.Version
}
