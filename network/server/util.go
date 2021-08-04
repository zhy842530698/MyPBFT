package server

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func GetInternet()(string,error){
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", errors.New("internal IP fetch failed, detail:" + err.Error())
	}
	defer conn.Close()

	// udp 面向无连接，所以这些东西只在你本地捣鼓
	res := conn.LocalAddr().String()
	res = strings.Split(res, ":")[0]
	return res, nil
}
func PortInUse(port int) bool {
	res := -1
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr %d", port)
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.Stdout = &outBytes
	cmd.Run()
	resStr := outBytes.String()
	r := regexp.MustCompile(`\s\d+\s`).FindAllString(resStr, -1)
	if len(r) > 0 {
		pid, err := strconv.Atoi(strings.TrimSpace(r[0]))
		if err != nil {
			res = -1
		} else {
			res = pid
		}
	}
	return res!=-1
}
func FindPort(port int )string {
	for PortInUse(port) {
		port++

	}
	return strconv.Itoa(port)
}