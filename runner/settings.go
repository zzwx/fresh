package runner

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	envSettingsPrefix   = "RUNNER_"
	mainSettingsSection = "Settings"
)

var settings = map[string]string{
	"version":           "1",
	"config_path":       "./.fresher.yaml",
	"root":              ".",
	"main_path":         "",
	"tmp_path":          "./tmp",
	"build_name":        "runner-build",
	"build_args":        "",
	"build_log":         "runner-build-errors.log",
	"valid_ext":         ".go, .tpl, .tmpl, .html",
	"no_rebuild_ext":    ".tpl, .tmpl, .html",
	"ignored":           "assets, tmp",
	"build_delay":       "600",
	"colors":            "1",
	"log_color_main":    "cyan",
	"log_color_build":   "yellow",
	"log_color_runner":  "green",
	"log_color_watcher": "magenta",
	"log_color_app":     "",
	"delve":             "false",
	"delve_args":        "",
}

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

func logColor(logName string) string {
	settingsKey := fmt.Sprintf("log_color_%s", logName)
	colorName := settings[settingsKey]

	return colors[colorName]
}

func loadEnvSettings() {
	for key := range settings {
		envKey := fmt.Sprintf("%s%s", envSettingsPrefix, strings.ToUpper(key))
		if value := os.Getenv(envKey); value != "" {
			settings[key] = value
		}
	}

}

func loadRunnerConfigSettings() {

	cfgPath := configPath()

	if _, err := os.Stat(cfgPath); err != nil {
		panic(err)
	}

	logger.Printf("Loading settings from %s", cfgPath)

	file, err := ioutil.ReadFile(cfgPath)

	if err != nil {
		panic(err)
	}

	var givenSettings map[string]string

	yaml.Unmarshal(file, &givenSettings)

	if givenSettings["version"] == "" {
		log.Fatalln("no version was setted on config yaml file.")
	}

	for key, value := range givenSettings {
		settings[key] = value
	}

}

func initSettings() {
	loadEnvSettings()
	loadRunnerConfigSettings()
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func root() string {
	return settings["root"]
}

func mainPath() string {
	return settings["main_path"]
}

func tmpPath() string {
	return settings["tmp_path"]
}

func buildName() string {
	return settings["build_name"]
}

func buildPath() string {
	p := filepath.Join(tmpPath(), buildName())
	if runtime.GOOS == "windows" && filepath.Ext(p) != ".exe" {
		p += ".exe"
	}
	return p
}

func buildArgs() string {
	return settings["build_args"]
}

func buildErrorsFileName() string {
	return settings["build_log"]
}

func buildErrorsFilePath() string {
	return filepath.Join(tmpPath(), buildErrorsFileName())
}

func configPath() string {
	return settings["config_path"]
}

func buildDelay() time.Duration {
	value, _ := strconv.Atoi(settings["build_delay"])

	return time.Duration(value)
}

func mustUseDelve() bool {

	var b bool
	var err error

	if b, err = strconv.ParseBool(settings["delve"]); err != nil {
		panic(err)
	}

	return b

}

func delveArgs() string {
	return settings["delve_args"]
}
