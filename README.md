# watchon

A lightweight Go development tool that watches your project files and automatically restarts a running command when changes are detected.

Inspired by tools like **nodemon**, `watchon` is designed for commands such as:

```bash
watchon go run .
```

or:

```bash
watchon npm run dev
```

## Features

* Watch files and directories using [fsnotify](https://github.com/fsnotify/fsnotify)
* Detect file creation, modification, deletion, and rename events
* Automatically restart the running development process when files change
* Reduce duplicate filesystem events
* Watch a specific directory with `--only`
* Pass commands and their arguments directly to the process being watched

## Installation

Clone the repository:

```bash
git clone https://github.com/NOAK-MANUEL/watchon.git
cd watchon
```

Build the binary:

```bash
go build -o watchon .
```

You can then place the binary somewhere available in your `PATH`.

## Usage

Run a Go application:

```bash
watchon go run .
```

Run a Node.js development server:

```bash
watchon npm run dev
```

Run another command:

```bash
watchon python app.py
```

The command after `watchon` is the process that will be monitored and restarted.

## Watch a Specific Directory

Use the `--only` flag to specify the directory to watch:

```bash
watchon --only ./src go run .
```

For example:

```bash
watchon --only ./backend go run .
```

This allows you to prevent unrelated files from triggering a restart.

## How It Works

`watchon` uses [fsnotify](https://github.com/fsnotify/fsnotify) to receive filesystem events from the operating system.

The basic workflow is:

```text
                watchon
                     │
                     ▼
              Watch directory
                     │
                     ▼
              File changes
                     │
                     ▼
                fsnotify
                     │
                     ▼
             Detect the event
                     │
                     ▼
            Stop running process
                     │
                     ▼
             Start it again
                     │
                     ▼
               Keep watching
```

For example:

```text
main.go changed
      │
      ▼
fsnotify detects WRITE
      │
      ▼
watchon stops the running process
      │
      ▼
watchon starts it again
```

## Duplicate Events

Filesystem notifications do not necessarily correspond one-to-one with user actions.

For example, saving a file may produce:

```text
Modified: main.go
Modified: main.go
Modified: main.go
```

An editor or operating system can generate multiple underlying filesystem events for a single save operation.

`watchon` therefore uses event timing to prevent rapid duplicate events from causing unnecessary restarts.

## Command Arguments

Commands are passed to the underlying process as separate arguments.

For example:

```bash
watchon go run .
```

is interpreted as:

```text
Executable: go
Arguments:   run .
```

Similarly:

```bash
watchon npm run dev
```

becomes:

```text
Executable: npm
Arguments:   run dev
```

## Project Structure

The project currently uses:

* **Go** — application language
* **Cobra** — CLI framework
* **fsnotify** — filesystem event notifications

## Development

Clone the repository and install dependencies:

```bash
go mod tidy
```

Run the project:

```bash
go run . go run .
```

Build it:

```bash
go build .
```

## Project Status

`watchon` is currently under active development.

The goal is to provide a simple, lightweight Go alternative to development watchers such as nodemon, while giving more control over how files are watched and processes are restarted.

Planned improvements include:

* Better recursive directory watching
* More reliable process termination
* Debounced file events
* Ignore patterns
* File-extension filtering
* Configuration files
* Improved Windows process handling
* Cross-platform process management
* Cleaner terminal output

## License

License information will be added when the project is released.
