package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/yuenshui/zinxplus/ziface"
)

// ConnectionManager 连接管理模块
type ConnectionManager struct {
	connections map[string]ziface.IConnection
	connLock    sync.RWMutex
}

// NewConnectionManager 创建一个链接管理
func NewConnectionManager() *ConnectionManager {
	fmt.Println("NewConnectionManager")
	return &ConnectionManager{
		connections: make(map[string]ziface.IConnection),
	}
}

// Add 添加链接
func (connMgr *ConnectionManager) Add(conn ziface.IConnection) {
	fmt.Println("ConnectionManager add")

	connMgr.connLock.Lock()
	// 将conn连接添加到ConnectionManager中
	connMgr.connections[conn.GetConnID()] = conn
	connMgr.connLock.Unlock()

	fmt.Println("connection add to ConnectionManager successfully: conn num = ", connMgr.Len())
}

// Remove 删除连接
func (connMgr *ConnectionManager) Remove(conn ziface.IConnection) {
	fmt.Println("ConnectionManager Remove")

	connMgr.connLock.Lock()
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	connMgr.connLock.Unlock()
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

// Get 利用ConnID获取链接
func (connMgr *ConnectionManager) Get(connID string) (ziface.IConnection, error) {
	fmt.Println("ConnectionManager Get")
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")

}

// Len 获取当前连接
func (connMgr *ConnectionManager) Len() int {
	fmt.Println("ConnectionManager Len")
	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()
	return length
}

// ClearConn 清除并停止所有连接
func (connMgr *ConnectionManager) ClearConn() {
	fmt.Println("ConnectionManager ClearConn")
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	connMgr.connLock.Unlock()
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}

// ClearOneConn  利用ConnID获取一个链接 并且删除
func (connMgr *ConnectionManager) ClearOneConn(connID string) {
	fmt.Println("ConnectionManager ClearOneConn")
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connections := connMgr.connections
	if conn, ok := connections[connID]; ok {
		//停止
		conn.Stop()
		//删除
		delete(connections, connID)
		fmt.Println("Clear Connections ID:  ", connID, "succeed")
		return
	}

	fmt.Println("Clear Connections ID:  ", connID, "err")

}
