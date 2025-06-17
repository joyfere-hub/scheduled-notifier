# Scheduled notifier

一个简单的定时通知工具，通过定时调用 API 的方式，拉取通知到 MacOS 与 Window 的系统通知。

## Usage

### MacOS 

```shell
brew install terminal-notifier
```

默认使用`~/.local/scheduled-notifier.yaml`作为配置文件，当然也可以配置环境变量`SCHEDULED_NOTIFIER_CONFIG_FILE_PATH`覆盖默认路径。
同时支持在使用时添加`-config scheduled-notifier.yaml`，配置格式如下：

```yaml
jobs:
  - type: rebuild_work_task # 客户端类型
    name: test # 配置名称
    interval: '@every 5m' # cron
    token: token # token 认证方式，具体认证由 JobClient 实现决定。
```