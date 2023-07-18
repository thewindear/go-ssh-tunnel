package tunnel

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Forwards struct {
	Items []*Forward
}

func (f Forwards) Run() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		sig := <-sigCh
		log.Printf("exit. sig: %s", sig.String())
		for i := 0; i < len(f.Items); i++ {
			log.Println("send stop.")
			f.Items[i].StopCh <- struct{}{}
		}
	}()
	var wg sync.WaitGroup
	for i := 0; i < len(f.Items); i++ {
		wg.Add(1)
		go func(forward *Forward) {
			if err := forward.Run(); err != nil {
				log.Fatalf("Error: run forward error %s", err)
			}
			wg.Done()
		}(f.Items[i])
	}
	wg.Wait()
}

func BuildWithConfigFile(configPath string) (forwards *Forwards, err error) {
	f, err := os.OpenFile(configPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	var items []*Forward
	err = yaml.NewDecoder(f).Decode(&items)
	forwards = &Forwards{Items: items}
	for i := 0; i < len(forwards.Items); i++ {
		forwards.Items[i].StopCh = make(chan struct{})
	}
	return
}
