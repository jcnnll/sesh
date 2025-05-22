# sesh - Terminal Workspace Launcher Tool

## Purpose
`sesh` is a CLI tool that I created for personal use. It uses `fzf` to fuzzy-select a workspace
and launch a tmux session that is ready for a coding session.

This tool is heavily inspired by [ThePrimeagen](https://github.com/ThePrimeagen)'s Dev workflow
and the [tmux-sessionizer](https://github.com/ThePrimeagen/.dotfiles/blob/master/bin/.local/scripts/tmux-sessionizer)
script used to launch tmux sessions.

## Primary Behavior
| Command                | Description                                                                           |
| ---------------------- | ------------------------------------------------------------------------------------- |
| `sesh`                 | Fuzzy-select a directory from all configured workspace paths and start a tmux session |
| `sesh <path>`          | Launch a tmux session directly from the provided path (absolute or relative)          |
| `sesh add <path>`      | Add a path to the config                                                              |
| `sesh remove [<path>]` | Remove a path from the config                                                         |

## Installation
### Clone the Repository

```sh
git clone https://github.com/yourusername/sesh.git
cd sesh
```

### Build
Ensure you have Go installed (version 1.20+ recommended).

Then run:
```sh
go build -o sesh
```

### Install

To install sesh globally (in your $PATH):
```sh
go install .
```

## Configuration
`sesh` uses a JSON configuration file to define workspace search paths and your preferred editor.
The config file is automatically created on first run if it doesn't exist.

By default, the config file is located at:
```bash
$HOME/.config/sesh/config.json
```

Config format:
```json
{
  "paths": [
    "/Users/jay/workspace/github.com",
    "/Users/jay/personal"
  ],
  "editor": "nvim"
}
```
- **paths**
  A list of absolute paths to directories where sesh will look for projects.
- **editor**
  Default to `nvim` but can be set to any terminal editor

## Dependencies

- `tmux` (>= 3.0 recommended)
- `fzf`

## Supported Platforms

- **Linux**
- **macOS**

> ❗ **Windows/WSL is not supported.**

## License

This project is licensed under the  Unlicense.

Feel free to fork, customize, or contribute as you see fit! 🎉

