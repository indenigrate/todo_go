# CLI Todo List Program

This is a simple command-line interface (CLI) todo list application written in Go. It allows users to manage their todos by adding, editing, deleting, toggling, and listing them.

## Features

- Add new todos with a title
- Edit existing todos by index
- Delete todos by index
- Toggle the completion status of todos
- List all current todos

## Requirements

- Go 1.18 or later

## Installation
## Installation Options

You have two installation methods:

- **Local Build**: Build and run the application locally.
- **System-Wide Installation**: Install the application system-wide for easier access.

Choose **local** if you want to test quickly, or **system-wide** for a more permanent setup.

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/indenigrate/todo_go.git
   cd todo_go
2. **Install dependencies**:
    ```bash
    go mod vendor
    go mod tidy
3. **Local build**:
    ```bash
    go build -o todo
    ./todo -h

4. **Install system wide**:
    ```bash
    go install
    todo_go -h

