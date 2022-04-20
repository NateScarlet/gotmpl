# gotmpl

[![build status](https://github.com/NateScarlet/gotmpl/workflows/go/badge.svg)](https://github.com/NateScarlet/gotmpl/actions)

Render [golang template](https://pkg.go.dev/text/template/) with command line.

[templates](./templates) folder has some example templates.

## Install

Download gotmpl from [release](https://github.com/NateScarlet/gotmpl/releases),
and put it in your PATH.

Or build from source

```shell
go get github.com/NateScarlet/gotmpl/cmd/gotmpl
```

## Usage

Render template to stdout

```shell
gotmpl ./template.gotmpl
```

Render template with data

```shell
gotmpl Key=Value ./template.gotmpl
```

Render template that file name has `=` in it.

```shell
gotmpl =./template=.gotmpl
```

Render to a file.

```shell
gotmpl -o ./output.txt ./template.gotmpl
```

Render multiple template to stdout.
first template executed, others can be used with `template` action.

```shell
gotmpl Key=Value ./template.gotmpl ./template2.gotmpl
```

Read data from stdin

```shell
echo '{"Key": "Value"}' | gotmpl -i - ./template.gotmpl
```

Read data from file

```shell
gotmpl -i data.json ./template.gotmpl
```

## Template

### Functions

- [sprig](https://masterminds.github.io/sprig/) functions
- `cwd`() current working directory.
- `args`() return command line argument list.
- `templateFiles`() return template file list.
- `__file__`() value of `-o` option, empty string if not specified.
- `absPath`(string) convert path to absolute.
- `lowerFirst`(string) lower first rune of string
- `upperFirst`(string) upper first rune of string

### Data

The "." object is a `map[string]interface{}`, so you can use `set` function to modify it.

There are some default data for convenient:

- `Name`: `index (__file__ | absPath | base | split ".") 0` if `__file__` not empty
- `Package` `__file__ | absPath | dir | base` if `__file__` not empty

Data merged in order:

- default data
- `-i` file data
- inline key value pair
