package kraken

import (
	"fmt"
	"golang.org/x/exp/rand"
	"net"
	"time"
)

func getActivePort() (port int, err error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		err = fmt.Errorf("rrror starting listener: %w", err)
		return
	}
	defer func() {
		_ = listener.Close()
	}()
	addr := listener.Addr().(*net.TCPAddr)
	port = addr.Port
	return
}

func randInt(min, max int) int {
	rand.Seed(uint64(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}

func JustSleep() {
	t := time.Duration(randInt(1, 3)) * time.Millisecond
	log.Debugf("just sleep: %s", t)
	time.Sleep(t)
}

func SleepForSlowOperation() {
	t := time.Duration(randInt(5, 15)) * time.Millisecond
	log.Debugf("sleep: %s for slow operation", t)
	time.Sleep(t)
}

func SleepRandSeconds(min, max int) {
	t := time.Duration(randInt(min, max)) * time.Millisecond
	log.Debugf("sleep rand time: %s", t)
	time.Sleep(t)
}

func prepareEventDispatchScript(event string) string {
	return fmt.Sprintf("const event = new MouseEvent('%s', { bubbles: true }); arguments[0].dispatchEvent(event);", event)
}
