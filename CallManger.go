package main
import (
	//"net"
	//"bufio"
	//"fmt"
	//"strings"
	//"os"
)


type ChannelCall struct{
	Channel string
	CallerIDNum string
	CallerIDName string
	ChannelState int
	Uniqueid string
	Showcallid bool 
	//CallerIDNum
}
type CallManger struct {
	//Clients list.List
	Calls map[string] *ChannelCall
}

var pCallManger *CallManger
func GetCallManger() *CallManger {
	if pCallManger == nil {
		pCallManger = new(CallManger)
		pCallManger.Calls =map[string] *ChannelCall{}
		//logger = log.New(os.Stdout, "AgentManger:", 0)
	}
	return pCallManger
}
func (p *CallManger) addCall(call* ChannelCall) bool{
	_,ok:=p.Calls[call.Uniqueid]
	if(ok){
		return false
	}
	p.Calls[call.Uniqueid] = call
	return true
}
func (p *CallManger) removeCall(uid string){
	//p.Calls[uid]=nil,false
	delete(p.Calls,uid)
}
func (p *CallManger) getCall(uid string)* ChannelCall{
	call,ok:=p.Calls[uid]
	if(ok){
		return call
	}
	return nil
}