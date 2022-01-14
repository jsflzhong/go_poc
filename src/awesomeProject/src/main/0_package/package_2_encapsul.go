package main

//导入其他自定义的包
import (
	encapsulation "awesomeProject/src/com.jsflzhong/0_package"
	"fmt"
)

/*
在Go语言中封装就是把抽象出来的字段和对字段的操作封装在一起，数据被保护在内部，程序的其它包只能通过被授权的方法，才能对字段进行操作。

封装的好处：
隐藏实现细节；
可以对数据进行验证，保证数据安全合理。

如何体现封装：
对结构体中的属性进行封装；
通过方法，包，实现封装。

封装的实现步骤：
将结构体、字段的首字母小写；
给结构体所在的包提供一个工厂模式的函数，首字母大写，类似一个构造函数；
提供一个首字母大写的 Set 方法（类似其它语言的 public），用于对属性判断并赋值；
提供一个首字母大写的 Get 方法（类似其它语言的 public），用于获取属性的值。

 */

func main() {
	testEncapsulation()
}

/**
使用:0_package包下的函数
*/
func testEncapsulation() {
	people := encapsulation.NewPeople()
	people.SetAge(1)
	age := people.GetAge()
	fmt.Println("@@@people",people,",age:",age)
}
