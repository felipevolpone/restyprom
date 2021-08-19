# restyprom

> a simple wrapper around [resty](https://github.com/go-resty/resty) to
report HTTP calls metrics to prometheus

If you're using [resty](https://github.com/go-resty/resty) and want to
have metrics of your HTTP calls, `restyprom` is here for you.

For now, these are the metrics available:
- Response time of called URLs
- Total of calls per URL and status code
- Total of success calls per URL
- Total of failure calls per URL

## Install

```shell
go get github.com/felipevolpone/restyprom
```

## Getting Started

```go
client := resty.New()
client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
    restyprom.Collect(r)
    return nil
})
restyprom.Init() // to register your metrics

resp, err := client.R().Get("https://httpbin.org/get")
```

If you're creating simple `resty` clients, you can use the `NewBasicClient`
func to wrap that code and use it in a simpler way, just  like that:

```go
client := restyprom.NewBasicClient()
restyprom.Init() // to register your metrics

resp, err := client.R().Get("https://httpbin.org/get")
```

## Example with Gin

```go
client := restyprom.NewBasicClient()
restyprom.Init()

r := gin.Default()
r.GET("/ping", func(c *gin.Context) {
    resp, _ := client.R().Get("https://httpbin.org/get")
    c.JSON(200, gin.H{
        "message": resp.Body(),
    })
})

r.GET("/metrics", gin.WrapH(promhttp.Handler()))
r.Run()
```

## Details

If you're not using `prometheus.DefaulRegister` you can use it this way:

```go
restyprom.InitWithRegister(yourRegister)
```