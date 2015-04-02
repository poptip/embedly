# extract

A go client for the embedly Extract API.

### Usage

```go
package main

import (
  "github.com/poptip/embedly"
)

func main() {
  c := embedly.NewClient(key)
  urls := []string{}
  options = embedly.Options{}
  responses, err := embedly.Extract(urls, options)
}
```

### Install

    go get github.com/poptip/embedly

### License

MIT
