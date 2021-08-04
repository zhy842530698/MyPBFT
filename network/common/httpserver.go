package common

//type Controller interface {
//	GenIdl()  interface{}
//	Do (interface{})interface{}
//}
//type httpserver interface {
//	getReq(writer http.ResponseWriter, request *http.Request)
//	getPrePrepare(writer http.ResponseWriter, request *http.Request)
//	getPrepare(writer http.ResponseWriter, request *http.Request)
//	getCommit(writer http.ResponseWriter, request *http.Request)
//	getReply(writer http.ResponseWriter, request *http.Request)
//}
type IController interface {
	GenIdl() interface{}
	Do(interface{}) interface{}
}
