package runner

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigPath = "./.fresh.yaml"
)

var ConfigPath string
var EnvPrefix = "RUNNER_"

type Settings struct {
	Version         string `yaml:"version"`
	Root            string `yaml:"root"`
	MainPath        string `yaml:"main_path"`
	TmpPath         string `yaml:"tmp_path"`
	BuildName       string `yaml:"build_name"`
	BuildArgs       string `yaml:"build_args"`
	RunArgs         string `yaml:"run_args"`
	BuildLog        string `yaml:"build_log"`
	ValidExt        string `yaml:"valid_ext"`
	NoRebuildExt    string `yaml:"no_rebuild_ext"`
	Ignore          string `yaml:"ignore"`
	BuildDelay      string `yaml:"build_delay"` // Number: Nanoseconds, otherwise - parse Duration
	Colors          bool   `yaml:"colors"`
	LogColorMain    string `yaml:"log_color_main"`
	LogColorBuild   string `yaml:"log_color_build"`
	LogColorRunner  string `yaml:"log_color_runner"`
	LogColorWatcher string `yaml:"log_color_watcher"`
	LogColorApp     string `yaml:"log_color_app"`
	Delve           bool   `yaml:"delve"`
	DelveArgs       string `yaml:"delve_args"`
	Debug           bool   `yaml:"debug"`
}

var settings Settings

func init() {
	ConfigPath = DefaultConfigPath

	settings.Version = "1"
	settings.Root = "."
	// settings.MainPath
	settings.TmpPath = "./tmp"
	settings.BuildName = "runner-build"
	// settings.BuildArgs
	// settings.RunArgs
	settings.BuildLog = "runner-build-errors.log"
	settings.ValidExt = ".go, .tpl, .tmpl, .html"
	settings.NoRebuildExt = ".tpl, .tmpl, .html"
	settings.Ignore = "assets, tmp/*"
	settings.BuildDelay = "600"
	settings.Colors = true
	settings.LogColorMain = "cyan"
	settings.LogColorBuild = "yellow"
	settings.LogColorRunner = "green"
	settings.LogColorWatcher = "magenta"
	//settings.LogColorApp
	//settings.Delve
	//settings.DelveArgs
	settings.Debug = true
}

/*var settings = map[string]string{
	"version":           "1",
	"config_path":       DefaultConfigPath,
	"root":              ".",
	"main_path":         "",
	"tmp_path":          "./tmp",
	"build_name":        "runner-build",
	"build_args":        "",
	"build_log":         "runner-build-errors.log",
	"valid_ext":         ".go, .tpl, .tmpl, .html",
	"no_rebuild_ext":    ".tpl, .tmpl, .html",
	"ignore":            "assets, tmp/*",
	"build_delay":       "600",
	"colors":            "true",
	"log_color_main":    "cyan",
	"log_color_build":   "yellow",
	"log_color_runner":  "green",
	"log_color_watcher": "magenta",
	"log_color_app":     "",
	"delve":             "false",
	"delve_args":        "",
	"debug":             "true",
}*/

var colors = map[string]string{
	"reset":          "0",
	"black":          "30",
	"red":            "31",
	"green":          "32",
	"yellow":         "33",
	"blue":           "34",
	"magenta":        "35",
	"cyan":           "36",
	"white":          "37",
	"bold_black":     "30;1",
	"bold_red":       "31;1",
	"bold_green":     "32;1",
	"bold_yellow":    "33;1",
	"bold_blue":      "34;1",
	"bold_magenta":   "35;1",
	"bold_cyan":      "36;1",
	"bold_white":     "37;1",
	"bright_black":   "30;2",
	"bright_red":     "31;2",
	"bright_green":   "32;2",
	"bright_yellow":  "33;2",
	"bright_blue":    "34;2",
	"bright_magenta": "35;2",
	"bright_cyan":    "36;2",
	"bright_white":   "37;2",
}

// logColor can return empty string if color is not found
func logColor(logName string) string {
	var clr string
	switch logName {
	case "main":
		clr = settings.LogColorMain
	case "build":
		clr = settings.LogColorBuild
	case "runner":
		clr = settings.LogColorRunner
	case "watcher":
		clr = settings.LogColorWatcher
	case "app":
		clr = settings.LogColorApp
	default:
		panic("unknown color type: " + logName)
	}
	return colors[clr]
}

// tagDetails takes the struct field from a Field's Tag and splits the tag data `yaml:"ignored,omitempty",fresh:"noenv"`
// to retrieve the field name and possible omit (if omitempty is found) and noenv (if found)
func tagDetails(field reflect.StructTag) (keyName string, omit bool, noenv bool) {
	yaml_ := field.Get("yaml")
	fresh := field.Get("fresh")
	split := strings.Split(yaml_, ",")
	if len(split) == 0 {
		return
	}
	keyName = split[0]
	for i := 1; i < len(split); i++ {
		switch split[i] {
		case "omitempty":
			omit = true
		}
	}
	split = strings.Split(fresh, ",")
	for i := 0; i < len(split); i++ {
		switch split[i] {
		case "noenv":
			noenv = true
		}
	}
	return
}

