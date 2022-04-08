# Decentralised Forums for Democratic Discussion

## Installing Dependencies

- Install Go, here is there get started page: [https://golang.google.cn/learn/](https://golang.google.cn/learn/)
- Install wails using their get started page: [https://wails.app/gettingstarted/](https://wails.app/gettingstarted/)

## Building the Application

`$ wails build`

(Alternatively for Linux if wails is installed locally:)

`$ ~/go/bin/wails build`

This will create a binary in the build folder: `/build/dforums-app`

## Running the Application

Simply run the binary !

## Automatically Generated Files

All files and directories will be created in the working directory.

- `dfd.log` (log file)
- `dfd-config.yaml` (config file)
- `database/` (local storage directory)
