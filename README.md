# flenv: FLag + ENV for Go

```
type Config struct {
  Host string `env:"HOST" default:"localhost" flag:"-h,--host" help:"Host to listen on"`
  Port string `env:"PORT" default:"80" flag:"-p,--port" help:"Port to listen on"`
}

func main() {
  var config Config
  flagSet, err := flenv.DecodeArgs(&config)
}
```

