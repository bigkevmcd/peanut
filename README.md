# peanut

Just a k8s resource parser.

This is a PoC for parsing resources.

## Building

```shell
$ go build ./cmd/peanut
```

## Running

```shell
$ peanut --kustomization-path ./path/to/kustomization.yaml
application: go-demo
name         namespace  replicas
go-demo-http production 1
redis        production 1
```

## Testing

```shell
$ go test ./...
```
