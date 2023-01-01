# wpool_gdl
worker pool implementation in golang

### available options
1. [ttl ] time to leave (=timeout)
2. [mrtc] max retrying count

### usage:
run `go get github.com/siamak4mo/wpool_gdl` and `go install github.com/siamak4mo/wpool_gdl` to get package.
like `example.go` use this import:
```
import (
	...

	wpool "github.com/siamak4mo/wpool_gdl"
)
```