func loadEnvSettings() {
	envKey := EnvPrefix + "CONFIG_PATH"
	if value := os.Getenv(envKey); value != "" {
		ConfigPath = value
	}
	t := reflect.TypeOf(Settings{})
	v := reflect.ValueOf(&settings).Elem()
	for i := 0; i < t.NumField(); i++ {
		keyName, _, noenv := tagDetails(t.Field(i).Tag)
		if keyName != "" && !noenv {
			envKey := fmt.Sprintf("%s%s", EnvPrefix, strings.ToUpper(keyName))
			if value := os.Getenv(envKey); value != "" {
				field := v.FieldByName(t.Field(i).Name)
				if field.CanSet() {
					if field.Kind() == reflect.String {
						field.SetString(value)
					} else if field.Kind() == reflect.Bool {
						if value == "1" || value == "true" {
							field.SetBool(true)
						} else {
							field.SetBool(false)
						}
					} else if field.Kind() == reflect.Uint {
						if u, err := strconv.ParseUint(value, 10, 64); err == nil {
							field.SetUint(u)
						}
					}
					//runnerLog("Set %q from env %q to %v", t.Field(i).Name, envKey, value)
				}
			}
		}
	}
	cleanupCommaSeparatedEntries()
}

// trimSpaceAndHangingComma removes a hanging comma that can be left by mistake or
// in a multiline setting for convenience of editing, but will mean
// an empty string when decoded
func trimSpaceAndHangingComma(input string) string {
	input = strings.TrimSpace(input)
	if strings.HasSuffix(input, ",") {
		input = input[0 : len(input)-1]
	}
	return input
}

func cleanupCommaSeparatedEntries() {
	settings.Ignore = trimSpaceAndHangingComma(settings.Ignore)
	settings.ValidExt = trimSpaceAndHangingComma(settings.ValidExt)
	settings.NoRebuildExt = trimSpaceAndHangingComma(settings.NoRebuildExt)
}

func loadRunnerConfigSettings() {
	cfgPath := configPath()
	if _, err := os.Stat(cfgPath); err != nil {
		mainLog("No config file found at %q. Using default settings", cfgPath)
	} else {
		mainLog("Loads settings from %q", cfgPath)
		file, err := ioutil.ReadFile(cfgPath)
		if err != nil {
			mainLog("Error reading config file %q: %v", cfgPath, err)
			os.Exit(1)
		}
		err = yaml.Unmarshal(file, &settings)
		if err != nil {
			mainLog("Error unmarshalling config file %q: %v", cfgPath, err)
			os.Exit(1)
		}
		// Check for wrong field names
		var givenSettings map[string]string
		yaml.Unmarshal(file, &givenSettings)
		t := reflect.TypeOf(Settings{})
		for key := range givenSettings {
			found := false
			for i := 0; i < t.NumField(); i++ {
				keyName, _, _ := tagDetails(t.Field(i).Tag)
				if keyName == key {
					found = true
					break
				}
			}
			if !found {
				switch key {
				case "ignored":
					// Archaic "ignored" name should be treated as "ignore"
					settings.Ignore = givenSettings["ignored"]
				default:
					mainLog("Unknown settings name: %q", key)
					os.Exit(1)
				}
			}
		}
	}
	cleanupCommaSeparatedEntries()
}

func SaveRunnerConfigSettings(configPath string) {
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(configPath)
			if err != nil {
				fmt.Printf("Cannot create file %q due to %v", configPath, err)
				os.Exit(1)
			}
			err = yaml.NewEncoder(file).Encode(settings)
			if err != nil {
				fmt.Printf("Error encoding default settings into %q due to %v", configPath, err)
				os.Exit(1)
			}
			fmt.Printf("%q generated\n", configPath)
			return
		}
	}
	fmt.Printf("%q already exists\n", configPath)
	os.Exit(1)
}

func initSettings() {
	loadEnvSettings()
	loadRunnerConfigSettings()
}

func root() string {
	return filepath.Clean(settings.Root)
}

func mainPath() string {
	// MainPath is representing a module full path, not a file path
	// No need to clean it
	return settings.MainPath
}

func tmpPath() string {
	return filepath.Clean(settings.TmpPath)
}

func buildPath() string {
	p := filepath.Join(tmpPath(), settings.BuildName)
	if runtime.GOOS == "windows" && filepath.Ext(p) != ".exe" {
		p += ".exe"
	}
	return p
}

func buildArgs() string {
	return settings.BuildArgs
}

func runArgs() string {
	return settings.RunArgs
}

func buildErrorsFileName() string {
	return filepath.Clean(settings.BuildLog) // In case a path is included
}

func buildErrorsFilePath() string {
	return filepath.Join(tmpPath(), buildErrorsFileName())
}

func configPath() string {
	return filepath.Clean(ConfigPath)
}

func buildDelay() time.Duration {
	if v, err := strconv.ParseInt(settings.BuildDelay, 10, 64); err == nil {
		if v > 0 {
			return time.Duration(v)
		}
	}
	if v, err := strconv.ParseFloat(settings.BuildDelay, 64); err == nil {
		if v > 0 {
			return time.Duration(int64(v))
		}
	}
	if v, err := time.ParseDuration(settings.BuildDelay); err == nil {
		return v
	}
	return time.Duration(600) // 600 nanoseconds
}

func mustUseDelve() bool {
	return settings.Delve
}

func delveArgs() string {
	return settings.DelveArgs
}

func isDebug() bool {
	return settings.Debug
}
