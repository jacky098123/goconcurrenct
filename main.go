package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sourcegraph/conc/pool"
)

func main() {
	p := pool.New().
		WithMaxGoroutines(4).
		WithContext(context.Background()).
		WithFirstError()
	for i := 0; i < 30; i++ {
		i := i
		p.Go(func(ctx context.Context) error {
			if i%3 == 0 {
				return errors.New(fmt.Sprintf("I will cancel all other tasks! %d", i))
			}
			//			<-ctx.Done()
			time.Sleep(1 * time.Second)
			return nil
		})
	}
	err := p.Wait()
	fmt.Println(err)
}
