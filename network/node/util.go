package node

import (
	"bytes"
	"net/http"
)

func Send(url string, msg []byte)  {
	buff :=bytes.NewBuffer(msg)
	http.Post("http://"+url,"application/json",buff)
}
