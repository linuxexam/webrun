package util

import (
	"fmt"
	"sync"
	"time"
)

var uuid_lock sync.Mutex
var uuid_counter int64

func GenerateUUID() string {
	uuid_lock.Lock()
	defer uuid_lock.Unlock()

	uuid_counter++
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), uuid_counter)
}
