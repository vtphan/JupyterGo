package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type Configuration struct {
	IP       string
	Port     int
	Database string
	Address  string
}

func informIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			ip4 := ipnet.IP.To4()
			if ip4 != nil {
				switch {
				case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
				case ip4[0] == 192 && ip4[1] == 168:
				default:
					return ip4.String()
				}
			}
		}
	}
	return ""
}

func init_config(filename string) *Configuration {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	config := &Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	if config.IP == "" {
		config.IP = informIPAddress()
	}
	config.Address = fmt.Sprintf("%s:%d", config.IP, config.Port)
	return config
}

func echo(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	fmt.Fprintf(w, content)
}

func main() {
	config_file := ""
	flag.StringVar(&config_file, "c", config_file, "json-formatted configuration file.")
	flag.Parse()
	if config_file == "" {
		flag.Usage()
		os.Exit(1)
	}

	Config := init_config(config_file)

	file, _ := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(file)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })
	http.HandleFunc("/echo", echo)

	log.Printf("Serving at: %s\n", Config.Address)
	fmt.Printf("Serving at: %s\n", Config.Address)
	err := http.ListenAndServe(Config.Address, nil)
	if err != nil {
		log.Fatal("Unable to serve carpo server at " + Config.Address)
	}
}