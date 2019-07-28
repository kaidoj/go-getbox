package pbx

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Getbox struct {
	Config *viper.Viper
}

//NewGetbox starts new getbox instance
func NewGetbox(config *viper.Viper) *Getbox {

	getbox := &Getbox{}
	getbox.Config = config

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
func (gbox *Getbox) Run() error {
	Sync(gbox.Config)
	ms := gbox.Config.GetInt("fetchers_interval") * 1000
	for range time.Tick(time.Duration(ms) * time.Millisecond) {
		res := Sync(gbox.Config)
		if !res {
			return errors.New("Sync failed")
		}
	}

	return nil
}

//RunOnce getbox
func (gbox *Getbox) RunOnce() error {
	if Sync(gbox.Config) {
		os.Exit(1)
	}
	return errors.New("Run once sync failed")
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
