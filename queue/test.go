package queue

import "fmt"

type TestQueue struct{}

var Test TestQueue

func (t *TestQueue) Exec(d string) {
	fmt.Println(d)
}
