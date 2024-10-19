# TIM - Tmux plugIn Manager

A straightforward plugin manager for tmux.

## Prerequisites

A working go install. Ensure that `$GOPATH/bin` is in your shell's path.

Oh and you probably need tmux.

## Installation

Install tim:

```bash
go install github.com/kjnsn/tim
```

Ask tim to load some plugins when tmux starts:

```bash
# Add anywhere in your ~/.tmux.conf

run "tim load"
```

That's it. Enjoy. I hope tim is a good friend.

If `tim` is not resolving in your path, try `~/go/bin/tim` instead.

## Managing plugins

Adding is as easy as:

```bash
tim add catppuccin/tmux-catppuccin
```

And removing is just as easy:

```bash
tim remove catppuccin/tmux-catppuccin
```

You can see and edit the list of plugins in the config file at
`~/.config/tim/tim.json`.

If you change any versions in the json configuration, just run
`tim add` again to sync.
