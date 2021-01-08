package connectionManager_test

import "testing"

func TestName(t *testing.T) {
	ch := make(chan int, 1024)
	ch <- 1
	close(ch)

	for range ch {
	}
}
