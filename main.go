package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/vadviktor/telegram-msg"
)

const banned = `\/:*?"<>|%`

func init() {
	flag.Usage = func() {
		u := `Searches for path names that may break/not appear in Samba shares.

Create a config file named find-samba-incompatible-path.json by filling in what is defined in its sample file.
`
		fmt.Fprint(os.Stderr, u)
	}
	flag.Parse()

	viper.SetConfigName("find-samba-incompatible-path")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
}

func main() {
	telegram := &telegram_msg.Telegram{}
	telegram.Create(viper.GetString("botToken"),
		viper.GetInt("targetId"))
	//telegram.SendSilent("Begin to search.")

	searchPaths := viper.GetStringSlice("searchPaths")
	var pathsFound int
	for _, p := range searchPaths {
		err := filepath.Walk(p,
			func(path string, info os.FileInfo, err error) error {
				if strings.ContainsAny(filepath.Base(path), banned) {
					if info.IsDir() {
						pathsFound += 1
						fmt.Println("[" + path + "]")
					} else {
						pathsFound += 1
						fmt.Println(path)
					}
				}
				return err
			})

		if err != nil {
			fmt.Println(err)
			telegram.Send(fmt.Sprintf("Error while walking path: %s\n",
				err.Error()))
		}
	}

	if pathsFound > 0 {
		telegram.SendMD(fmt.Sprintf("*Found %d paths*", pathsFound))
	}

	//telegram.SendSilent("Finished.")
}
