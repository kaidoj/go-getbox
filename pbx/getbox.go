package pbx

import (
	"fmt"
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

	r := &Request{gbox.Config, ""}
	var f Fetch
	f = &Fetcher{r}
	project := f.ProjectToSync()

	fmt.Printf("Name is %s; Id is %s", project.Name, project.Id)
	fmt.Printf("\r\nStatus is %d", project.Render.Status)
	fmt.Printf("\r\nURL is %s", project.Render.URL)

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
