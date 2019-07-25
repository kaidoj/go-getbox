package pbx

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Getbox struct {
	Config *viper.Viper
	fetch  Fetcher
}

//NewGetbox starts new getbox instance
func NewGetbox(config *viper.Viper) *Getbox {

	getbox := &Getbox{}
	getbox.Config = config
	r := &Request{}
	r.Config = config
	m := &Move{r}
	getbox.fetch = NewFetcher(r, m)

	if len(os.Args) <= 1 {
		getbox.Run()
		return getbox
	}

	if os.Args[1] == "run_once" {
		getbox.RunOnce()
	} else {
		getbox.Run()
	}

	return getbox
}

//Run getbox forever
func (gbox *Getbox) Run() {
	gbox.fetch.ProjectsToSync()
	ms := gbox.Config.GetInt("fetchers_interval") * 1000
	for range time.Tick(time.Duration(ms) * time.Millisecond) {
		gbox.fetch.ProjectsToSync()
	}
}

//RunOnce getbox
func (gbox *Getbox) RunOnce() *Getbox {
	if gbox.fetch.ProjectsToSync() {
		os.Exit(1)
	}
	return gbox
}

//IsNrCoresSet checks if number of cores is set in config
//will panic if not set
func IsNrCoresSet(config *viper.Viper) bool {
	nrCores := config.GetInt("nr_of_cores")
	if nrCores == 0 {
		log.Panicln("Number of cores not set in config")
	}

	return true
}
