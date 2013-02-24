// CTIDB.go
package main

import (
	"fmt"
	"mysql"
	"time"
	//"strconv"
)

/*
func main() {
	fmt.Println("Hello World!")
}
*/
//const AGENT_GROUPS_TABLE = "agent_groups_t"
//const AGENT_GROUP_TABLE = "agent_group_t"

type CtiDb struct {
	db	*mysql.Client
	QueryS	[]interface{}
	Ch	chan int
}

var pCtiDb *CtiDb

func GetCtiDb() *CtiDb {
	if pCtiDb == nil {
		pCtiDb = new(CtiDb)
		b := pCtiDb.InitDb()
		if b {
			fmt.Println("cti db connectok ")
		} else {
			fmt.Println("cti db connect error ")
		}
		fmt.Println("GetCtiDb new ")
	}
	return pCtiDb
}
func (p *CtiDb) InitDb() bool {
	//GetConf().DBAddr
	p.Ch = make(chan int)
	db, err := mysql.DialTCP(GetConf().DBAddr, GetConf().DBUser, GetConf().DBPasswd, GetConf().DBName)
	if err != nil {
		return false
	}
	p.db = db
	go p.DBQueryRun()
	return true
}
func (p *CtiDb) CloseDB() {
	if p.db != nil {
		p.db.Close()
	}
}
//date('Y-m-d H:i:s');
func TimeToString(d time.Time) string {
	//if(d == nil)
	if d.IsZero() {
		return "null"
	}
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second())
}
func TimeToStringSQL(d time.Time) string {
	//if(d == nil)
	if d.IsZero() {
		return "null"
	}
	return fmt.Sprintf("\"%04d-%02d-%02d %02d:%02d:%02d\"", d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second())
}

func (p *CtiDb) DBQueryRun() {
	for {
		_ = <-p.Ch
		for len(p.QueryS) > 0 {
			//fmt.Printf("Query run\r\n")
			switch t := p.QueryS[0].(type) {
			case string:
				//v:= p.QueryS[0].(string)
				err := p.db.Query(t)
				if err != nil {
					//fmt.Printf("Query error %v by:%s\n", err,v)
					GetLogTag("ctidb").Err("error:%v by %s", err, t)
				}
			case DBQuery:
				//fmt.Printf("Query run----\r\n")
				//v:=p.QueryS[0].(DBQuery)
				t.Run()
			}
			p.QueryS = p.QueryS[1:]
		}

	}
}


