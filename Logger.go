package main



import (
	"winlog"
	"fmt"
	"log"
	"time"
	//"container/list"
)
type Priority int
const (
	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)
func PriorityStr(p Priority)(str string){
	switch p{
		case LOG_EMERG:
			str = "Emerg"
		case LOG_ALERT:
			str = "Alert"
		case LOG_CRIT:
			str = "Crit"
		case LOG_ERR:
			str = "Err"
		case LOG_WARNING:
			str = "Warning"
		case LOG_NOTICE:
			str = "Notice"
		case LOG_INFO:
			str = "Info"
		case LOG_DEBUG:
			str = "Debug"
	}
	return ""
}
type LogWriter struct {
	priority Priority
	prefix   string
}
func (w *LogWriter) Write(p []byte) (n int, err error){
    str:=string(p)
    winlog.Log(str)
    return 0,nil
}
func (w *LogWriter) writeString(p Priority, s string) {
	if w.priority < p {
		return 
	}
	header :=TimeToString16(time.Now())

	header+=" "

	str:=fmt.Sprintf("%s %s: %s\r\n", PriorityStr(p), w.prefix, s)
	str=header+str
	winlog.Log(str)
}
func NewLog(prefix string, flag int) *log.Logger{
	w:=LogWriter{}
	//return &log.Logger{out: w, prefix: prefix, flag: flag}
	return log.New(&w,prefix,flag)
}
type LoggerManger struct{
	Loggers [] interface {}
}
var pLoggerManger *LoggerManger
func GetLoggerManger() *LoggerManger{
	if pLoggerManger == nil {
		pLoggerManger = new(LoggerManger)
	}
	return pLoggerManger
}
func SetLoggerPriority(prefix string, priority Priority){
	//e := GetLoggerManger().Loggers.Front()
	//for ; e != nil; e = e.Next() {
	//	if e.Value.(*LogWriter).prefix == prefix {
	//		e.Value.(*LogWriter).priority = priority
	//		break
	//	}
	//}
	//p :=GetLoggerManger()
	//for _,v :=range p.Loggers {
	//	if(v.(*LogWriter).prefix == prefix){
	//		//w = v.(*LogWriter)
	//		v.(*LogWriter).priority = priority
	//		break
	//	}
	//}
	w:=GetLogTag(prefix)
	w.priority = priority
}
func GetLogExistTag(tag string) bool{
	p :=GetLoggerManger()
	for _,v :=range p.Loggers {
		if(v.(*LogWriter).prefix == tag){
			//w = v.(*LogWriter)
			//break
			return true
		}
	}
	return false
}

func GetLogTag(tag string)(w *LogWriter){
	p :=GetLoggerManger()
	for _,v :=range p.Loggers {
		if(v.(*LogWriter).prefix == tag){
			w = v.(*LogWriter)
			break
		}
	}
	if(w==nil){
		lw:=LogWriter{
			priority:LOG_WARNING,
			prefix:tag,
			}
		p.Loggers = append(p.Loggers,&lw)
		w = &lw
	}
	return 
}
func NewLogger(prefix string, priority Priority) *LogWriter{
	w:=LogWriter{
		priority:priority,
		prefix:prefix,
	}
	GetLoggerManger().Loggers = append(GetLoggerManger().Loggers,&w)
	//.PushBack(&w)
	return &w
}
//func (w *Writer) Printf(format string, v ...interface{}) {
//	fmt.Sprintf(format, v...)
//}
func (w *LogWriter) Alert(format string, v ...interface{}) (err error) {
	//_, err = w.writeString(LOG_ALERT, m)
	w.writeString(LOG_ALERT, fmt.Sprintf(format, v...))
	return err
}
func (w *LogWriter) Crit(format string, v ...interface{}) (err error) {
	//_, err = w.writeString(LOG_CRIT, m)
	//w.writeString(LOG_CRIT, m)
	w.writeString(LOG_CRIT, fmt.Sprintf(format, v...))
	return err
}
// Err logs a message using the LOG_ERR priority.
func (w *LogWriter) Err(format string, v ...interface{}) (err error) {
	//_, err = w.writeString(LOG_ERR, m)
	//w.writeString(LOG_ERR, m)
	w.writeString(LOG_ERR, fmt.Sprintf(format, v...))
	return err
}
// Warning logs a message using the LOG_WARNING priority.
func (w *LogWriter) Warning(format string, v ...interface{}) (err error) {
	//_, err = w.writeString(LOG_WARNING, m)
	//w.writeString(LOG_WARNING, m)
	w.writeString(LOG_WARNING, fmt.Sprintf(format, v...))
	return err
}

// Notice logs a message using the LOG_NOTICE priority.
func (w *LogWriter) Notice(format string, v ...interface{}) (err error) {
	//_, err = w.writeString(LOG_NOTICE, m)
	//w.writeString(LOG_NOTICE, m)
	w.writeString(LOG_NOTICE, fmt.Sprintf(format, v...))
	return err
}

// Info logs a message using the LOG_INFO priority.
func (w *LogWriter) Info(format string, v ...interface{}) (err error) {
	//_, err = w.writeString(LOG_INFO, m)
	//w.writeString(LOG_INFO, m)
	w.writeString(LOG_INFO, fmt.Sprintf(format, v...))
	return err
}

// Debug logs a message using the LOG_DEBUG priority.
func (w *LogWriter) Debug(format string, v ...interface{}) (err error) {
	//_, err = w.writeString(LOG_DEBUG, m)
	//w.writeString(LOG_DEBUG, m)
	w.writeString(LOG_DEBUG, fmt.Sprintf(format, v...))
	return err
}