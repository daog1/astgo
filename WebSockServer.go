// WebSockServer.go
package main

import (
	//"fmt"
	"websocket"
	"net/http"
	"time"
	"fmt"
	//"io"
)
/*
func main() {
	fmt.Println("Hello World!")
}*/
func TimeToString16(d time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d", d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute())
}
type JsonEvent struct {
	TypeName string   `json:"type"`
	Name     string   `json:"name"`
	Id       int      `json:"id"`
	Ack      string   `json:"ack"`
	Args     []string `json:"args"`
}
type JsonReply struct {
	TypeName string   `json:"type"`
	Id       int      `json:"id"`
	Args     Message `json:"args"`
}
type AgentClent struct {
	Conn    *websocket.Conn
}
func (p *AgentClent) SendMsg(msg JsonReply){
	websocket.JSON.Send(p.Conn,msg)
}
//type 
func EchoServer(ws *websocket.Conn){ 
//io.Copy(ws, ws) 
	//fmt.Printf("\r\nconn start\r\n")
	GetLogTag("websocket").Info("conn start")
	for{
		je := JsonEvent{}
		err,msg:=websocket.JSON.Receive(ws, &je)
		if(err!=nil){
			GetLogTag("websocket").Err("recv err:%v",err)
			if(msg!=nil){
				msgstr :=string(msg)
				GetLogTag("websocket").Err("recv msg:%v",msgstr)
			}
			break;
		}
		if ws.UserData == nil {
			if je.Name == "agentlogin" {
				ag:=new(AgentClent)
				ag.Conn = ws
				GetAstAmi().addEventHandler(je.Args[0], ag)
				ws.UserData = ag
				//GetAgentManger().LoginAgent(je.Id, je.Args[0], je.Args[1], ws)
			}
		} else {
			//ws.UserData.(*AgentClient).RecvMsg(je.Id, je.Name, je.Args)
		}
		//fmt.Printf("%v",je)
		GetLogTag("websocket").Info("recv %v",je)
	}

	GetLogTag("websocket").Info("conn close")
	//fmt.Printf("\r\nconn close\r\n")
	if(ws.UserData!=nil){
			/*agent:=ws.UserData.(*AgentClient)
			if(agent!=nil){
				agent.Removehandler()
				if(agent.Phone!=nil){
					agent.Phone.UnBindAgent()
					agent.Phone = nil
				}
				GetAgentManger().RemoveAgent(agent)
				ws.UserData =nil
			}*/
		}
	ws.Close()
	//je := JsonEvent{}
	//ws.JSON.Receive(ws, &je)
}
func init() {
	//GetLogTag("AgentClient").
	SetLoggerPriority("websocket",LOG_DEBUG)
}
func StartWebSocketServer(){
	//fmt.Println("StartWebSocketServer ")
	GetLogTag("websocket").Info("StartWebSocketServer")
	urls:="/"+GetConf().Approot+"/echo"
	GetLogTag("websocket").Info("websocke urls %s",urls)
	http.Handle(urls, websocket.Handler(EchoServer));
	//http.Handle("/echo", websocket.Handler(EchoServer))
	//http.Handle("/echo", Handler(echoServer))
	http.Handle("/", http.FileServer(http.Dir("www/")))
	if err := http.ListenAndServe(GetConf().AgentAddr,nil); err != nil {
		//log.Fatal("ListenAndServe:", err)
	}
}
