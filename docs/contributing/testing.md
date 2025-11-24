# Testing

## Run

```bash
go test ./...
```

## Focused

```bash
go test ./promise -run TestPromise
```

## Principles

- Unit first, no global state mutation.
- Async tests: use timeouts conservatively.
