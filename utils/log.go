package utils

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

var (
	dir      = "logs/"
	filename = "getbox.log"
)

//LogOutput logs output to file
//Filename format yyyy-mm-dd-filename
func LogOutput(config *viper.Viper) *os.File {

	logsDir := config.GetString("logs_dir")
	if logsDir != "" {
		dir = logsDir
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModeDir)
	}

	file, err := os.OpenFile(dir+DateFileName(filename), os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

//DateFileName formats filename to yyyy-mm-dd-filename.log
func DateFileName(filename string) string {
	t := time.Now().Local()
	s := t.Format("2006-01-02")
	return s + "-" + filename
}
