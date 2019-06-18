package main

import (
	"fmt"
	"github.com/Banyango/Alligator/config"
	"github.com/Banyango/Alligator/reverseProxy"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
)

const banner = `
   _____  .__  .__  .__              __                
  /  _  \ |  | |  | |__| _________ _/  |_  ___________ 
 /  /_\  \|  | |  | |  |/ ___\__  \\   __\/  _ \_  __ \
/    |    \  |_|  |_|  / /_/  > __ \|  | (  <_> )  | \/
\____|__  /____/____/__\___  (____  /__|  \____/|__|   
        \/            /_____/     \/                   

Proxy Starting up on port 8080
`

func main() {
	fmt.Println(banner)

	tomlBytes, err := ioutil.ReadFile("./alligator.toml")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := config.New(string(tomlBytes))
	if err != nil {
		log.Fatal(err)
	}

	handler := reverseProxy.New(*conf)

	if err := http.ListenAndServe(":8080", handler.Build()); err != nil {
		log.Fatal(err)
	}
}
