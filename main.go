package main

import (
	"MyPBFT_1/network/server"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
)

func main() {
	var (
		Serverurl string
		ServerPort string
		ServerID string
		PeerUrl string
		PeerPort string
		PublicKey string

	)
	flag.StringVar(&Serverurl,"server","","-server=127.0.0.1")
	flag.StringVar(&ServerPort,"port","","-port=10024")
	flag.StringVar(&ServerID,"serverid","","-serverid=1001")
	flag.StringVar(&PeerUrl,"peerurl","","-peerurl=127.0.0.1")
	flag.StringVar(&ServerPort,"peerport","","-peerport=20024")
	flag.StringVar(&PublicKey,"publickey","","")
	//flag.StringVar(&PrivateKey,"privatekey","","")
	flag.Parse()
	if  Serverurl==""{
			url,err:=server.GetInternet()
			if err!=nil{
				fmt.Println(err)
				return
			}
			Serverurl=url
	}
	if ServerPort==""{
		//ServerPort ="10024"
		ServerPort=server.FindPort(10024)

	}
	if ServerID==""{
		ServerID=Serverurl+strconv.Itoa(rand.Intn(10000))
	}
	if PeerUrl==""&&Serverurl!="" {
		PeerUrl = Serverurl
	}
	if PeerPort=="" {
		PeerPort=server.FindPort(20024)
	}
	//create instance of Server
	server:=server.NewServer(ServerID,Serverurl+":"+ServerPort,Serverurl+":"+ServerPort,PublicKey)
	//Start server

	server.Start()
}
//func usage()  {
//	fmt.Fprintf(os.Stdout, "please run \"%s --help\" and get help info\n", os.Args[0])
//	os.Exit(-1)
//}
