package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main(){
	////common.GenRsaKey(2048)
	//////加密
	////data := []byte("hello world")
	////encrypt := common.RSA_Encrypt(data, "./RSA/10024/public.pem")
	////fmt.Println(string(encrypt))
	////
	////// 解密
	////decrypt := common.RSA_Decrypt(encrypt, "./RSA/10024/private.pem")
	////fmt.Println(string(decrypt))
	//
	//	//获取时间戳
	//
	//	timestamp := int64(1627201141729072800)
	//
	//	fmt.Println(timestamp)
	//
	//	//格式化为字符串,tm为Time类型
	//
	//	tm := time.Unix(timestamp, 0)
	//
	//	fmt.Println(tm.Format("2006-01-02 03:04:05 PM"))
	//
	//	fmt.Println(tm.Format("02/01/2006 15:04:05 PM"))
	//
	//
	//
	//	//从字符串转为时间戳，第一个参数是格式，第二个是要转换的时间字符串
	//
	//	tm2, _ := time.Parse("01/02/2006", "02/08/2015")
	//
	//	fmt.Println(tm2.Unix())




	//base64Sig:= common.GetSign([]byte("zhy241626972757919569600"), "./RSA/10024/private.pem")
	//fmt.Println(string(base64Sig))
	//flag:= common.VerifySign([]byte("zhy241626972757919569600"), base64Sig, "./RSA/10024/public.pem")
	//fmt.Println(flag)
	//var Map map[string]string
	//Map = make(map[string]string,0)
	//Map["app"]="123"
	//for _,value:=range Map{
	//	fmt.Println(value)
	//}
	//s:="123"
	//res :=s+strings.Repeat(s,3)
	//fmt.Println(res)
//	l:=list.New()
//	//入队,压栈
//	l.PushBack(1)
//	l.PushBack(2)
//	l.PushBack(3)
//	l.PushBack(4)
//	//出队
//	i1:=l.Front()
//	l.Remove(i1)
//	fmt.Println(i1)
//// 出栈
//	i4:=l.Back()
//	l.Remove(i4)
//	fmt.Println(i4.Value)
	fmt.Println(computeString("3*[a2*[c]]3*[a2*[c]]"))
}
func computeString( str string ) string {
	// write code here
	strStack := []string{}
	numStack := []int{}
	num, res := 0, ""
	for _, char := range str {
		if char >= '0' && char <= '9' {
			n, _ := strconv.Atoi(string(char))
			num = num * 10 + n
		} else if char == '[' {
			strStack = append(strStack, res)
			numStack = append(numStack, num)
			num, res = 0, ""
		} else if char == ']' {
			str := strStack[len(strStack)-1]
			count := numStack[len(numStack)-1]
			strStack = strStack[:len(strStack)-1]
			numStack = numStack[:len(numStack)-1]
			res = str + strings.Repeat(res, count)
		} else if char == '*'{
			continue
		} else {
			res += string(char)
		}
	}
	return res
}