package runner

import (
	"testing"
)

func TestIsWatchedExt(t *testing.T) {
	settings.ValidExt = `.go, "", .tpl, ".tmpl", .html, ".g,m"`
	tests := []struct {
		file     string
		expected bool
	}{
		{"test.go", true},
		{"test.tpl", true},
		{"test.tmpl", true},
		{"test.html", true},
		{"test.css", false},
		{"test.g,m", true},          // comma in file type
		{"test-no-extension", true}, // "" allows for files that end with no extension
	}
	// valid_ext: .go, .tpl, .tmpl, .html

	for _, test := range tests {
		actual := isWatchedExt(test.file)

		if actual != test.expected {
			t.Errorf("Expected %v, got %v (%q)", test.expected, actual, test.file)
		}
	}
}

func TestShouldRebuild(t *testing.T) {
	settings.ValidExt = ".go, .tpl, .tmpl, .html, .ts, .scss, .css, .tsx, .json, .txt"
	settings.NoRebuildExt = ".tpl, .tmpl, .html"
	settings.Ignore = "src/main\\script.ts, build\\bundle.js, tmp\\*, src\\main\\node_modules, src\\main\\node_modules\\*"
	tests := []struct {
		eventName string
		expected  bool
	}{
		{`"src\\main\\script.ts": MODIFY`, false},
		{`"src/main/script.ts": MODIFY`, false},
		{`"src\\main/validscript.ts": MODIFY`, true},
		{`"tmp/somefile.go": MODIFY`, false},
		{`"tmp/deep/somefile.go": MODIFY`, true},
		{`"test.go": MODIFY`, true},
		{`"test.tpl": MODIFY`, false},
		{`"test.tmpl": DELETE`, false},
		{`"unknown.extension": DELETE`, true},
		{`"no_extension": ADD`, true},
		{`"./a/path/test.go": MODIFY`, true},
	}

	for _, test := range tests {
		actual := shouldRebuild(test.eventName)

		if actual != test.expected {
			t.Errorf("Expected %v, got %v (event was %v)", test.expected, actual, test.eventName)
		}
	}
}

func TestIsIgnored(t *testing.T) {
	settings.Ignore = "a, b, b/*, s/*, m/**, tmp, tmp/*, build/file.ts, assets/*, app/controllers/, \"app/controllers/*\", app/views/*, \"dir with space\", \"dir\\sub\""
	tests := []struct {
		dir      string
		expected bool
	}{
		{"a", true}, // "a" means a, but not sub-folders nor sub-sub-folders of a
		{"a/sub", false},
		{"a/sub/sub", false},

		{"b", true}, // "b, b/*" means b, and b sub-folders, but not b sub-sub-folders
		{"b/sub", true},
		{"b/sub/sub", false},

		{"s", false}, // "s/*" means s sub-folders but not sub-sub-folders too, and not the s itself
		{"s/sub", true},
		{"s/sub/sub", false},

		{"m", false}, // "m/**" means m sub-folders and all its sub-folders too, but not the m itself
		{"m/sub", true},
		{"m/sub/sub", true},

		{"assets", false}, // "assets/*" is in the list, but not the "assets".
		{"assets/node_modules", true},
		{"./build\\file.ts", true}, // regular files can be ignored.
		{"build/anotherFile.ts", false},
		{"tmp", true},
		{"tmp/pid", true},
		{"app", false},
		{"app/controllers", true},
		{"app/controllers/user", true},
		{"app/views", false},         // because "app/views" is not listed.
		{"app/views/ is good", true}, // because "app/views/*" is listed.
		{"dir with space", true},
		{"./dir with space", true},
		{"./dir/sub", true},
		{"\\dir\\sub", false}, // not the root dir.
		{"dir\\sub", true},    // now fine.
	}
	for _, test := range tests {
		actual := isIgnored(test.dir)
		if actual != test.expected {
			t.Errorf("Expected %v, got %v (for %q)", test.expected, actual, test.dir)
		}
	}
}
