package main

import "fmt"

/*
接收器:
	位置:定义在普通方法的名字前,但在关键字func后.
	语法:func (aaa *T) bbb
	作用:与普通方法类似,唯一的不同,就是用上述接收器部分, 强行把该函数与指定的结构体绑定, 使得本函数可以用"指定结构体."的方式调用.

在Go语言中，结构体就像是类的一种简化形式，那么类的方法在哪里呢？
在Go语言中有一个概念，它和方法有着同样的名字，并且大体上意思相同，
Go 方法是作用在接收器（receiver）上的一个函数，接收器是某种类型的变量，因此方法是一种特殊类型的函数。

接收器类型可以是（几乎）任何类型，不仅仅是结构体类型，任何类型都可以有方法，甚至可以是函数类型，
可以是 int、bool、string 或数组的别名类型，但是接收器不能是一个接口类型，
因为接口是一个抽象定义，而方法却是具体实现，如果这样做了就会引发一个编译错误invalid receiver type…。

接收器也不能是一个指针类型，但是它可以是任何其他允许类型的指针，
一个类型加上它的方法等价于面向对象中的一个类，一个重要的区别是，
在Go语言中，类型的代码和绑定在它上面的方法的代码可以不放置在一起，它们可以存在不同的源文件中，唯一的要求是它们必须是同一个包的。

类型 T（或 T）上的所有方法的集合叫做类型 T（或 T）的方法集。

因为方法是函数，所以同样的，不允许方法重载，即对于一个类型只能有一个给定名称的方法，
但是如果基于接收器类型，是有重载的：具有同样名字的方法可以在 2 个或多个不同的接收器类型上存在，比如在同一个包里这么做是允许的。

提示
在面向对象的语言中，类拥有的方法一般被理解为类可以做的事情。
在Go语言中“方法”的概念与其他语言一致，只是Go语言建立的“接收器”强调方法的作用对象是接收器，也就是类实例，而函数没有作用对象。

接收器的格式如下：
	func (接收器变量 接收器类型) 方法名(参数列表) (返回参数) {
		函数体
	}

对各部分的说明：
	接收器变量：接收器中的参数变量名在命名时，官方建议使用接收器类型名的第一个小写字母，而不是 self、this 之类的命名。
		例如，Socket 类型的接收器变量应该命名为 s，Connector 类型的接收器变量应该命名为 c 等。
	接收器类型：接收器类型和参数类似，可以是指针类型和非指针类型。
	方法名、参数列表、返回参数：格式与函数定义一致。



注意:
	接收器有两种, 接收器根据接收器的类型可以分为"指针接收器"、"非指针接收器"
	两种接收器在使用时会产生不同的效果，根据效果的不同，两种接收器会被用于不同性能和功能要求的代码中。

	1.指针接收器: 参数是类型的指针, 外部调用时对结构体的字段的值的修改会一直有效.
	2.非指针接收器: 参数是类型,不是指针, 外部调用时对结构体的字段的值的修改会在调用结束后无效, 只是一份复制.

指针和非指针接收器的使用:
	在计算机中，小对象由于值复制时的速度较快，所以适合使用非指针接收器，
	大对象因为复制性能较低，适合使用指针接收器，在接收器和参数间传递时不进行复制，只是传递指针。

*/
func main() {
	receiver_pointer()

	receiver_nonPointer()
}

/*
种类一: 定义和使用"指针接收器".

接收器: 用特殊语法定义的函数, 使这个函数成为指定的结构体的函数(写在结构体的外面).

指针类型的接收器由一个结构体的指针组成，更接近于面向对象中的 this 或者 self。

注意:
	由于指针的特性，调用方法时，修改接收器指针的任意成员变量，在方法结束后，修改都是有效的!!

结果:
	在指针接收器中为其绑定的结构体赋值: &{[1]}
	在调用指针接收器之后的外部, 打印其绑定的结构体的字段值是否真的有改变: &{[1]}

注意:
	由于指针的特性, 调用方法时改变结构体的值, 调用后这些修改都是有效的. 上面的字段被改成了1
 */
func receiver_pointer() {
	bag := new(Bag)
	bag.PointerTypeReceiver(1)
	fmt.Println("在调用指针接收器之后的外部, 打印其绑定的结构体的字段值是否真的有改变:",bag)
}

/*
种类二: 定义和使用"非指针接收器".
当方法作用于非指针接收器时，Go语言会在代码运行时将接收器的值复制一份，在非指针接收器的方法中可以获取接收器的成员值，但修改后无效。

结果:
	在非指针接收器中为其绑定的结构体赋值: {[1]}
	在调用非指针接收器之后的外部, 打印其绑定的结构体的字段值是否真的有改变: &{[]}

注意:
	结果没有改变!
 */
func receiver_nonPointer() {
	bag := new(Bag)
	bag.NonPointerTypeReceiver(1)
	fmt.Println("在调用非指针接收器之后的外部, 打印其绑定的结构体的字段值是否真的有改变:",bag)
}

type Bag struct {
	items []int
}

/*
种类一: 定义和使用"指针接收器".

定义一个receiver接收器, 使得该函数成为上面自定义结构起的函数,强行绑定.
然后该函数就可以被上面结构体的指针直接调用了.

注意:
	由于指针的特性，调用方法时，修改接收器指针的任意成员变量，在方法结束后，修改都是有效的!!

名为接收器的部分:(bag *Bag)

注意:
	每个方法只能有一个接收器.
 */
func (bag *Bag) PointerTypeReceiver(itemId int)  {
	bag.items = append(bag.items, itemId)
	fmt.Println("在指针接收器中为其绑定的结构体赋值:",bag)
}

/*
种类一: 定义和使用"非指针接收器".

定义一个receiver接收器, 使得该函数成为上面自定义结构起的函数,强行绑定.
然后该函数就可以被上面结构体的指针直接调用了.

注意:
	当方法作用于非指针接收器时，Go语言会在代码运行时将接收器的值复制一份，在非指针接收器的方法中可以获取接收器的成员值，但修改后无效!!

名为接收器的部分:(bag Bag)

注意:
	每个方法只能有一个接收器.
 */
func (bag Bag) NonPointerTypeReceiver(itemId int)  {
	bag.items = append(bag.items, itemId)
	fmt.Println("在非指针接收器中为其绑定的结构体赋值:",bag)
}




