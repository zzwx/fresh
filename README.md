# Fresh

Fresh is a command-line tool that builds and (re)starts your web application every time you save a Go or template file or any desired files you specify using configuration.

It was forked from [fresh](https://github.com/gravityblast/fresh) because the author [Andrea Franz](http://gravityblast.com) set it as unmaintained.
One of the later forks was named `fresher`, and I pulled it from [Roger Russel's](https://github.com/roger-russel/fresher.git) repo.
Then I renamed it back to **fresh** because I don't see a reason why not. I'm used to the name, and I can simply replace original fresh
with my own and don't think about it.

Fresh can be useful on Windows vs `go run` because of executable is built at one location, causing the Windows Firewall to bother only once with its popup window. 

This fork is particularly addressing the following goals:

* Work with any **folder separator**, so that one configuration file acts the same way on different platforms.
* Allow quotes `"` to surround names in the list of folders, and file extensions for more control over spaces and commas in file names.
* Allow patterns like `*` and `**` in the list of "ignore" as well as other pattern symbols, defined by [filepath.Match](https://pkg.go.dev/path/filepath?tab=doc#Match) and [doublestar](https://github.com/bmatcuk/doublestar) for `**`.
* More control over ignoring:
    * Ignore sub-folders but not the folder itself and nor sub-sub-folders (using `assets/*`).
    * Ignore sub-folders and sub-sub-folders but not the folder itself (using `assets/**`).
    * Ignore wild-carded patterns (`bootstrap-*/**`)
* Allow individual files to be ignored, because if a file type is monitored, but one particular file of that type shouldn't be (for instance, because it's auto-generated), there was no way to configure that.
* Allow for `ignore` (now default) in the settings (it just bothered me to use past tense). `ignored` still works for backward-compatibility.
* Set `debug` setting `false` to remove unnecessary output.
* Check for wrong settings names.

Converting to `yaml` configuration allows for multi-line values (with at least one space padding on every line) to be used for long option values.
Also, comments are possible after `#` symbol.

## Backward Compatibility

For the most part, simple renaming of the `runner.conf` into `.fresh.yaml` should do the job.
Then start `fresh` and check all the watched and ignored folders and files. 

The biggest change from the original "fresh" is that the sub and sub-sub folders **have to be specified more precisely**.

Rename `"ignored"` to `"ignore"` if you like. 
You can split the values into several lines with at least one space for padding, due to `yaml` syntax.

In this edition, sub-folders of an ignored folder are not automatically ignored, unless a `/**` is used. 

In short,

* Specifying just `a` as ignored results in:
    * `a` will be ignored
    * `a/sub` will **not** be ignored
    * `s/sub/sub` will **not** be ignored
* Specifying just `s/*` as ignored results in:
    * `s` will **not** be ignored
    * `s/sub` will be ignored
    * `s/sub/sub` will **not** be ignored
* Specifying just `m/**` as ignored results in:
    * `m` will **not** be ignored
    * `m/sub` will be ignored
    * `m/sub/sub` will be ignored

To emulate full ignore similar to the way it worked in original `fresh`, simply comma-separate `a` and `a/**`, which will make both the folder and all the sub and sub-sub folders ignored.

## Installation

    go get github.com/zzwx/fresh

## Usage

    cd /path/to/myapp

Start fresh:

    fresh

Fresh will watch for file events, and every time you create/modify/delete a file it will build and restart the application.

If `go build` returns an error, it will log it in the "tmp" folder.

`fresh` uses `./.fresh.yaml` for configuration by default, but you may specify an alternative config file path using `-c`:

```bash
fresh -c other.yaml
```

Here is a sample config file with the default settings:

```yaml
version: 1 # [not required]
root: . # [required] the root folder where the project is
main_path: # the folder where main.go is if it was not in root. example: %root%/cmd/
tmp_path: ./tmp
build_name: runner-build
build_args: # build args
build_log: runner-build-errors.log
valid_ext: .go, .tpl, .tmpl, .html # the extension that it will be watching
no_rebuild_ext: .tpl, .tmpl, .html
ignore: # ignore both folders and files
  assets,
  tmp
build_delay: 600
colors: true
log_color_main: cyan
log_color_build: yellow
log_color_runner: green
log_color_watcher: magenta
log_color_app:
debug: true
```

More examples can be seen [here](https://github.com/zzwx/fresh/tree/master/docs/_examples)

## Changes tracking

* [Original repo](https://github.com/gravityblast/fresh/commit/0fa698148017fa2234856bdc881d9cc62517f62b)
* [Fresher repo](https://github.com/roger-russel/fresher/commit/da1959ee8a25a760339c9f2c9b8160ce1105c02f)


## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
