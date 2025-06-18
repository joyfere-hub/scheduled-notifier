# Scheduled notifier

A simple scheduled notifier tool, scheduled call api, then notifier.

## Usage

Use makefile build.

```shell
make all
```

### MacOS 

On `MacOS`, `terminal-notifier` must be installed to receive system notifications.

```shell
brew install terminal-notifier
```

Package dmg (Don't forge install `create-dmg`).

```shell
make package-mac
```

Run an executable file.

```shell
build/mac/ScheduledNotifier.app/Contents/MacOS/ScheduledNotifier
```

Open an app.

```shell
open build/mac/ScheduledNotifier.app
```

### Windows

Use [go-toast](https://git.sr.ht/~jackmordaunt/go-toast) to receive system notifications, See [adaptive-interactive-toasts](https://learn.microsoft.com/en-us/windows/apps/design/shell/tiles-and-notifications/adaptive-interactive-toasts)。

```shell
build/windows/ScheduledNotifier.exe
```

### Config

Config order:
1. `~/.local/scheduled-notifier.yaml`
2. `env.SCHEDULED_NOTIFIER_CONFIG_FILE_PATH` 
3. `ScheduledNotifier -config scheduled-notifier.yaml`

### Build

Base template, more template see `internal/job/*_template.md`:

```yaml
jobs:
  - type: rebuild_work_task     # api client type, see internal/job
    name: test                  # job name
    interval: '@every 5m'       # cron
    token: token                # The specific authentication is determined by the JobClient implementation.
```

## TODO
- [ ] Jenkins
- [ ] Gitlab
- [ ] more ...

## Other

I'm a Rookie. 