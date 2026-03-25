package logo

import (
	"fmt"

	"github.com/aceld/zinx/v3/zconf"
)

var zinxLogo = `                                        
              ██                        
              ▀▀                        
 ████████   ████     ██▄████▄  ▀██  ██▀ 
     ▄█▀      ██     ██▀   ██    ████   
   ▄█▀        ██     ██    ██    ▄██▄   
 ▄██▄▄▄▄▄  ▄▄▄██▄▄▄  ██    ██   ▄█▀▀█▄  
 ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀  ▀▀    ▀▀  ▀▀▀  ▀▀▀ 
                                        `
var topLine = `┌──────────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└──────────────────────────────────────────────────────┘`

func PrintLogo() {
	fmt.Println(zinxLogo)
	fmt.Println(topLine)
	fmt.Printf("%s [Github] https://github.com/aceld                    %s\n", borderLine, borderLine)
	fmt.Printf("%s [tutorial] https://www.yuque.com/aceld/npyr8s/bgftov %s\n", borderLine, borderLine)
	fmt.Printf("%s [document] https://www.yuque.com/aceld/tsgooa        %s\n", borderLine, borderLine)
	fmt.Println(bottomLine)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		zconf.GlobalObject.Version,
		zconf.GlobalObject.MaxConn,
		zconf.GlobalObject.MaxPacketSize)
}
