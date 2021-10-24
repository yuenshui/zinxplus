package znet

import "github.com/yuenshui/zinxplus/ziface"

// Request 请求
type Request struct {
	conn   ziface.IConnection //已经和客户端建立好的 链接
	msg    ziface.IMessage    //客户端请求的数据
	Params interface{}
}

// GetMsgMap 获取服务器
func (r *Request) GetMsgMap() *map[uint32]interface{} {
	s := r.conn.GetServer()
	return (*s).GetMsgMap()
}

// GetServer 获取服务器
func (r *Request) GetServer() *ziface.IServer {
	return r.conn.GetServer()
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

// GetSessionId 获取请求的消息的ID
func (r *Request) GetSessionId() string {
	return r.GetConnection().GetConnID()
}

// SetParam 设置参数
func (r *Request) SetParam(p interface{}) {
	r.Params = p
}

// GetParam 获取参数
func (r *Request) GetParam() interface{} {
	return r.Params
}

// SendMsg 发送信息到客户端
func (r *Request) SendMsg(msg interface{}) {
	r.GetConnection().SendBuffMsg(msg)
}
