# Raft
## Non-technical setup of Urbit on Windows

### Status 

Under construction.

### Motivation

A non-technical Windows user should be able to double-click an `.exe`, click through a GUI menu, and be up and running with a ship by the time the GUI has closed.

When complete, this will be a small application in golang which will compile to an `.exe`.
This app will initally have a command-line interface, but is planned to eventually use Fyne for the GUI and will:

1. Download and install Docker.
1. Prompt the user for configuration information (e.g., new vs existing ship, comet vs planet, where to store Urbit-related files on bare metal, etc.).
1. Pull and run the `tloncorp/urbit` image.
1. Provide the user with the `+code` required to log in to Landscape.

### Build instructions
```
env GOOS=windows GOARCH=amd64 go build -o raft.exe cmd/main/main.go
```

### Acknowledgments

Special thanks to ~botter-nidnul for inspiration, encouragement, and breaking ground with
[this excellent guide](https://gist.github.com/botter-nidnul/bc55769afe006de6f93b27390e5d1267).
