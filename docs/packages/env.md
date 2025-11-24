# env

Structured decoding of environment variables into tagged structs.

## Decode

```go
type Config struct {
    Port int    `env:"PORT,default=8080"`
    Mode string `env:"MODE,required"`
}

var c Config

if err := env.Decode(&c); err != nil { 
    panic(err) 
}
```

## Export Metadata

```go
info, _ := env.Export(&c)

for _, f := range info { 
    fmt.Println(f.EnvVar, f.Value) 
}
```
