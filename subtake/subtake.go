package subtake

import (
	"log"
	"sync"
)

type Options struct {
	Domains  string
	Threads  int
	Timeout  int
	Output   string
	Ssl      bool
	All      bool
	Verbose  bool
	Config   string
}

type Subdomain struct {
	Url string
}

/* Start processing from the defined options. */
func Process(o *Options) {
	urls := make(chan *Subdomain, o.Threads*10)
	list, err := open(o.Domains)
	if err != nil {
		log.Fatalln(err)
	}

	wg := new(sync.WaitGroup)

	for i := 0; i < o.Threads; i++ {
		wg.Add(1)
		go func() {
			for url := range urls {
				url.dns(o)
			}

			wg.Done()
		}()
	}

	for i := 0; i < len(list); i++ {
		urls <- &Subdomain{Url: list[i]}
	}

	close(urls)
	wg.Wait()
}
