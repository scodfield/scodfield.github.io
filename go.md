1. 新的数据类型interface,channel,channel的运算符'<-',defer,recover与panic
2. select语句的case必须是一个通信操作,select随机选一个可运行的case,如果没有则阻塞,直到有case可运行,比较感兴趣的是如果有多个可运行的case,将会如何选
   能保证公平嘛，有优先级取舍嘛
3. 函数的形参就像定义在函数体内部的局部变量,这样就很好理解值传递,在调用函数时,将实际参数值复制一份赋值给形参,传递到函数中,所以值传递时对形参的修改不会
   影响实参的值;在此基础上引用传递,就是将实际参数的地址复制一份复制给形参,所以此时函数体内的局部变量形参所指向的是内存中的一个地址,任何对形参的修改
   都将改变地址内保存的值,而实参也指向该地址,表现上就是对形参的改变会同时影响到实参;函数另外需要注意的点是闭包和方法,闭包是匿名函数;至于方法,
   是带了接受者的函数,接受者是一个结构体对象或者指针,由于go中没有面向对象的概念,而像c++/java这种面向对象的语言,实现类的方法是编译器隐式的给函数加了
   一个this指针(第一个参数),而在go语言中,这个this指针需要显式的指明,方法的语法格式: func (var_name var_type) f_name [return_type] {} 
   方法示例: 
      type Circle struct {radius float64}  // 声明一个结构体
      func (c Circle) get_area() float64 { return 3.14* c.radius * c.radius } // 声明一个方法,可由Circle对象调用
4. 如何理解go语言中的闭包
5. go中的切片是对数组的抽象,数组在声明时已确定大小,切片则不同,可以认为切片是动态数组,它的长度是不固定的,可以追加元素,追加元素时切片的容量会扩大
   可以通过声明一个未指定大小的数组来定义切片,或通过make()来创建切片,定义方式如下:
   var slice1 []type;  // 未指定大小的数组
   var slice1 []type = make([]type, len) // make()创建,len是数组的长度,也是切片的初始长度,该方式可简写为: slice1 := make([]type,len)
   make() 有一个可选参数capacity,用于指定切片的容量
   切片的初始化方式: slice1 := []int {1,2,3}  or slice1 := arr[[start_index]:[endindex]] or slice2 := slice1[s_index:e_index]
   切片是可索引的,可通过len()函数获取切片的长度,cap()函数获取切片的容量
   可通过append()函数向切片追加元素,例:append(slice1,0) // 添加一个元素   append(slice1,1,2,3) // 同时添加多个元素
   可通过copy()函数在切片间拷贝,例:copy(slice2,slice1) // 将slice1的数据拷贝到slice2中,copy()需要注意的是声明slice2的时候要指定len的大小,
   即: slice2 := make([]int, 5),否则上述copy失败,错误声明:slice2 := make([]int,0) 则执行上述copy之后,len(slice2)=0,cap(slice2)=0,slice2=[]
   切片由三部分组成:指向底层数组的指针,len,cap
   基于数据或切片创建新切片时,新切片的大小和容量计算公式:slice2 := slice[i:j] 设定slice的cap=k,则slice2的len=j-i, cap = k-i
   上述的计算公式中len的计算好理解,对于cap而言,为啥是k-i,因为slice2的底层数组指针指向下表为i的元素,那么新的cap=原来的cap - i
   切片和数据的不同点,在函数调用时,切片是引用传递,数组是值传递,但是可以显式的取数据的地址
   所谓底层数组的指针,可通过一个例子来说明:var slice1 = []int{0,1,2,3,4,5}  slice2 := slice1[:3] slice1[1] = 10,则打印slice2 = [0,10,2]
   切片append()时,如果cap不足会扩大,新的cap计算方式:ceil((cap+len(arr))/2) * 2, 
   例:slice1 := []int {0,1,2,3} slice2 := slice1[2:] 
      slice2 = append(slice2,4,5,6,7,8)  // slice2.cap = ceil((slice2.cap + len([4,5,6,7,8])/2) * 2
