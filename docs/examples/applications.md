# Application Examples

## Minimal Service

```go
type Config struct { 
    Port int `env:"PORT,default=8080"` 
}

var config Config

env.MustDecode(&config)

resp := httpx.Get("https://example.com")

fmt.Println(config.Port, resp.Status().Code)
```

## Processing Pipeline

```go
numbers := slice.Of(1,2,3,4)

filtered := slice.Filter(numbers, func(v int) bool { 
    return v%2==0 
})

res := slice.Map(filtered, func(v int) int { 
    return v*v 
})

fmt.Println(res)
```
