// Package config handles the configuration loading and validation for the retour application.
// It supports loading settings from both command line arguments and a TOML configuration file.
// Command line arguments take precedence over configuration file settings.
package config

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// getDefaultDBPath returns the default path for the SQLite database file
func getDefaultDBPath() string {
	return filepath.Join(".local", "share", "retour", "history.db")
}

// Mode represents the operating mode of the application.
type Mode string

const (
	// InteractiveMode indicates the application should run in an interactive shell
	InteractiveMode Mode = "interactive"
	// QueryMode indicates the application should execute a SQL query and exit
	QueryMode Mode = "query"
)

// TimeRange represents the time period over which to filter command history.
type TimeRange string

const (
	// Today filters commands executed today
	Today TimeRange = "today"
	// Yesterday filters commands executed yesterday
	Yesterday TimeRange = "yesterday"
	// LastWeek filters commands executed in the last week
	LastWeek TimeRange = "thelastweek"
	// AllTime includes all commands in history
	AllTime TimeRange = "alltime"
)

// ResultFilter represents how to filter commands based on their exit status.
type ResultFilter string

const (
	// AllResults includes both successful and failed commands
	AllResults ResultFilter = "all"
	// SuccessResults includes only commands that completed successfully
	SuccessResults ResultFilter = "success"
	// FailedResults includes only commands that failed
	FailedResults ResultFilter = "failed"
)

// Config holds all configuration for the application
// Config holds all the configuration settings for the retour application.
type Config struct {
	// Database configuration
	ConnectionString string `toml:"connection_string"`
	RetentionPeriod  string `toml:"retention_period"`

	// Command filtering
	ExclusionPatterns []string `toml:"exclusion_patterns"`
	Limit             int      `toml:"limit"`
	WorkingDirectory  string

	// Runtime options
	Mode      Mode
	Query     string
	Result    ResultFilter
	TimeRange TimeRange
}

// LoadConfig loads the configuration from both the config file and command line flags
// LoadConfig creates a new Config by combining settings from command line arguments
// and a TOML configuration file. Command line arguments take precedence over file settings.
// If the config file doesn't exist, default values will be used for file-only settings.
//
// The fsys parameter should be a filesystem containing the config file at .config/retour/config.toml
// The args parameter should be the command line arguments (including the program name as args[0])
func LoadConfig(fsys fs.FS, args []string) (*Config, error) {
	config := &Config{
		Mode:              InteractiveMode,
		Query:             "",
		Result:            AllResults,
		TimeRange:         AllTime,
		ExclusionPatterns: []string{},
	}

	configPath, err := parseCommandLine(config, args)
	if err != nil {
		return nil, err
	}

	if err := readConfig(config, fsys, configPath); err != nil {
		return nil, err
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func readConfig(config *Config, fsys fs.FS, configPath string) error {
	configFile, err := fsys.Open(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		// Set default connection string when no config file exists
		config.ConnectionString = getDefaultDBPath()
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	defer configFile.Close()

	if _, err := toml.NewDecoder(configFile).Decode(config); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}

	// Set default connection string if not specified in config
	if config.ConnectionString == "" {
		config.ConnectionString = getDefaultDBPath()
	}

	return nil
}

func parseCommandLine(config *Config, args []string) (string, error) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.Usage = usage

	flags.StringVar(&config.Query, "q", "", "SQL query to execute")
	flags.StringVar(&config.Query, "query", "", "SQL query to execute")

	flags.IntVar(&config.Limit, "l", 100, "Limit the number of results returned")
	flags.IntVar(&config.Limit, "limit", 100, "Limit the number of results returned")

	flags.StringVar(&config.WorkingDirectory, "w", "", "Filter by working directory")
	flags.StringVar(&config.WorkingDirectory, "working-directory", "", "Filter by working directory")

	result := ""
	flags.StringVar(&result, "r", string(AllResults), "Filter results (success, failed, all)")
	flags.StringVar(&result, "result", string(AllResults), "Filter results (success, failed, all)")

	timeRange := ""
	flags.StringVar(&timeRange, "t", string(AllTime), "Time range (today, yesterday, thelastweek, alltime)")
	flags.StringVar(&timeRange, "time-range", string(AllTime), "Time range (today, yesterday, thelastweek, alltime)")

	defaultConfigPath := filepath.Join(".config", "retour", "config.toml")
	configPath := ""
	flags.StringVar(&configPath, "c", defaultConfigPath, "Config file path")
	flags.StringVar(&configPath, "config", defaultConfigPath, "Config file path")

	if err := flags.Parse(args[1:]); err != nil {
		return "", fmt.Errorf("failed to parse command line flags: %w", err)
	}

	config.Result = ResultFilter(result)
	config.TimeRange = TimeRange(timeRange)
	if config.Query != "" {
		config.Mode = QueryMode
	}

	// Check if config file exists only if explicitly specified
	if configPath != defaultConfigPath {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return "", fmt.Errorf("config file %q does not exist", configPath)
		}
	}

	return configPath, nil
}

func validateConfig(config *Config) error {
	switch config.Mode {
	case InteractiveMode, QueryMode:
		// valid
	default:
		return fmt.Errorf("invalid mode: %s", config.Mode)
	}

	switch config.TimeRange {
	case Today, Yesterday, LastWeek, AllTime:
		// valid
	default:
		return fmt.Errorf("invalid time range: %s", config.TimeRange)
	}

	switch config.Result {
	case SuccessResults, FailedResults, AllResults:
		// valid
	default:
		return fmt.Errorf("invalid result filter: %s", config.Result)
	}

	if config.WorkingDirectory != "" {
		if _, err := os.Stat(config.WorkingDirectory); err != nil {
			return fmt.Errorf("invalid working directory: %w", err)
		}
	}

	if config.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0, got %d", config.Limit)
	}

	if config.ConnectionString == "" {
		return errors.New("connection string is empty")
	}

	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `Retour - Command History Manager

Usage:
  retour [options]

Options:
  -q, --query string      Execute a SQL query on the command history
  -r, --result string     Filter results by execution status (success|failed|all) [default: all]
  -t, --time-range string Time range to search (today|yesterday|thelastweek|alltime) [default: alltime]
  -c, --config string     Config file path [default: $HOME/.config/retour/config.toml]
  -l, --limit int         Limit the number of results returned [default: 100]
  -w, --working-directory Filter by working directory
  -h, --help              Show this help message

Examples:
  retour                           # Interactive mode
  retour -q "SELECT * FROM cmds"   # Query mode
  retour -r failed                 # Show failed commands
  retour -t today -r success       # Show today's successful commands
`)
}
