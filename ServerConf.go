// ServerConf.go
package main

import (
	//"fmt"
	"encoding/xml"
	"os"
)


type ServerConf struct {
	
	XMLName   xml.Name `xml:"server"`
	AgentAddr string   `xml:"agentAddr"`
	IPpbxAddr string   `xml:"ippbxAddr"`
	License   string   `xml:"license"`
	Approot   string   `xml:"approot"`
	AmiUser   string   `xml:"amiuser"`
	AmiPassword string `xml:"amipasswd"`
	DBAddr    string   `xml:"dbaddr"`
	DBName    string   `xml:"dbname"`
	DBUser    string   `xml:"dbuser"`
	DBPasswd  string   `xml:"dbpasswd"`
	Rpcaddr   string 	`xml:"rpcaddr"`
}

var conf *ServerConf

func LoadConf() {
	conf = &ServerConf{
		IPpbxAddr: "None",
		AgentAddr: "None",
		License:   "None",
		Approot:   "None",
		AmiUser:   "None",
		AmiPassword: "None",
		DBAddr:	"None",
		DBName:	"None",
		DBUser:	"None",
		DBPasswd:	"None",
		Rpcaddr:"None",
	}
	f, err := os.Open("serverconf.xml")
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()
	d := xml.NewDecoder(f)
	//var m interface {}
	err = d.Decode(&conf)
}
func GetConf() *ServerConf {
	if conf == nil {
		LoadConf()
	}
	return conf
}
