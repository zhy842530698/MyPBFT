package common

import "fmt"

type Myerror struct {
	Msg string
	No int
}
func(e Myerror)Error()string{
	//e.msg="失败"
	//e.errno=-1
	return fmt.Sprintf("err is happening %d ,errinfo %s",e.No,e.Msg)
}
