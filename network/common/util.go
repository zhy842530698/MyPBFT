package common

import (
	"reflect"
	"strconv"
	"strings"
)
//验证ip地址
func ValidIPAddress(IP string) string {
	v4CharNum := strings.Count(IP, ".")
	v6CharNum := strings.Count(IP, ":")
	if v4CharNum > 0 && v6CharNum > 0{
		return "Neither"
	}
	if v4CharNum > 0{
		return ValidIPv4(IP)
	}
	if v6CharNum > 0{
		return ValidIPv6(IP)
	}
	return "Neither"
}

func ValidIPv4(IP string) string {
	strs := strings.Split(IP, ".")
	if len(strs) != 4{
		return "Neither"
	}
	for i := range strs{
		if len(strs[i]) == 0{
			return "Neither"
		}
		if strs[i][0] == '0' && len(strs[i]) > 1{
			return "Neither"
		}
		curNum, err := strconv.Atoi(strs[i])
		if err != nil || curNum >= 256{
			return "Neither"
		}
	}

	return "IPv4"
}

func ValidIPv6(IP string) string {
	strs := strings.Split(IP, ":")
	if len(strs) != 8{
		return "Neither"
	}
	for i := range strs{
		if len(strs[i]) > 4 || len(strs[i]) == 0{
			return "Neither"
		}
		for j := range strs[i]{
			if strs[i][j] >= 'g' && strs[i][j] <= 'z' || strs[i][j] >= 'G' && strs[i][j] <= 'Z'{
				return "Neither"
			}
		}
	}
	return "IPv6"
}
func Struct2Map(obj interface{})map[string]interface{}  {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
func StructAtoB (a interface{},b interface{}) map[string]interface{}{
	//第一步,先将结构体转化为map方便后续遍历
	amap := Struct2Map(a)
	bmap := Struct2Map(b)
	//开始遍历A结构体的字段
	for k1, v1 := range amap {
		for k2,v2:= range bmap {
			typea := reflect.TypeOf(v1)
			typeb := reflect.TypeOf(v2)
			if k1 == k2 && typea==typeb {
				bmap[k2] = v1
			}
			//if k1==k2 && typea!=typeb {
			//
			//
			//	if typeb == tInterface {
			//
			//		//因为我们假设开放给前端的类型是字符串,所以我们直接用string转换为datetime
			//		//1.断言v1为string
			//		vs, ok := v1.(string)
			//		if ok {
			//			fmt.Println("开始转换number")
			//			if vs == "" {
			//				//v2 = nil
			//				//bmap[k2]=v2
			//				continue
			//			}
			//			bmap[k2] = v1
			//
			//		}
			//	}
			//}
		}
	}
	return bmap

}