5. go并发,通过go关键字开启goroutine线程,goroutine是轻量级线程,由go运行时进行管理,goroutine语法格式: go f_name(parameters), 通道channel可用于
   两个goroutine之间的通信,通道可由关键字chan和make()函数创建,例: ch := make(chan int),默认情况下通道是没有缓冲区的,发送端发送数据,同时必须有
   接收端接收数据,(经过测试,貌似没有缓冲区的通道,存储数据的部分是个栈结构,带缓冲区的是个队列),通道可通过make的第二个参数设置缓冲区大小,带缓冲区的
   通道允许发送端的数据发送和接收端的数据获取处于异步状态,也就是说带缓冲区的通道,可以先将数据放在缓冲区,等待接收端接收数据,如果不带缓冲区,发送方
   会阻塞直到接收端从通道接收数据,如果带缓冲,发送方会阻塞直到被发送的值被拷贝到缓冲区,缓冲区已满,则发送方会一直阻塞,直到接收方接收数据
   close()函数关闭通道
6. 参考系列
   a> 深度解密Go语言之map https://mp.weixin.qq.com/s/2CDpE5wfoiNXm1agMAq4wA
   b> [译] 我是如何在大型代码库上使用 pprof 调查 Go 中的内存泄漏 https://juejin.im/post/5ce11d1ee51d4510601117fd
   c> Golang 多版本管理器 https://github.com/voidint/g
   d> Go面试必考题目之method篇 https://mp.weixin.qq.com/s/US7MnIJfekJRazioxyWQhg
   e> 异常检测的N种方法，阿里工程师都盘出来了 https://mp.weixin.qq.com/s/w7SbAHxZsmHqFtTG8ZAXNg
7. Go内存分配:https://mp.weixin.qq.com/s?__biz=MzUzMjk0ODI0OA==&mid=2247483835&idx=1&sn=da048d277a12937e911d7fcbcf1ed11c&chksm=faaa35cccdddbcdaf38fe9e2060138164ad53c2d9e328d88944364cfd98f6991101846f7912e&mpshare=1&scene=23&srcid=#rd
8. Go中的指针分为两类:类型指针,允许对数据进行修改,但不能进行偏移和运算;切片,由指向起始元素的指针,元素数量和容量组成,切片比原始指针具备更强大的特性
   更安全,切片发生越界时,运行时会报宕机,并打印堆栈,而原始指针会崩溃,go中指针定义后,没有分配变量时,它的值为nil,这一点和c/c++不同,c/c++中声明后赋值
   前,指针变量可能指向任意地址,也就是野指针,另一个就是go是自动回收内存,不需要c/c++中手动回收,手动回收造成的问题包括,忘了回收导致内存泄漏,释放后未将
   指针赋值为NULL,导致悬空指针,以及多次释放导致的程序崩溃,由此可知,c/c++中指针的问题主要还是指针运算(越界访问,缓冲区溢出)和释放
   参考:http://c.biancheng.net/view/21.html
9. Go unsafe:https://mp.weixin.qq.com/s/JpHRe_XN9cqrP3KC8dOMqA
10. Go反射,反射是go的高级主题之一,反射就是程序能够在运行时检查变量和值,求出它们的类型,reflect包实现了运行时反射,reflect包会帮助识别interface{}变量
    的底层具体类型和值,reflect包的几个具体类型和方法如下:
    reflect.Type,reflect.Value表示interface{}变量的具体类型和值,对应的方法是reflect.TypeOf()和reflect.ValueOf();
    reflect.Kind 与reflect.Type类似,不同之处在于Type表示interface{}的实际类型,而Kind表示该类型的特定类别;
    reflect.NumField()方法返回结构体中字段的数量,reflect.Field(i int)方法返回字段i的reflect.Value;
    Int()和String()方法可以分别取出reflect.Value的特定数据类型
    参考: https://studygolang.com/articles/13178
