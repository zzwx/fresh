# Fresher

Fresher is a command-line tool that builds and (re)starts your web application every time you save a Go or template file like fresh but fresher.

It has been forked from [fresh](https://github.com/gravityblast/fresher) because the author [Andrea Franz](http://gravityblast.com) set it as unmaintained.

If the web framework you are using supports the Fresher runner, it will show build errors on your browser.

It currently works with [Traffic](https://github.com/pilu/traffic), [Martini](https://github.com/codegangsta/martini) and [gocraft/web](https://github.com/gocraft/web).

## Installation

    go get github.com/roger-russel/fresher

## Usage

    cd /path/to/myapp

Start fresher:

    fresher

Fresher will watch for file events, and every time you create/modify/delete a file it will build and restart the application.

If `go build` returns an error, it will log it in the "tmp" folder.

`fresher` uses `./.fresher.yaml` for configuration by default, but you may specify an alternative config file path using `-c`:

    fresher -c other_runner.yaml

Here is a sample config file with the default settings:
    version: 1
    root:              .                        // the root folder where the project is
    main_path:                                  // the folder where main.go is if it was not in root. exemple: /cmd/
    tmp_path:          ./tmp
    build_name:        runner-build
    build_args:                                 // build args
    build_log:         runner-build-errors.log
    valid_ext:         .go, .tpl, .tmpl, .html  // the extension that it will be watching
    no_rebuild_ext:    .tpl, .tmpl, .html
    ignored:           assets, tmp              // ignorade folders
    build_delay:       600
    colors:            1
    log_color_main:    cyan
    log_color_build:   yellow
    log_color_runner:  green
    log_color_watcher: magenta
    log_color_app:

More examples can be see [here](./docs/_examples/.)

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
