# resultGroup

[![CI](https://github.com/orlovm/resultgroup/actions/workflows/go.yml/badge.svg)](https://github.com/orlovm/resultgroup/actions/workflows/go.yml)
[![codebeat badge](https://codebeat.co/badges/712c7ef7-4ac9-4df9-96b0-bd5c2d547e30)](https://codebeat.co/projects/github-com-orlovm-resultgroup-main)
[![GitHub release](https://img.shields.io/github/release/orlovm/resultgroup.svg?label=version)](https://github.com/orlovm/resultgroup/releases/latest)  

resultGroup is a simple and flexible Go library for managing the results and errors of concurrent tasks. It is inspired by the beloved HashiCorp's [go-multierror/Group](https://github.com/hashicorp/go-multierror/blob/master/group.go).

## Motivation

The need to aggregate results from multiple sources is a common requirement in many applications. While libraries such as [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) and [go-multierror](https://github.com/hashicorp/go-multierror/blob/master/group.go) can be used for this purpose, they often require additional boilerplate code and the use of mutexes or channels for synchronization. Result Group aims to streamline this process by offering a generic solution that minimizes boilerplate and simplifies the management of concurrent tasks.


## Usage 
  
> **Note**
> ResultGroup works with Go 1.20 and is compatible with Go 1.20 wrapped errors.    
  
To use Result Group, follow these steps:

1. Import the package:

```go
import "github.com/orlovm/resultgroup"
```

2. Create a new Result Group:

```go
group := resultgroup.Group[ResultType]{}
```

Replace `ResultType` with the type of the results you expect to collect.

3. Alternatively, create a new Result Group with an error threshold:

```go
ctx := context.Background()
threshold := 1
group, ctx := resultgroup.WithErrorsThreshold[ResultType](ctx, threshold)
```

4. Run concurrent tasks using the `Go` method:

```go
group.Go(func() ([]ResultType, error) {
    // Your concurrent task logic here
})
```

5. Wait for all tasks to complete and collect the results:

```go
results, err := group.Wait()
```

`err` could be used as usual go 1.20 wrapped error, or be easily unwrapped with `Unwrap() []error`

Here's a complete example that demonstrates how to use Result Group to fetch data from multiple sources concurrently:

```go
package main

import (
	"context"
	"fmt"

	"github.com/orlovm/resultgroup"
)

type Data struct {
	Source string
	Value  int
}

func fetchData(source string) ([]Data, error) {
	// Simulate fetching data from the source
	data := []Data{
		{Source: source, Value: 1},
		{Source: source, Value: 2},
	}

	return data, nil
}

func main() {
	sources := []string{"source1", "source2", "source3"}

	ctx := context.Background()
	group, _ := resultgroup.WithErrorsThreshold[Data](ctx, 1)

	for _, source := range sources {
		source := source
		group.Go(func() ([]Data, error) {
			return fetchData(source)
		})
	}

	results, err := group.Wait()
	if err != nil {
		fmt.Println("Error:", err)
                fmt.Println("Wrapped errors", err.Unwrap())
	}

	for _, result := range results {
		fmt.Printf("Source: %s, Value: %d\n", result.Source, result.Value)
	}
}
```  

See [tests](https://github.com/orlovm/resultgroup/blob/main/group_test.go) for more examples.
## License

[MIT](LICENSE)
