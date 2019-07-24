package pbx

import (
	"log"

	"github.com/spf13/viper"
)

//Main getbox
type Getbox struct {
	Config *viper.Viper
}

//Run getbox
func (gbox *Getbox) Run() *Getbox {

	gbox.isNrCoresSet()

	r := &Request{}
	r.Config = gbox.Config

	m := &Move{r}
	f := &Fetch{r, m, 0}

	f.ProjectsToSync()

	return gbox
}

//checks if number of cores is set in config
//will panic if not set
func (gbox *Getbox) isNrCoresSet() bool {
	nrCores := gbox.Config.GetInt("nr_of_cores")
	if nrCores == 0 {
		log.Panicln("Number of cores not set in config")
	}

	return true
}
