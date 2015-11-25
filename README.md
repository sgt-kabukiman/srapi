speedrun.com API Client for Go
==============================
[![GoDoc](https://godoc.org/github.com/sgt-kabukiman/srapi?status.svg)](https://godoc.org/github.com/sgt-kabukiman/srapi)

This Go package implements a client for the
[speedrun.com API](https://github.com/speedruncom/api). It's not 100% complete
and a relatively direct mapping of API structures to Go ``struct`` values.

Installation
------------

```
go get github.com/sgt-kabukiman/srapi
```

Usage
-----

```go
package main

import (
	"fmt"

	"github.com/sgt-kabukiman/srapi"
)

func main() {
	// optional sorting
	sort := &srapi.Sorting{"name", srapi.Descending}

	// optional pagination
	cursor := &srapi.Cursor{2, 2} // offset, max

	// optional embeds
	embeds := srapi.NoEmbeds

	regions, err := srapi.Regions(sort, cursor, embeds)
	if err != nil {
		panic(err) // err is an srapi.*ApiError struct, containing more information
	}

	fmt.Printf("regions = %+v\n", regions)
}
```

License
-------

This code is licensed under the MIT license.
