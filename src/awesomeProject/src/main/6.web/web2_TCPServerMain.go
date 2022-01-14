package main

import myWeb "awesomeProject/src/com.jsflzhong/6_web"

/*
net 包中有相应功能的函数，函数定义如下：
	func ListenTCP(net string, laddr *TCPAddr) (l *TCPListener, err os.Error)
	func (l *TCPListener) Accept() (c Conn, err os.Error)

ListenTCP 函数会在本地 TCP 地址 laddr 上声明并返回一个 *TCPListener，net 参数必须是 "tcp"、"tcp4"、"tcp6"，
如果 laddr 的端口字段为 0，函数将选择一个当前可用的端口，可以用 Listener 的 Addr 方法获得该端口。
 */

func main() {
	myWeb.BuildTCPServer()
}



