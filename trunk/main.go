package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var (
	config      *Config
	channelPool *ChannelPool
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var configFile string
	root := filepath.Dir(os.Args[0])
	flag.StringVar(&configFile, "config", filepath.Join(root, "config.json"), "configuration file path, default is ./config.json")
	flag.Parse()

	config = NewConfigFromFile(configFile)
	channelPool = NewChannelPool()

	http.Handle("/qrlogin/public/", http.StripPrefix("/qrlogin/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/qrlogin", indexHandler)
	http.HandleFunc("/qrlogin/get_channel", getChannel)
	http.HandleFunc("/qrlogin/check", checkLogin)
	http.HandleFunc("/qrlogin/login", readyToLogin)
	http.HandleFunc("/qrlogin/confirm", confirmLogin)

	log.Println("ListenAndServe: ", config.Listen)
	err := http.ListenAndServe(config.Listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
