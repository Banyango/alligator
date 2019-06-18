package main

import (
	"fmt"
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
	// todo replace with env.
	fmt.Println(banner)

	http.ListenAndServe(":8080", nil)
}
