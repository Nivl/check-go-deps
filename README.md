# check-deps

[![Go Report Card](https://goreportcard.com/badge/github.com/nivl/check-deps)](https://goreportcard.com/report/github.com/nivl/check-deps)

Check dep is a program that parses the `go.mod` a of package (current
directory) and prints all the modules that need to be updated, or that haven't
received any updates in the last 6 months (they might no longer be maintained).

## Install

`go get -u github.com/Nivl/check-deps`

## usage

`check-deps [flags]`

| Flag              | Description                                             |
| ----------------- | ------------------------------------------------------- |
| --check-old       | check for modules without updates for the last 6 months |
| --ignore, -i      | coma separated list of packages to ignore               |
| --check-indirects | check indirect modules                                  |
