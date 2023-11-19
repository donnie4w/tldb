// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tldb
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file
package tc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	. "github.com/donnie4w/tldb/keystore"
	"golang.org/x/net/websocket"
)

var cmd = &cmdcli{}

type cmdcli struct{}

func (this *cmdcli) Connect() {
	Init()
	addr, _ := StoreAdmin.GetOther("admin")
	auth, _ := StoreAdmin.GetOther("admintauth")
	admintls, _ := StoreAdmin.GetOther("admintls")
	protocol := "ws://"
	if admintls == "1" {
		protocol = "wss://"
	}
	urlstr := protocol + addr + "/local"
	origin := "http://tldb-admin"
	var ws *websocket.Conn
	var err error
	if strings.HasPrefix(urlstr, "wss:") {
		config := &websocket.Config{TlsConfig: &tls.Config{InsecureSkipVerify: true}, Version: websocket.ProtocolVersionHybi13}
		if config.Location, err = url.ParseRequestURI(urlstr); err == nil {
			if config.Origin, err = url.ParseRequestURI(origin); err == nil {
				ws, err = websocket.DialConfig(config)
			}
		}
	} else {
		ws, err = websocket.Dial(urlstr, "", origin)
	}
	if err == nil && ws != nil {
		defer ws.Close()
		if err = websocket.Message.Send(ws, auth); err == nil {
			var byt []byte
			if err = websocket.Message.Receive(ws, &byt); err == nil {
				port := string(byt)
				this.connect2(port)
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (this *cmdcli) connect2(port string) {
	url := "ws://127.0.0.1:" + port
	if ws, err := websocket.Dial(url, "", "http://tldb-cmd"); err == nil {
		defer ws.Close()
		fmt.Println(">>>tldb command [EXIT-> exit+Enter or q! + Enter]")
		fmt.Println(usage)
		fmt.Print(">>>Enter your command:\n")
		stdin(func(txt string) bool {
			cmdtxt := strings.TrimSpace(txt)
			if cmdtxt == "exit" || cmdtxt == "q!" {
				return false
			}
			if err = websocket.Message.Send(ws, cmdtxt); err == nil {
				for {
					var byt []byte
					if err = websocket.Message.Receive(ws, &byt); err == nil {
						if len(byt) == 4 && string(byt) == "exit" {
							return false
						}
						if len(byt) > 3 && strings.HasPrefix(string(byt[:3]), ">>>") {
							s := string(byt)
							fmt.Print(s)
							break
						} else if !strings.HasPrefix(cmdtxt, "export") {
							s := string(byt)
							fmt.Print(s)
						}
					}
					if strings.HasPrefix(cmdtxt, "export") {
						if ss := strings.Split(cmdtxt, " "); len(ss) >= 2 && len(byt) > 0 {
							dir, _ := os.Getwd()
							name := fmt.Sprint(dir, "/", ss[1], ".gz")
							os.WriteFile(name, byt, 0666)
							fmt.Println("export ", name, " successful")
						}
					}
				}
			}
			return true
		})
	}
}

func stdin(f func(txt string) bool) {
	fmt.Print("\n>>>")
	reader := bufio.NewReader(os.Stdin)
	txt, _ := reader.ReadString('\n')
	txtbs := []byte(txt)
	cmdt := strings.TrimSpace(string(txtbs[:len(txtbs)-1]))
	if f(compressSpace(cmdt)) {
		stdin(f)
	}
}

func compressSpace(str string) string {
	reg := regexp.MustCompile("\\s+")
	return strings.TrimSpace(reg.ReplaceAllString(str, " "))
}

var usage = `>>>[addnode xxx:add cluster node][load(or loadforce) xxx xxx:load datafilepath datetime][pwd(or add) xxx: add or alter account]
>>>[del xxx: del account][close: close node]
>>>[other command : pwdcli  addcli  pwdmq  addmq  delcli  delmq  monitor init export]`
