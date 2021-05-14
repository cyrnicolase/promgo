# Promgo

## Install
```shell
$ go get -u github.com/cyrnicolase/promgo
```

## Usage

```go

var (
    // APIRequestTotal 接口请求量
    APIRequestTotal promgo.Counter
)

func init() {
    rdb := redis.NewClient(&redis.Options{
        Addr: `:6379`,
    })

    APIRequestTotal = promgo.NewCounter(rdb, promgo.CounterOptions{
        Name: `api_request_total`,
        Help: `api request counter`,
        Labels: []string{`method`, `endpoint`},
    })

    promgo.GetDefaultRegistry().MustRegister(APIRequestTotal)
}

func main() {
    http.HandleFunc(`/index`, func(rw http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), time.Second)
        defer cancel()

        APIRequestTotal.Inc(ctx, promgo.ConstLables{
            `method`:   r.Method,
            `endpoint`: r.URL.Path,
        }) 

        fmt.Fprint(rw, `hello`)
    })

    http.HandleFunc(`/metrics`, promgo.Render())
    http.ListenAndServe(`:1234`, nil)
}

```