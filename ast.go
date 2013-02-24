
package main

import (
	"net"
	"bufio"
	"fmt"
	"strings"
	"os"
	"httpjsonrpc"
)
var bufLen = 1024 * 10
var lineLen = 2 // windows : "\r\n" 2 ; linux : "\n" 1
type Message map[string]string

type AstAmi struct{
	Conn net.Conn
	Cbmap map[string] interface{}
}

/*
len:2msg event: {map[ChannelStateDesc: Down Context: from-sip Uniqueid: 13432772
39.9 Privilege: call,all Exten: AccountCode: CallerIDNum: ChannelState: 0 Channe
l: SIP/8002-00000009 Event: Newchannel CallerIDName:]}start recv
len:2start recv
len:2msg event: {map[ChannelStateDesc: Ringing CallerIDName: 8006 CallerIDNum: 8
006 ChannelState: 5 Uniqueid: 1343277239.9 Privilege: call,all Channel: SIP/8002
-00000009 Event: Newstate]}start recv
len:2msg event: {map[ChannelState: 6 CallerIDNum: 8006 Privilege: call,all Uniqu
eid: 1343277239.9 CallerIDName: 8006 ChannelStateDesc: Up Event: Newstate Channe
l: SIP/8002-00000009]}start recv
len:2start recv
len:2start recv
len:2start recv
len:2msg event: {map[Cause: 16 CallerIDName: 8006 Event: Hangup Channel: SIP/800
2-00000009 Privilege: call,all Uniqueid: 1343277239.9 CallerIDNum: 8006 Cause-tx
t: Normal Clearing]}start recv*/
//需要设置exten =>8001,1,Dial(SIP/8001,20,rto)
func (p *AstAmi) recvMsg(reader *bufio.Reader) (Message,bool){
	msg:=Message{}
	println("start recv")
	for {
		line,_:=reader.ReadBytes('\n')
        //println(string(line))
        linestr:=string(line)
        if(len(line)>2){
        	linestr:=strings.TrimSpace(linestr)
        	argv:=strings.Split(linestr, ":")
        	if(len(argv)>=2){
        		//fmt.Printf("msg :%v",argv)
        		key := strings.TrimSpace(argv[1])
        		msg[argv[0]] =key
        	}

        }else if(len(line) == 0){
        	return msg,false
        }else{
        	fmt.Printf("len:%d", len(line))
        	return msg,true
        }
	}
    return msg,true
}
func (p *AstAmi) addEventHandler(channel string,f interface {} ){
	p.Cbmap[channel] = f
}
func (p *AstAmi) connect(ctype string,addr string) (error){
	
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// handle erro
	}
	//conn.SetReadDeadline()
	p.Conn = conn
	return err
}
func (p *AstAmi) startrecv(){
	buf :=bufio.NewReader(p.Conn)
	for {
			msg,rok:=p.recvMsg(buf)
			if(rok){
				_,ok:=msg["Event"]
				if(ok){
					//fmt.Printf("msg event:%s {%v}",msg)
					p.handleEvent(msg)
				}
			}else{
				break
			}
        }
}
func rpc_incoming(callerid string,channel string){
	p2:=[]interface{}{}
	p2=append(p2,callerid)
	rpcaddr:="http://"+GetConf().Rpcaddr+"/server.php"
	res, err:=httpjsonrpc.Call(rpcaddr,"rpcgetCallerInfo",1,p2)
	if err!=nil{
        //log.Fatalf("Err: %v", err)
        fmt.Printf("err:%v\n", err)
	}
	r,_:=res["result"]
	p:=[]interface{}{}
	p=append(p,channel)
	p=append(p,callerid)
	p=append(p,channel)
	p=append(p,"incoming")
	p=append(p,r)
	_, err=httpjsonrpc.Call(rpcaddr,"rpcaddToCallHistory",1,p)
	if err!=nil{
        fmt.Printf("err:%v\n", err)
	}
}
func (p *AstAmi) handleEvent(msg Message){
	channel,ok:=msg["Channel"]
	if ok {
		vl:=strings.Split(channel, "/")
		//fmt.Printf("%v", vl)
		ch:=strings.Split(vl[1], "-")
		ext,b:=GetCtiDb().GetAgent(ch[0])
		if(b){
			fmt.Printf("msg event:%s ",msg)
			etype,_:=msg["Event"]
			if(etype == "Newchannel"){
				uid,_:=msg["Uniqueid"]
				call:=ChannelCall{
					Channel:ch[0],
					Uniqueid:uid,
					Showcallid:false,
				}
				GetCallManger().addCall(&call)
			}else if(etype == "Hangup"){
				uid,_:=msg["Uniqueid"]
				GetCallManger().removeCall(uid)
			}else if(etype == "Newstate"){
				uid,_:=msg["Uniqueid"]
				calleridnum,ok:=msg["CallerIDNum"]
				if(ok){
					if(len(calleridnum)>2){
						call:=GetCallManger().getCall(uid)
						if(call!=nil){
							if call.Channel!=calleridnum {
								call.CallerIDNum = calleridnum
								if(call.Showcallid == false){
									fmt.Printf("incomingcall: {%v}",calleridnum)
									call.Showcallid = true
									GetCtiDb().WriteIncomingcall(uid, calleridnum, call.Channel)
									rpc_incoming(calleridnum,ext.Asterisk_extension)
									//GetCtiDb().WriteIncomingcall(uid, calleridnum, call.Channel)
									//GetCtiDb().addToCallHistory(ext.Userid)
								}
							}

							
						}
					}
				}

			}
		}
		//handle,ok:=p.Cbmap[ch[0]]
		//if(ok){
		//	jr := JsonReply{
		//		TypeName :"Event",
		//		Id:0,
		//		Args:msg,
		//	}
		//	handle.(*AgentClent).SendMsg(jr)
		//	fmt.Printf("msg event: {%v}",handle)
		//}
		
		//if(ch[0] == "1001"){
		//	fmt.Printf("msg event: {%v}",msg)
		//}
		//fmt.Printf("%v", ch[0])
		//SIP/8002-00000005
		//str :=strings.TrimRight(channel, "SIP/")
		//strings.Trim(s, cutset)
		//fmt.Printf("msg channel:%s ",str)
	}

}
func (p *AstAmi) login(usr string,pass string){
	
	//var buf		[512]byte
	p.Conn.Write([]byte("Action: Login\r\nUsername: "+usr+"\r\nSecret: "+pass+"\r\n\r\n"));
	
}
func (p *AstAmi) startcall(agent string,pass string){
	
	//var buf		[512]byte
	p.Conn.Write([]byte("Action: Originate\r\nChannel: Agent/1001"+"\r\nContext: from-sip\r\n"+"Exten: 8002\r\n"+"Priority: 1\r\n\r\n"));
	
}
func CmdMain(c chan int) {
	r := bufio.NewReader(os.Stdin)
	for {
		s, _ := r.ReadString('\n')
		//s=strings.Trim(s," \r\n")
		s = strings.TrimSpace(s)
		//strings.split()
		//if s == "log" {
		//	Log()

		if s == "quit" {
			c <- -1
			break
		}else if s == "startcall" {
			GetAstAmi().startcall("1001", "123456")
		}
		fmt.Printf("%s", s)
	}
}
//结束的时候先收到一行，再收到\r\n结束。还是很好搞定。
var pAstAmi *AstAmi=nil
func GetAstAmi() *AstAmi {
	if pAstAmi == nil {
		pAstAmi = new(AstAmi)
		pAstAmi.Cbmap =map[string] interface{}{}
		//logger = log.New(os.Stdout, "AgentManger:", 0)
	}
	return pAstAmi
}
func main() {
	//ast:=AstAmi{}
	//SIP/8002-00000005
	//teststr :="SIP/8002-00000005"
	//str :=strings.TrimLeft(teststr, "/")
	//vl:=strings.Split(teststr, "/")
	//fmt.Printf("%v", vl)
	//	ch:=strings.Split(vl[1], "-")
	//	fmt.Printf("%v", ch[0])
		//strings.Trim(s, cutset)
	GetAstAmi().connect("tcp", GetConf().IPpbxAddr)
	GetAstAmi().login(GetConf().AmiUser,GetConf().AmiPassword)
	TestCtiDb()
	go StartWebSocketServer()
	go GetAstAmi().startrecv()
	//teststartrecv()
	ci := make(chan int)
	go CmdMain(ci)
	//GenLastTime(2)
	for {
		_ = <-ci
		break
	}
}
/*
package main

import (
        "net"
        "bufio"
        //"os"
)

var bufLen = 1024 * 10
var lineLen = 2 // windows : "\r\n" 2 ; linux : "\n" 1

func main() {
      //  reader := bufio.NewReader(os.Stdin)
      //  print("Input ip : ") ; svrIp, _ := reader.ReadBytes('\n')
       // print("Input userName : ") ; usrName, _ := reader.ReadBytes('\n')
       // print("Input passwd : ") ; pwd, _ := reader.ReadBytes('\n')
        conn,err := net.Dial("tcp", "192.168.0.115:5038")
        defer conn.Close()
        if err != nil {
                println("Error : ",err.Error()) //In go 1 , use err.Error() ,not err.String()
        }
        conn.Write([]byte("Action: login\r\nUserName: "+
        "mark"+"\r\nSecret: "+
        "123456"+"\r\n\r\n"))
        reader:=bufio.NewReader(conn)
        for {
               // p := make([]byte,bufLen)
                //sz, _ := bufio.NewReader(conn).Read(p)
        	line, _ := reader.ReadString('\n')
        	println(string(line))
                //println(string(p[0:sz]))
        }
}//*/
