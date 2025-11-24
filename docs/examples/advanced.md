# Advanced Examples

## Promise Composition Chain
```go
final := promise.Compose(
    promise.Of(func() (int,error){ return 2,nil }),
    func(v int) *promise.Promise[int] { return promise.Of(func() (int,error){ return v*10,nil }) },
).Then(func(v int) int { return v+1 })
fmt.Println(final.Join()) // 21
```
## Env Decode + Use
```go
type AppCfg struct { Port int `env:"PORT,default=8080"` }
var cfg AppCfg
env.MustDecode(&cfg)
fmt.Println(cfg.Port)
```
