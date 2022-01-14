package main

import photoServer "awesomeProject/src/com.jsflzhong/6_web/photoServer"

/**********
文件上传和浏览的服务器

注意, 每个handler相当于java中的一个controller,
但其本身不能作为网络endpoint存在, 需要在main方法中用:http.HandleFunc("/xxx", funcName) 来将其"注册"成一个endpoint,之后才可以被网络访问.
 */

/*
注册所有endpoints, 即网络端点.
只有用http.HandleFunc()这样的函数, 把我们自定义的函数注册后, 那些自定义的函数才能被网络访问到. 类似Controller中的handler.
 */
func main() {
	photoServer.RegistEndpoints()
}





