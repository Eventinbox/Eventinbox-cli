# EventInbox CLI

Command-line tool for [EventInbox](https://eventinbox.pro) — webhook delivery
infrastructure for developers. Send test events, inspect deliveries, tail
delivery logs, and replay deliveries from your terminal.

## Install

```sh
go install github.com/Eventinbox/Eventinbox-cli@latest
```

This installs the `Eventinbox-cli` binary into `$(go env GOPATH)/bin`. Make
sure that directory is on your `PATH`.

Prebuilt binaries for macOS, Linux, and Windows (amd64 + arm64) are also
attached to each [GitHub release](https://github.com/Eventinbox/Eventinbox-cli/releases).

## Authentication

The `deliveries`, `logs`, and `replay` commands talk to the EventInbox API and
need an API key and workspace ID. Create an API key in the dashboard under
**Settings → API keys**, then export both values:

```sh
export EI_API_KEY=your_api_key
export EI_WORKSPACE_ID=your_workspace_id
```

You can also pass them per command with `--api-key` and `--workspace`. The API
base URL defaults to `https://api.eventinbox.pro` and can be overridden with
`--api-url`.

## Commands

### `send` — send a test event to an endpoint

```sh
eventinbox send payment.created --tenant acme --endpoint payments --payload '{"amount":5400}'
```

### `deliveries` — list and inspect deliveries

```sh
eventinbox deliveries --status delivered
```

### `logs` — tail live delivery logs

```sh
eventinbox logs --tail
```

### `replay` — replay a delivery by ID

```sh
eventinbox replay b4901ce3-d019-4c5c-96a2-d7990b045b7b
```

Run `eventinbox --help` or `eventinbox <command> --help` for the full list of
flags.

## Documentation

Full documentation and API reference: <https://eventinbox.pro/docs>

## License

[MIT](./LICENSE)
