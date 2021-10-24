package znet

import (
	"fmt"
	"net"
	"reflect"

	"github.com/yuenshui/zinxplus/utils"
	"github.com/yuenshui/zinxplus/ziface"
)

// Server 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler ziface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr ziface.IConnectionManager

	// 消息列表
	MsgMap *map[uint32]interface{}

	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)

	packet ziface.Packet
}

// NewServer 创建一个服务器句柄
func NewServer(option *utils.ServerOptions) ziface.IServer {

	s := &Server{
		Name:       option.Name,
		IPVersion:  "tcp4",
		IP:         option.Host,
		Port:       option.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnectionManager(),
		packet:     NewDataPack(),
	}

	//for _, opt := range opts {
	//	opt(s)
	//}

	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	//开启一个go去做服务端Linster业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}

		//已经监听成功
		fmt.Println("start Zinx server  ", s.Name, " succ, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法

		cID := utils.SessionId()

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				err := conn.Close()
				if err != nil {
					continue
				}
				continue
			}

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, conn, cID, s.msgHandler)

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msg interface{}, router ziface.IRouter) {
	fmt.Println("Server AddRouter")
	msgID := s.GetMsgId(msg)

	if msgID != 0 {
		s.msgHandler.AddRouter(msgID, router)
		fmt.Println("AddRouter success:", msgID, reflect.TypeOf(msg))
	} else {
		fmt.Println("AddRouter fail:", msgID, reflect.TypeOf(msg))
	}
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnectionManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetMsgMap(MsgMap *map[uint32]interface{}) {
	s.MsgMap = MsgMap
}

func (s *Server) GetMsgMap() *map[uint32]interface{} {
	return s.MsgMap
}

func (s *Server) GetMsgId(msg interface{}) uint32 {
	msgID := uint32(0)
	for k, v := range *s.MsgMap {
		if reflect.TypeOf(msg) == reflect.TypeOf(v) {
			msgID = k
		}
	}

	return msgID
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

func (s *Server) Packet() ziface.Packet {
	return s.packet
}

func init() {
}
