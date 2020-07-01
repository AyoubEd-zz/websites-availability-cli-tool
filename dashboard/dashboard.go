package dashboard

import (
	"fmt"
	"time"

	"github.com/ayoubed/datadog-home-project/statsagent"
)

func Run(urls []string, interval, span int) {
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case t := <-ticker.C:
				fmt.Printf("[last %vs stats]---------------------%v-------------------\n", interval, t)
				res := statsagent.GetStats(urls, span)
				for k, v := range res {
					fmt.Println("> ", k)
					fmt.Printf("%+v\n", v)
				}
			}
		}
	}()
}
