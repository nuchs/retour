package main_test

import (
	"testing"
	"testing/fstest"

	rt "github.com/nuchs/retour"
)

func TestTimeRange(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want rt.TimeRange
	}{
		{
			name: "Default",
			args: []string{"cmd"},
			want: rt.AllTime,
		},
		{
			name: "Short form today",
			args: []string{"cmd", "-t", "today"},
			want: rt.Today,
		},
		{
			name: "Short form yesterday",
			args: []string{"cmd", "-t", "yesterday"},
			want: rt.Yesterday,
		},
		{
			name: "Short form lastweek",
			args: []string{"cmd", "-t", "thelastweek"},
			want: rt.LastWeek,
		},
		{
			name: "Short form alltime",
			args: []string{"cmd", "-t", "alltime"},
			want: rt.AllTime,
		},
		{
			name: "Long form today",
			args: []string{"cmd", "--time-range", "today"},
			want: rt.Today,
		},
		{
			name: "Long form yesterday",
			args: []string{"cmd", "--time-range", "yesterday"},
			want: rt.Yesterday,
		},
		{
			name: "Long form lastweek",
			args: []string{"cmd", "--time-range", "thelastweek"},
			want: rt.LastWeek,
		},
		{
			name: "Long form alltime",
			args: []string{"cmd", "--time-range", "alltime"},
			want: rt.AllTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := rt.LoadConfig(makeConfigFile(t), tt.args)
			if err != nil {
				t.Fatalf("LoadConfig() unexpected error = %v", err)
			}

			if got := config.TimeRange; got != tt.want {
				t.Errorf("TimeRange = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want rt.ResultFilter
	}{
		{
			name: "Default",
			args: []string{"cmd"},
			want: rt.AllResults,
		},
		{
			name: "Short form all",
			args: []string{"cmd", "-r", "all"},
			want: rt.AllResults,
		},
		{
			name: "Short form success",
			args: []string{"cmd", "-r", "success"},
			want: rt.SuccessResults,
		},
		{
			name: "Short form failed",
			args: []string{"cmd", "-r", "failed"},
			want: rt.FailedResults,
		},
		{
			name: "Long form all",
			args: []string{"cmd", "--result", "all"},
			want: rt.AllResults,
		},
		{
			name: "Long form success",
			args: []string{"cmd", "--result", "success"},
			want: rt.SuccessResults,
		},
		{
			name: "Long form failed",
			args: []string{"cmd", "--result", "failed"},
			want: rt.FailedResults,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := rt.LoadConfig(makeConfigFile(t), tt.args)
			if err != nil {
				t.Fatalf("LoadConfig() unexpected error = %v", err)
			}

			if got := config.Result; got != tt.want {
				t.Errorf("Result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMode(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantMode rt.Mode
		wantSQL  string
	}{
		{
			name:     "Default",
			args:     []string{"cmd"},
			wantMode: rt.InteractiveMode,
			wantSQL:  "",
		},
		{
			name:     "Short form query",
			args:     []string{"cmd", "-q", "SELECT * FROM cmds"},
			wantMode: rt.QueryMode,
			wantSQL:  "SELECT * FROM cmds",
		},
		{
			name:     "Long form query",
			args:     []string{"cmd", "--query", "SELECT * FROM cmds"},
			wantMode: rt.QueryMode,
			wantSQL:  "SELECT * FROM cmds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := rt.LoadConfig(makeConfigFile(t), tt.args)
			if err != nil {
				t.Fatalf("LoadConfig() unexpected error = %v", err)
			}

			if got := config.Mode; got != tt.wantMode {
				t.Errorf("Mode = %v, want %v", got, tt.wantMode)
			}
			if got := config.Query; got != tt.wantSQL {
				t.Errorf("Query = %v, want %v", got, tt.wantSQL)
			}
		})
	}
}

func TestConfigFile(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		wantConn   string
		wantRet    string
		wantExcl   []string
		wantLimit  int
	}{
		{
			name:       "Empty config",
			configFile: "",
			wantConn:   ".local/share/retour/history.db",
			wantRet:    "",
			wantExcl:   []string{},
			wantLimit:  100,
		},
		{
			name: "Full config",
			configFile: `
connection_string = "test.db"
retention_period = "30d"
exclusion_patterns = ["^sudo", "^ssh"]
limit = 50
`,
			wantConn:  "test.db",
			wantRet:   "30d",
			wantExcl:  []string{"^sudo", "^ssh"},
			wantLimit: 50,
		},
		{
			name: "Partial config",
			configFile: `
connection_string = "test.db"
`,
			wantConn:  "test.db",
			wantRet:   "",
			wantExcl:  []string{},
			wantLimit: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := fstest.MapFS{".config/retour/config.toml": &fstest.MapFile{Data: []byte(tt.configFile)}}

			config, err := rt.LoadConfig(fsys, []string{"cmd"})
			if err != nil {
				t.Fatalf("LoadConfig() unexpected error = %v", err)
			}

			if got := config.ConnectionString; got != tt.wantConn {
				t.Errorf("ConnectionString = %v, want %v", got, tt.wantConn)
			}
			if got := config.RetentionPeriod; got != tt.wantRet {
				t.Errorf("RetentionPeriod = %v, want %v", got, tt.wantRet)
			}
			if got := len(config.ExclusionPatterns); got != len(tt.wantExcl) {
				t.Errorf("ExclusionPatterns length = %v, want %v", got, len(tt.wantExcl))
			} else {
				for i, want := range tt.wantExcl {
					if got := config.ExclusionPatterns[i]; got != want {
						t.Errorf("ExclusionPattern[%d] = %v, want %v", i, got, want)
					}
				}
			}
			if got := config.Limit; got != tt.wantLimit {
				t.Errorf("Limit = %v, want %v", got, tt.wantLimit)
			}
		})
	}
}

func TestBadCommandLine(t *testing.T) {
	// Test cases for configuration loading failures
	tests := []struct {
		name       string
		args       []string
		want       string
		skipConfig bool // If true, don't create the config file
	}{
		{
			name: "Invalid result filter",
			args: []string{"cmd", "-r", "invalid"},
			want: "invalid result filter: invalid",
		},
		{
			name: "Invalid time range",
			args: []string{"cmd", "-t", "invalid"},
			want: "invalid time range: invalid",
		},
		{
			name:       "Bad config path",
			args:       []string{"cmd", "-c", "invalid"},
			want:       "config file \"invalid\" does not exist",
			skipConfig: true,
		},
		{
			name: "Invalid limit",
			args: []string{"cmd", "--limit", "0"},
			want: "limit must be greater than 0, got 0",
		},
		{
			name: "Invalid working directory",
			args: []string{"cmd", "--working-directory", "/nonexistent/path"},
			want: "invalid working directory: stat /nonexistent/path: no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a valid config file for command line validation tests
			fsys := makeConfigFile(t)
			if tt.skipConfig {
				fsys = &fstest.MapFS{}
			}

			_, err := rt.LoadConfig(fsys, tt.args)
			if err == nil {
				t.Fatal("Want error, got nil")
			}
			if errMsg := err.Error(); errMsg != tt.want {
				t.Errorf("Got = %v, want %v", errMsg, tt.want)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want int
	}{
		{
			name: "Default values",
			args: []string{"cmd"},
			want: 100,
		},
		{
			name: "Short form limit",
			args: []string{"cmd", "-l", "50"},
			want: 50,
		},
		{
			name: "Long form limit",
			args: []string{"cmd", "--limit", "25"},
			want: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := makeConfigFile(t)

			config, err := rt.LoadConfig(fsys, tt.args)
			if err != nil {
				t.Fatalf("LoadConfig() unexpected error = %v", err)
			}

			if got := config.Limit; got != tt.want {
				t.Errorf("Limit = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestLimitAndWorkingDir(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "Default values",
			args: []string{"cmd"},
			want: "",
		},
		{
			name: "Short form working directory",
			args: []string{"cmd", "-w", "/tmp"},
			want: "/tmp",
		},
		{
			name: "Long form working directory",
			args: []string{"cmd", "--working-directory", "/tmp"},
			want: "/tmp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := makeConfigFile(t)

			config, err := rt.LoadConfig(fsys, tt.args)
			if err != nil {
				t.Fatalf("LoadConfig() unexpected error = %v", err)
			}

			if got := config.WorkingDirectory; got != tt.want {
				t.Errorf("WorkingDirectory = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeConfigFile(t *testing.T) *fstest.MapFS {
	t.Helper()
	fsys := fstest.MapFS{}
	fsys[".config/retour/config.toml"] = &fstest.MapFile{
		Data: []byte(`
connection_string = "test.db"
retention_period = "30d"
`),
	}
	return &fsys
}