type DBQuery struct {
	Run func()
}
type AstExten struct{
	Userid int64
	Asterisk_extension string
	Use_asterisk string
}
func (p *CtiDb)GetAgent(agent string)(exten AstExten,b bool){
	//group := []string{}
	b = false
	if p.db == nil {
		return 
	}
	querystr := "SELECT userid, asterisk_extension,use_asterisk FROM vtiger_asteriskextensions"
	querystr=querystr+" where asterisk_extension="+agent;
	//fmt.Printf("Query str %v", querystr)
	//"SELECT * FROM " + AGENT_GROUPS_TABLE
	d := DBQuery{}
	ch := make(chan int)
	d.Run = func() {
		fmt.Printf("start query")
		err := p.db.Query(querystr)
		if err != nil {
			fmt.Printf("Query error %v", err)
			ch <- 0
			return
		}
		res, err := p.db.StoreResult()
		if res.RowCount() > 0 {
			r := res.FetchRows()
			exten.Userid = r[0][0].(int64)
			exten.Asterisk_extension = r[0][1].(string)
			exten.Use_asterisk = r[0][2].(string)
			err = res.Free()
			if err != nil {
				fmt.Printf("res free error %v", err)
				ch <- 0
				return
			}
			ch <- 1
			return
		}else {
			err = res.Free()
			if err != nil {
				fmt.Printf("res free error %v", err)
				ch <- 0
				return
			}
			ch <- 0
		}
	}
	p.QueryS = append(p.QueryS, d)
	p.Ch <- 1
	r := <-ch
	if r == 1 {
		b = true
		return 
	}
	return 
}
func (p *CtiDb) WriteIncomingcall(uid string,from string,to string){
	//b = false
	if p.db == nil {
		return 
	}
	querystr := "INSERT INTO vtiger_asteriskincomingcalls (refuid, from_number, from_name, to_number, callertype, flag, timer) VALUES(\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%d\",\"%d\")"
	querystr =fmt.Sprintf(querystr,uid,from, from,to,"SIP",0,time.Now().Unix())
	fmt.Printf("Query str %v", querystr)
	p.QueryS = append(p.QueryS, querystr)
	p.Ch <- 1
} 
func (p *CtiDb) getUniqueID(table string) int64 {
	querystr2:="update "+table+"_seq set id = id + 1"
	d2 := DBQuery{}
	ch := make(chan int)
	id := int64(0)
	d2.Run = func() {
		err:= p.db.Query(querystr2)
		if err != nil {
			fmt.Printf("Query error %v", err)
			ch <- 0
			return
		}

		querystr3:="SELECT id FROM "+table+"_seq"
		err = p.db.Query(querystr3)
		if err != nil {
			fmt.Printf("Query error %v", err)
			ch <- 0
			return
		}else{
			res, _ := p.db.StoreResult()
			if res.RowCount() > 0 {
				r := res.FetchRows()
				id= r[0][0].(int64)
				fmt.Printf("Query id %v", id)
			}
			err = res.Free()
			if err != nil {
				fmt.Printf("res free error %v", err)
				ch <- 0
				return
			}
		}
		ch <- 1
		return 
	}
	p.QueryS = append(p.QueryS, d2)
	p.Ch <- 1
	r:= <-ch
	if r ==1 {
		fmt.Printf("Query id %v", id)
		return id
	}
	return 0
	//r3 := <-ch
	//	$ok = $this->Execute("update $seq with (tablock,holdlock) set id = id + 1");
}
func (p *CtiDb)addToCallHistory(userid int64){
	crmID:=p.getUniqueID("vtiger_crmentity")
	fmt.Printf("Query crmID %v", crmID)
	timeOfCall:=time.Now()
	querystr:=fmt.Sprintf("insert into vtiger_crmentity values (\"%d\",\"%d\",\"%d\",\"%d\",\"%s\",\"%s\",\"%s\",\"%s\",%s,%s,\"%d\",\"%d\",\"%d\")", 
		crmID,userid,userid,0,"PBXManager", "",TimeToString(timeOfCall),TimeToString(timeOfCall),"null","null",0,1,0)
	d2 := DBQuery{}
	ch := make(chan int)
	//id := int64(0)
	fmt.Printf("Query str %v\n", querystr)
	d2.Run = func() {
		err:= p.db.Query(querystr)
		if err != nil {
			fmt.Printf("Query error %v", err)
			ch <- 0
			return
		}
		ch <- 1
		return
	}
	p.QueryS = append(p.QueryS, d2)
	p.Ch <- 1
	_ = <-ch//*/
}
func TestCtiDb() {
	GetCtiDb().GetAgent("1001")
	//GetCtiDb().getUniqueID("vtiger_crmentity")
	//GetCtiDb().addToCallHistory(5)
	//GetCtiDb().WriteIncomingcall("1223233434", "8002", "1001")
	//g, b := GetCtiDb().GetAgentGroups()
	//if b {
	//	fmt.Printf("GetAgentGroups ok%v", g)
	//}
	//groups, b := GetCtiDb().GetAgentByGroups("3001")
	//if b {
	//	fmt.Printf("GetAgentGroups ok%v", groups)

	//}
	//cs := GetCtiDb().GetIVRChannels()
	//fmt.Printf("%v", cs)
	/*b := GetCtiDb().AuthNameAndPasswd("3001", "123456")
	if b {
		fmt.Printf("AuthNameAndPasswd ok")
	}
	g, b := GetCtiDb().GetAgentGroups()
	if b {
		fmt.Printf("GetAgentGroups ok%v", g)
	}
	groups, b := GetCtiDb().GetAgentByGroups("3001")
	if b {
		fmt.Printf("GetAgentGroups ok%v", groups)

	}*/
	//GetCtiDb().WriteCallRecord("12323xx434","werwerer",time.Now(),time.Now())
	//GetCtiDb().WriteCallRecord("123234xxx34","werwerer",time.Now(),time.Now())
}
