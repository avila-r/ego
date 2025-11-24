# dotenv

Enhanced dotenv parsing & loading with support for export, quoting, expansion.

## Load
```go
if err := dotenv.Load(); err != nil { log.Fatal(err) }
```
## Overload
```go
dotenv.Overload(".env.local")
```
## Exec
```go
_ = dotenv.Exec(nil, "go", []string{"run","main.go"}, false)
```
