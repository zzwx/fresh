package runner

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLogColor(t *testing.T) {
	test := []struct {
		color    string
		expected string
	}{
		{color: "main", expected: "36"},
		{color: "build", expected: "33"},
		{color: "runner", expected: "32"},
		{color: "watcher", expected: "35"},
		{color: "app", expected: ""},
	}

	for _, v := range test {
		actual := logColor(v.color)
		if actual != v.expected {
			t.Errorf("Expected %v, got %v (for %q)", v.expected, actual, v.color)
		}
	}
}

func TestLoadEnvSettings(t *testing.T) {

	os.Setenv(EnvPrefix+"BUILD_DELAY", "800")
	loadEnvSettings()
	if settings.BuildDelay != 800 {
		t.Errorf("Expected %d, got %d", 800, settings.BuildDelay)
	}
	os.Setenv(EnvPrefix+"DELVE", "true")
	loadEnvSettings()
	if settings.Delve != true {
		t.Errorf("Expected %v, got %v", true, settings.Delve)
	}
	EnvPrefix = "R__"
	os.Setenv("R__DELVE", "false")
	loadEnvSettings()
	if settings.Delve != false {
		t.Errorf("Expected %v, got %v", false, settings.Delve)
	}

}

func TestBuildPath(t *testing.T) {
	expected := filepath.Join(settings.TmpPath, settings.BuildName)
	if runtime.GOOS == "windows" {
		expected += ".exe"
	}
	got := buildPath()
	if got != expected {
		t.Errorf("Expected %v, got %v", expected, got)
	}

	settings.TmpPath = "\\temp"
	settings.BuildName = "/sub/b.exe"
	expected = filepath.Clean(filepath.Join(settings.TmpPath, settings.BuildName))
	if runtime.GOOS == "windows" && !strings.HasSuffix(expected, ".exe") {
		expected += ".exe"
	}
	got = buildPath()
	if got != expected {
		t.Errorf("Expected %v, got %v", expected, got)
	}

}

func TestCleanupCommaSeparatedEntries(t *testing.T) {
	settings.Ignore = "\t\rtmp,\nsettings,\nbuild,\n\r"
	expected := `tmp,
settings,
build`
	cleanupCommaSeparatedEntries()
	if settings.Ignore != expected {
		t.Errorf("Expected %v, got %v", expected, settings.Ignore)
	}
	settings.Ignore = `
tmp,

`
	expected = `tmp`
	cleanupCommaSeparatedEntries()
	if settings.Ignore != expected {
		t.Errorf("Expected %v, got %v", expected, settings.Ignore)
	}
	settings.Ignore = `
tmp, "",

`
	expected = `tmp, ""`
	cleanupCommaSeparatedEntries()
	if settings.Ignore != expected {
		t.Errorf("Expected %v, got %v", expected, settings.Ignore)
	}

}
