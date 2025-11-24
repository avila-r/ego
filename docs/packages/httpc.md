# httpx (httpc)

Lightweight HTTP helpers: `Get`, typed `Status`, `Body`, basic request builder.

## Simple GET
```go
resp := httpx.Get("https://example.com")
if resp.HasFailed() { log.Fatal(resp.Error()) }
fmt.Println(resp.Status().Label)
```
## Query Build
```go
url, _ := httpx.Url("https://api").WithQueryParam("page","1").Build()
```
