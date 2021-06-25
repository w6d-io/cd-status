# CI STATUS

This tool tekton and pods events and notify all webhooks recorded


To trigger the watching a payload has to bu post to `/watch/play` URI

## Installation

This chart [ci-status](https://github.com/w6d-io/charts/tree/main/charts/ci-status) can be used to install this tool

## Configuration

### Main

| Parameter  | Description                                                       | Default |
|------------|-------------------------------------------------------------------|---------|
| `listen`   | address and port to bind the application                          | `:8080` |
| `timeout`  | Time in minute of non-activity when the application stop watching | `60`    |
| `hooks`    | list of url and scope                                             | `[]`    |

### Hooks

| Parameter  | Description                               | Default  |
|------------|-------------------------------------------|----------|
| `url`      | full url to send notification             | `nil`    |
| `scope`    | scope define when and what events to send | `nil`    |

#### url

schemes supported

- http
- https
- kafka

#### scopes

List of scopes supported.

- timeout
- update
- and

scope also support regex

## run

The configuration file is mandatory

```shell
$ > /ci-status --config config.yaml --log-format text --log-level debug
```
