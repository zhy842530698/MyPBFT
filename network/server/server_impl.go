package server

import (
	//"MyPBFT_1/consensus"
	"MyPBFT_1/network/common"
	"MyPBFT_1/network/controller"
	controller2 "MyPBFT_1/network/controller/pbft"
	"MyPBFT_1/network/node"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	node 			*node.Node
	url 			string
	peerurl         string
	contype        Connectiontype
	longcon  	map[string]*net.Conn
	router   	*httprouter.Router
	urls     map[string]string
	ac        *Access
}
type Access struct {
	maxAccess int
	// 桶,用来限制最大并发数
	bucket chan struct{}
	// 用来判断何时桶内已空，否则需要循环判断
	wg sync.WaitGroup
	closed int32
}
type Connectiontype  int
const(
	longcon Connectiontype= iota
	shortcon Connectiontype = iota
)
func Makechan(){
	common.DataQueue=make(chan interface{},1024)
	common.RespQueue=make(chan *common.Myerror ,1024)
	common.HttpQueue = make(chan interface{},1024)
}
func NewServer(nodeID ,url,peerurl,public string)*Server  {
	node := node.NewNode(nodeID,url,public)
	Makechan()
	servermap:=make(map[string]string)

	server := &Server{node, url,peerurl,0,nil,httprouter.New(),servermap,nil}

	server.ac = NewAccessor(1024)

	server.setRoute()

	go server.Listen()

	return server
}
func (server *Server) Start() {
	fmt.Printf("Server is Listening by %s....Peer is Listening by %s..\n",server.url,server.peerurl)
	if err:=http.ListenAndServe(server.url,server);err!=nil {
		fmt.Println(err)
	}

}

//设置路由
func (server *Server)setRoute()  {

	go GracefulExit(server)
	server.AddRoute("POST","/join",&controller.JoinController{})
	server.AddRoute("POST","/req",&controller2.PBFTRequestController{})
	server.AddRoute("POST","/preprepare",&controller2.PBFTPrePrepareController{})
	server.AddRoute("POST","/prepare",&controller2.PBFTPrepareController{})
	server.AddRoute("POST","/commit",&controller2.PBFTCommitController{})
	server.AddRoute("POST","/reply",&controller2.PBFTReplyController{})
}
func (s *Server) AddRoute(method, path string, ctrl common.IController) {
	var handle httprouter.Handle = func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// http框架虽然会panic-recover但是自己也必须recover，因为接入记录panic后不会正常消去
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				w.Write([]byte("Server is busy111."))
				stack := make([]byte, 2048)
				stack = stack[:runtime.Stack(stack, false)]
				//f := "PANIC: %s\n%s"
				//logger.Error(f, err, stack)
			}
		}()
		err := s.ac.InControl()
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Server is busy!!"))
			return
		}

		defer s.ac.OutControl()
		nt := time.Now()
		// 打印输入请求

		// 解析输入参数
		idl := ctrl.GenIdl()
		body, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, idl)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("request info parse err"))
			return
		}

		do := func(r *http.Request, w http.ResponseWriter) {
			var data []byte
			resp := ctrl.Do(idl)
			if resp == nil {
				data = []byte(`{"errno":0,"errmsg":"ok"}`)
			}
			data, _ = json.Marshal(resp)
			et := time.Now().Sub(nt)
			fmt.Printf("request_urli=%s||response=%s||proc_time=%s\n",r.URL,string(data),et.String())

			//	r.URL, string(data), et.String())
			w.WriteHeader(200)
			w.Write(data)
		}
		fmt.Println("-------------")

		do(r, w)
	}
	s.router.Handle(method, path, handle)
}
func (server *Server) Listen(){
	for true {
		select {
		case msg:=<-common.HttpQueue:
			server.node.MsgEntrance<-msg
		}
	}
}
//实现http接口
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}


func GracefulExit(svr *Server) {
	sigc := make(chan os.Signal, 0)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	<-sigc
	println("closing agent...")
	svr.GracefulExit()
	println("agent closed.")
	os.Exit(0)
}
func (s *Server) GracefulExit() {
	s.ac.Stop()
}





