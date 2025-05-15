# sesh - Terminal Workspace Launcher Tool

## Purpose
`sesh` is a CLI tool that I created for personal use. It uses `fzf` to fuzzy-select a workspace
and launch a tmux session that is ready for a coding session.

This tool is heavily inspired by [ThePrimeagen](https://github.com/ThePrimeagen)'s Dev workflow
and the [tmux-sessionizer](https://github.com/ThePrimeagen/.dotfiles/blob/master/bin/.local/scripts/tmux-sessionizer)
script used to launch tmux sessions.

## Primary Behavior
| Command                | Description                                                                                          |
| ---------------------- | ---------------------------------------------------------------------------------------------------- |
| `sesh`                 | Fuzzy-select a directory from all configured workspace paths (`SESH_PATHS`) and start a tmux session |
| `sesh <path>`          | Launch a tmux session directly from the provided path (absolute or relative)                         |
| `sesh add <path>`      | Add a path to the `SESH_PATHS` list (stored in env var or config file)                               |
| `sesh remove [<path>]` | Remove a path from the `SESH_PATHS` list (with `fzf` if no path provided)                            |
| `sesh list`            | List all currently configured paths                                                                  |
| `sesh config`          | Show current configuration (editor, session behavior, etc.)                                          |

- **🔍 Fuzzy Selection**
  When sesh is run without arguments, it will scan the configured top-level directories and pass them through fzf for selection.

- **⚙️ Environment-Based Configuration**
  Controlled using environment variables
  ```bash
  export SESH_PATHS="~/projects ~/sandbox ~/experiments"
  export SESH_EDITOR="nvim"
  ```
- **🖥 tmux Integration**
  If a tmux session for the selected path exists, it attaches to it. Otherwise, it creates a new session, launching two windows:

    - One window with the configured editor (nvim by default).
    - One window for a general-purpose terminal.

- **📁 Custom Per-Project Session Layouts**
  If a .sesh file is found inside the project directory, it overrides the default session setup.
This file defines a custom layout in terms of windows, splits (vertical/horizontal), and startup commands.

## 🛠 Default Session Behavior
### Without `.sesh` File
The default editor is defined in the `SESH_EDITOR` (default nvim)

1. Open a tmux session named after the directory basename.
2. Create two windows:
    - editor: runs nvim
    - terminal: runs bash or default shell

### 📄 .sesh File Structure (Optional Per Project)
```yaml

windows:
  - name: "editor"
    panes:
      - command: "nvim"
      - command: "bash"
        split: "vertical"  # This pane is split vertically from the previous one

  - name: "dev"
    panes:
      - command: "cd backend && go run main.go"
      - command: "tail -f logs/app.log"
        split: "horizontal"  # This is split horizontally relative to the first
      - command: "bash"
        split: "vertical"    # This is split vertically from the second
```

> This layout is used instead of the default one when a `.sesh` file is present in the selected directory.

## 📦 Configuration
### Environment Variable

| Variable      | Purpose                                                           |
| ------------- | ----------------------------------------------------------------- |
| `SESH_PATHS`  | Space-separated list of paths to search for directories to launch |
| `SESH_EDITOR` | Default editor command to use in the session (default: `nvim`)    |

> You can manage this with your shell profile

## 🔧 Example CLI Usage

```bash
sesh                         # Pick a directory from SESH_PATHS with fzf
sesh ~/projects/myapp        # Launch a session in ~/projects/myapp directly
sesh add ~/projects          # Add a new path to SESH_PATHS
sesh remove ~/sandbox        # Remove a path from SESH_PATHS
sesh list                    # Show all current paths
sesh config                  # Display current configuration
```

## 🔗 Dependencies

- `tmux` (>= 3.0 recommended)
- `fzf`

## 🖥 Supported Platforms

- **Linux**
- **macOS**

> ❗ **Windows/WSL is not supported.**

## 📜 License

This project is licensed under the MIT License.

Feel free to fork, customize, or contribute as you see fit! 🎉

