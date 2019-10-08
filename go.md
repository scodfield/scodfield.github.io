1. 新的数据类型interface,channel,channel的运算符'<-',defer,recover与panic
   a> go的数据类型中,只有切片slice,字典map和信道channel是引用类型,其它都是值类型
   b> defer延迟调用函数,含有defer语句的函数,在函数返回前,调用另一个函数,main() {defer hello() xxx },在main()函数返回前调用hello()函
   数,不仅适用于函数,还可以延迟调用(结构体)方法,defer语句的实参取值是在执行defer语句的时候,而非在调用延迟函数的时候,
   main() { a := 5 defer printA(a)
   a = 10 fmt.Println("before defer a:" a)} func printA(a int) { fmt.Println("in defer a:",a) },上述延迟调用结果:"in defer a 5"
   defer栈,当在一个函数内多次调用defer语句时,Go会把defer调用放入一个栈中,再按照LIFO的顺序执行(可以实现字符串逆序输出)
   当一个函数应该在与当前代码流无关的环境下调用时,可以使用defer,比如用到sync.WaitGroup的地方,在协程函数内声明一个: defer wg.Done()
   c> mutex 是sync包中的一个结构体类型,它主要定义了Lock()和Unlock()这两个方法,用于在Go中提供了一种处理竟态条件(race confition)的
   加锁机制(locking mechanism),可确保在某时刻只有一个协程在临界区(critical section)运行,防止出现竟态条件,如果有协程持
   有了锁(Lock()),当其它协程试图获得该锁时,这些协程会被阻塞,知道mutex变量解锁锁定(Unlock())为止,另外一种实现加锁的方法是使用非缓冲
   信道,在操作临界区前发送数据,操作后接收数据即可
   d> func panic(interface{})会终止程序的执行,并在延迟defer函数(如果定义了延迟函数)执行完之后,程序控制返回执行该panic()函数
   调用的调用方,退出过程一直进行直到当前协程的所有函数都退出,然后打印panic信息,接着打印堆栈跟踪(stack trace),发生panic时,recover
   可重新获得对该程序的控制,panic-recover与其它语言的try-catch-finally类似,调用panic的两种场景:发生不能恢复的错误,程序无法在继续
   运行下去(如,端口被其它程序占用,导致当前web应用绑定端口失败);发生了一个编程上的错误(如,一个非空的参数,在实际调用时为nil),上述
   退出过程会一直持续到main()(主协程)退出为止,跟踪堆栈则是从当前调用panic的协程函数开始直到main(),运行时错误(如,数组越界)也会
   导致panic,等价于调用了内置函数panic,其参数由接口类型runtime.Error给出;
   e> func recover()interface{} 也是内建函数,用于重新获得panic协程的控制,只有在延迟函数内部,调用recover才有用,在延迟函数内
   调用recover,可以取到panic的错误信息,停止panic续发事件(panicking sequence),程序恢复正常运行,如果在defer外部调用recover,不能停止panic
   续发事件,只有在同一个go协程中调用recover才行,recover不能恢复一个不同协程的panic,recover恢复panic之后,堆栈跟踪就被释放,但是可以通过
   runtime/debug.PrintStack()函数打印堆栈跟踪,导入debug包,并在defer的recover函数里调用即可:
   import "runtime/debug" 
   defer recov()   
   func recov(){if r := recover(); r != nil
     fmt.Println("recovered: ", r) debug.PrintStack()}
   在Printf()函数中,使用'%T'格式说明符,可以打印出变量的类型,unsafe包提供了一个Sizeof()函数,该函数接收变量并返回它的字节大小,
   不过unsafe包可能会带来可移植性问题,因此需要小心使用;
   go有着非常严格的强类型特征,没有自动类型提升/转换,比如在C语言中,整型可以和浮点型变量相加,但go中会报错(invalid operation),操作符
   两侧的变量类型必须一致,类型转换为T(v),其中v为变量,T为系统变量类型(int,string,float),类型转换举例:i := int(56.78)
   go中字符串是字节的切片,'+'操作符可用于拼接字符串,go中的字符串兼容unicode编码,并使用utf-8编码,字符串的遍历,可以通过:
   for xxx len(str){},不过这种方式是逐个遍历每个字节,对于特殊字符可能会有问题,可使用下述rune切片转化,另一种遍历方式为for range{},
   该方法可以避免上述特殊字符问题(自动区分代码点),Go中的字符串是不可变的,它可以通过类似于数组下标的方式访问,但是一旦赋值,不能通过
   下标方式更改,可以将字符转换为rune切片,再更改即可;
   rune是go内建类型,int32的别称,rune表示一个代码点(code point),代码点无论占用多少个字节,都可以用一个rune来表示,在遍历字符串时,
   特殊字符占用2个字节以上的情况下,如果适用'%c'会发现输出字符和原字符不一致,这时候就需要rune类型,
   例: runes := []rune("Señor"),这时候再用"%c",即可正确打印"ñ",使用方法就是把字符串转化为一个rune切片
   Go结构体可以声明匿名字段和匿名结构体,匿名字段默认该字段类型即是该字段名称,type n_s struct {name string int},
   var n_var n_s, n_var.int = 18
   如果结构体中有匿名的结构体类型字段,则该匿名结构体里的字段就称为提升字段,提升字段就像输入外部结构体一样,外部结构体变量
   和指针可以直接通过"."点号操作符访问,
   type Addr struct {state, country string} type Person struct { name string age int Addr}, var p Person fmt.Println(p.state)
   结构体名称以大写字母开头,则是其它包可以访问的导出类型(exported type),结构体是值类型,如果每一个字段都是可比较的,那么该结构体
   也是可比较的,如果两个结构体变量的对应字段相等,则这两个结构体变量是相等的,如果结构体包含不可比较字段(如:map类型字段),
   则结构体变量也不可比较(go只有三种引用类型);
   Go方法(methods),方法就是一个函数,它在func关键字和函数名之间加入了一个特殊的接收器类型,接收器可以是结构体类型,也可以是
   非结构体类型,接收器可以在方法内部访问,接收器类型类似于面向对象中的this指针,所以可以用结构体和方法模拟面向对象编程,同时一个方法可以
   被定义到多个不同类型的结构体上,可以用来实现多态,结构体可以使用指针接收器和值接收器,二者的区别是,在指针接收器的方法内部的改变,
   对调用者是可见,而值接收器则不然(类似函数的值和指针传递),使用指针接收器的场景:在方法内部对接收器的改变必须对调用者可见时;
   结构体字段太多,拷贝结构体的代价太大时;其它场景使用值接收器都可以接受,
   与匿名结构体字段一样,如果匿名结构体实现了一个方法,外部结构体也可以直接访问该方法,比如:func (a Addr) fullAddr{fmt.Println(xxx)},
   可以直接:p.fullAddr(),值接收器和函数的值参数类似,但也有区别,当函数定义为值参数时,它只能接受值参数,但是当一个方法为值接收器时,它可以同时
   接受值接收器和指针接收器,同样的指针参数和指针接收器,指针参数的函数只接受指针,使用指针接收器的方法可以使用指针接收器和值接收器
   对于在非结构体类型上定义方法,需保证方法的接收器类型定义和方法的定义在同一个包中,比如在自己项目的main包中如下定义:
   package main func (a int) add
   (b int) {}  func main() {},
   则编译报错:cannot define new methods on non-local type int, 这是因为add方法的定义和int这个系统类型的定义不在同
   一个包中,解决办法是在当前包中,为内置类型int创建一个类型别名,创建以该类型别名为接收器的方法:
   type myInt int func (a myInt) add(b myInt) myInt{}
2. a> select语句的case必须是一个通信操作,select随机选一个可运行的case,如果没有则阻塞,直到有case可运行,比较感兴趣的是如果有多个可运行
   的case,将会如何选择,能保证公平嘛，有优先级取舍嘛,由以下解析可知,go底层将所有case语句打乱顺序,一个一个循环检测是否channel是否可读或可写,
   select底层解析:https://mp.weixin.qq.com/s?__biz=MzUzMjk0ODI0OA==&mid=2247483766&idx=1&sn=eb605a64bed0b2066a12083f26fb04b6&chksm=faaa3501cdddbc177121ba14a6604743d5ea881ca8299d5609ac8eb9b6eca4f2a142ad5aabfd&token=1212449367&lang=zh_CN&scene=21#wechat_redirect
   b> 若协程阻塞,比如一个RPC请求阻塞了goroutine的执行,此时go运行时会新建M,继续运行P中的其它协程,那么对于一个大型系统来说,若RPC调用量
      非常的大,此时改如何控制生成的M的数量,大量系统线程的生成和回收会对系统性能及go的GC造成很大压力;
   c> gc时首先是标记根对象,根对象root一般包括:全局变量,各个G中statck上的变量等
      golang gc的第一阶段一般称之为:stack scan,具体描述是:collect pointers from globals and goroutine stacks,
      gc mark的最后一个阶段一般称为:mark termination,termination的第一个子阶段就是rescan(同时这个阶段也需要STW),
      对rescan的描述是:rescan globals/changed stacks
      综上可知,除了全局变量,所有active&changed goroutine stacks上的变量都要扫描
3. 函数的形参就像定义在函数体内部的局部变量,这样就很好理解值传递,在调用函数时,将实际参数值复制一份赋值给形参,传递到函数中,所以值传递时
   对形参的修改不会影响实参的值;在此基础上引用传递,就是将实际参数的地址复制一份复制给形参,所以此时函数体内的局部变量形参所指向的是
   内存中的一个地址,任何对形参的修改都将改变地址内保存的值,而实参也指向该地址,表现上就是对形参的改变会同时影响到实参;函数另外需要注意的是
   闭包和方法,闭包是匿名函数;
   至于方法,是带了接受者的函数,接受者是一个结构体对象或者指针,由于go中没有面向对象的概念,而像c++/java这种面向对象的语言,实现类的方法
   是编译器隐式的给函数加了一个this指针(第一个参数),而在go语言中,这个this指针需要显式的指明
   方法的语法格式: func (var_name var_type) f_name [return_type] {} 
   方法示例: 
      type Circle struct {radius float64}  // 声明一个结构体
      func (c Circle) get_area() float64 { return 3.14* c.radius * c.radius } // 声明一个方法,可由Circle对象调用
4. 如何理解go语言中的闭包,闭包是匿名函数与匿名函数所引用的环境的组合,这种特性使得匿名函数不用通过参数传递的方式,就可以直接引用外部变量,
   类似于常规函数直接引用全局变量,闭包同时保存有一个中间状态,当闭包作为函数返回值,并赋值给一个变量时,该变量不仅保存了匿名函数的地址,
   同时还有整个闭包的状态,状态会一直保存在这个被赋值的变量中,直到变量被销毁,整个闭包也被销毁
5. go中的切片是对数组的抽象,数组在声明时已确定大小,切片则不同,可以认为切片是动态数组,它的长度是不固定的,可以追加元素,追加元素时
   切片的容量会扩大,可以通过声明一个未指定大小的数组来定义切片,或通过make()来创建切片,定义方式如下:
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
   上述的计算公式中len的计算好理解,对于cap而言,为啥是k-i,因为slice2的底层数组指针指向下表为i的元素,那么新的cap=原来的cap - i,也就是说切片的容量
   是从创建切片索引开始的底层数组中元素的个数
   切片和数据的不同点,在函数调用时,切片是引用传递,数组是值传递,但是可以显式的取数据的地址
   所谓底层数组的指针,可通过一个例子来说明:var slice1 = []int{0,1,2,3,4,5}  slice2 := slice1[:3] slice1[1] = 10,则打印slice2 = [0,10,2]
   切片append()时,如果cap不足会扩大,新的cap计算方式:ceil((cap+len(arr))/2) * 2, 
   例:slice1 := []int {0,1,2,3} slice2 := slice1[2:] 
      slice2 = append(slice2,4,5,6,7,8)  // slice2.cap = ceil((slice2.cap + len([4,5,6,7,8])/2) * 2
   切片持有对底层数组的引用,只有切片在内存中,数组就不会被垃圾回收,如果有一个较大的数据,但我们只处理其中的一部分数据,由该数组创建切片时,只要仍在使用
   切片,数组就不会被回收,解决办法是使用func copy(dst, src []T)int函数来生成一个切片副本,这样可以使用新切片,原始数组被垃圾回收
   记一个小坑:在爬取豆瓣电影top250的时候,解析函数parseUrl有一个slice参数,在外层调用声明:var results []Result,解析后在外层调用时,发现results为空,
   按说切片是地址传递,为啥最后还是空呢,这是因为在parseUrl函数中调用了append函数,在添加元素的时候,切片的空间不够,分配了新的空间,在外层调用函数中
   results还是指向原来的地址,所以正确的做法是在外层调用函数中:results = parseUrl(url,results),相应的parseUrl/2函数要将最后新的results返回,具
   体实现参见douban_top250.go
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
   f> 浅谈Go语言实现原理:https://draveness.me/golang/
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
11. Go包(package)用于组织go源代码,提供了更好的可重用性,可读性和可维护性,所有可执行的go程序都必须包含一个main函数,main函数应该放置于main包中
    package packagename 这一行代码指定了该源文件属于哪一个包,放在每一个源文件的第一行,go install dir_name 编译dir_name命名的程序,会在dir_name
    中寻找包含main函数的文件,并生成dir_name命名的二进制或可执行文件,属于某一个包的源文件都应该放置于一个单独的文件夹内,例如dir_name程序里的一个工具
    类的源文件,在dir_name文件夹内创建util文件夹,并在新util内创建util.go(./dir_name/util/util.go),util.go的第一句就是:package util
    所有包都可以包含一个init函数,init函数没有任何参数和返回值,也不能显式的调用它,常用于初始化任务及开始执行之前检查程序的正确性
    包的初始化顺序如下:
      a> 首先初始化被导入的包,如果一个包导入了另一个包,则会先初始化被导入的包,一个包可能会被多次导入,但只会被初始化一次;
      b> 包级别的变量;
      c> 然后调用init函数,包可以有多个init函数(分布于多个文件中)
    导入包却不使用它,编译器编译时会报错,也就是说导入不使用的包在go中是非法的,这么做的目的是为避免导入过多未使用的包,导致编译时间显著增加,有时会先
    导入暂时用不到的包,此时可通过空白标识符来规避上述问题(与erlang类似),比如: var _ = math.Sqrt(3 * 5)
12. Go协程(goroutine),是与其它函数或方法并发运行的函数或方法,go通过协程和信道(channel)来处理并发
    go协程可以看作是轻量级的线程,创建一个go协程的成本很小(与erlang类型,spawn一个erlang进程消耗也很小),go协程相比系统线程的优势如下:
    a> 相比线程,go协程的成本极低,堆栈大小只有若干kb,并可根据应用增减,线程必须指定堆栈大小,且固定不变
    b> go协程会复用数量较少的os线程(类似erlang中beam的调度线程),若该线程中的某一个协程发生阻塞,系统会新建os线程,并把其它剩余的协程迁移到新线程上
    c> go协程之间使用信道来进行通信(erlang则采用基本消息的actor模型)
    启动一个新协程,协程的调用会立即返回,程序接着执行下一行代码,与直接调用函数的不同就在于,并不会等待新协程的任何返回值;如果希望运行go协程,则主协程
    必须继续运行,主协程终止,则程序终止
    注: go程序的主函数main()会运行在一个特殊的协程上,称为go主协程main goroutine
13. Go信道(channel),是协程之间通信的管道,所有信道都关联了一个类型,信道只能传输这种类型的数据,传输其它类型的数据都是非法的(erlang毕竟是动态语言
    可以发任意类型的消息,只要接收方可以有效处理即可),信道的零值为nil,零值没有任何意义,与map和切片一样,应该用make()来创建信道,信道通过箭头'<-'来
    标识是发送数据还是接受数据,'<-'箭头方向背离信道变量,表明从信道取数据,'<-'箭头方向指向信道变量,表明向信道写数据
    信道的发送和接受都是阻塞的,当向信道发送数据时,程序控制会在发送数据的语句处阻塞,直到有其它go协程从该信道读取数据时才会解锁阻塞,读取过程类似,
    如果没有协程向该信道写入数据,则读取协程会一直阻塞,这种特性使得协程之间通信时不需要使用锁或其它条件变量(erlang采用异步消息,当没有可处理的消息时
    陷入空转,不过erlang启动beam时可以指定启动的os调度线程数量,而go如果出现大量协程阻塞,岂不是要启动N多os线程?)
    信道需要避免的是死锁,当向一个信道发送数据,而没有任何协程接收,或者从某一个信道读取数据,但是没有任何协程向该信道写数据,此时会触发panic,返回运行时
    错误,一般是fatal error:all goroutine are asleep - deadlock
    通常创建的信道都是双向信道,既可以发送数据,也可以接收数据,不过在某些情况下可以把双向信道转换为单向信道(send/receive-only),但是反过来不行,比如
    可以在主协程创建双向信道,新协程函数的信道参数定义为单向信道类型,send-only: func xxx (sendOnlyChannel chan<- int){}
    receive-only: func xxx (receiveOnleyChannel <-chan int) {},在主协程创建双向信道:add_ch := make(chan int) go xxx(add_ch)  
    add_res := <-add_res fmt.Println(add_res),则在xxx协程里add_ch为只能发送数据的单向信道,而在主协程里是一个双向信道,仍然可以接收数据
    close(ch_variable)用来关闭信道,一般由发送数据的协程来关闭信道,在接收数据时,可通过多一个参数来判断信道是否关闭: v, ok := <- ch_variable
    当ok == false时,表示此时正在读取一个已关闭的信道,v 则是信道类型的零值,比如读取一个已关闭的chan int,此时v的值为0,这种方式需要自己判断信道是否
    已关闭,最常用的做法是使用for range循环,for range循环会一直从信道中接收数据,一旦关闭了chan,循环就自动结束:for v := range ch_variable {xxx}
    缓冲信道(buffered channels),通过向make()函数再传递一个表示容量的参数(表示缓冲区大小),就可以创建一个缓冲信道,无缓冲信道的接收和发送都是阻塞的
    而缓冲信道只在缓冲区已满的情况下才会阻塞向信道发送数据,同样也只在缓冲为空的情况下才会阻塞从信道接收数据
    WaitGroup用于等待一批协程执行结束,来自于sync包,程序会一直阻塞,直到这些协程全部执行完毕,WaitGroup是一个结构体类型,WaitGroup使用计数器来工作
    WaitGroup.Add()方法增加计数,WaitGroup.Done()方法减少计数,WaitGroup.Wait()方法会阻塞调用它的协程,直到计数器为0时才会解除阻塞
    缓冲信道的重要应用之一就是工作池,工作池就是一组等待任务分配的线程,一旦完成所分配的任务,就可以继续等待分配任务
    select,select语句用于在多个发送/接收信道操作中进行选择,select语句会一直阻塞,直到有发送/接收操作准备就绪,如果有多个信道操作准备完毕,select会
    随机选择其中的一个执行,select语法与switch相似,只不过每个case语句都是信道操作,使用default语句执行默认操作,防止select语句一直阻塞;
14. Go Runtime,go的运行时管理着调度(scheduler),垃圾回收(gc)和goroutine的运行环境
    go运行时负责运行goroutine,并把它们映射到操作系统线程上,每个goroutine由一个G结构体来表示,结构体字段用来维护此goroutine的栈和状态,运行时管理G
    并把它们映射到Logical Processor(P),P可以看做是是一个抽象的资源或者上下文,操作系统线程(M)获取P以便运行G
    当使用go命令创建goroutine时,P会先将创建的G放入local queue(由P维护),如果local queue已满,则将G放入global queue(schedt结构体),为了运行G,M需要
    持有P,M从P的queue中弹出一个goroutine并执行,当M执行一些G之后,如果local queue为空,它会随机选择一个P,并从它的queue中取走一半G到自己的queue中
    继续执行,当goroutine执行阻塞的系统调用时,如果P中还有一些G等待执行,则运行时会把P从这个阻塞线程中摘除(detach),唤醒一个空闲的M或者创建一个
    新M来服务于这个P,当阻塞的M继续的时候,goroutine放入global queue中等待调用,M则park它自己(休眠),加入到空闲线程中,运行时会在以下情况下阻塞并切换
    新的goroutine:blocking syscall(eg. file operation); newwork input; channel operations; primitives in the sync package(eg. waitGroup?)
    Go通过GODEBUG环境变量来跟踪运行时的调度器,eg: GODEBUG=scheddetail=1,schedtrace=1000 ./program ,GODEBUG输出主要包括M,P,G的概念及相关状态
    Go还有一个图形化的工具go tool trace 用于查看程序和运行时的状况
    goroutine的执行是可以被抢占的,运行时在启动程序时,会自动创建一个系统线程运行sysmon()函数,sysmon()在整个程序声明周期一直运行,负责监视各个go协程
    的运行状态,判断是否需要垃圾回收等,sysmon()调用retake()函数,retake()会遍历所有的P,如果一个P处于执行状态,并已经连续执行很长时间,就会被抢占,
    retake()调用preemptone()将P的stackguard0置为stackPreempt,这将导致P中正在执行的G在下一次函数调用时,栈空间检查失败,进而触发morestack()等一
    系列函数调用,morestack() --> newstack() --> gopreempt_m() --> goschedImpl() --> schedule(), goschedImpl()函数中会调用dropg()函数
    将G与P和M解除绑定,再调用globrunqput()将G放入global queue,最后调用schedule()为P设置新的可执行G
    Go调度器的schedule()和findrunnable()函数,goroutine调度是在P中进行的,每当运行时需要调度时会调用schedule()函数(proc1.go),schedul()函数先调用
    runqget()从当前P的local queue中取一个可执行的G,如果队列为空,继续调用findrunnable()函数,findrunnable()按照以下顺序取得G:调用runqget()函数
    从当前P的队列中取一个可执行G(与schedule()相同); 调用globrunqget()从全局队列中取可执行的G; 调用netpoll()取异步调用结束的G,该次为非阻塞调用,
    立即返回; 调用runqsteal()从其它P的队列中偷(类似与erlang的迁移) , 如果以上四步还未成功,继续执行以下低优先级的工作:如果处于垃圾回收的标记阶段,则
    进行垃圾回收标记工作; 再次调用globrunqget()从全局队列取可执行的G; 再次调用netpoll()取异步调用结束的G,该次调用未阻塞调用 , 如果还没有获得G,当前
    M停止执行,返回runnable()函数从头开始执行,如果findrunnable()正常返回一个可执行的G,schedule()函数会调用execute()函数执行该G,execute()函数调用
    gogo()函数,gogo()从G.sched结构体中恢复出G上次被调度器暂停时的寄存器现场(SP,PC),然后继续执行
    另外一种关于G长时间占用M的说法:
      go启动时启动sysmon,记录P中G的计数schedtick,schedtick会在每执行一个G之后递增; 如果检查到schedtick一直没有递增,说明这个P一直在执行同一个G,
      如果超过10ms(协程的切换时间片),就在这个G的栈里加一个tag标记; G在执行的时候,如果遇到非内联函数调用,就会检查一次这个标记,中断自己,把自己放到
      队列末尾,执行下一个G; 如果没有遇到非内联函数调用,就会一直执行这个G,直到结束,如果是死循环,且GOMAXPROCS=1,那么只有一个P&M,队列中的其它G不会
      执行,如果GOMAXPROCS大于1,则会新建M,并将P迁移到新的M,继续执行下一个G
    参考: https://studygolang.com/articles/10094; https://studygolang.com/articles/10095;
         GC: GC:https://studygolang.com/articles/7516; goroutine调度: https://studygolang.com/articles/10115
15. Go的垃圾回收,笔记参见miscellaneous.md
    参考:
    https://segmentfault.com/a/1190000012597428
    https://mp.weixin.qq.com/s?__biz=MzU4ODczMDg5Ng==&mid=2247483688&idx=1&sn=46742e533886fe8b2fb91d79cf5144eb&scene=21#wechat_redirect
16. some worthy tips
    a> go中,函数是一等公民,函数和其它类型一样,可以被赋值,传递给函数或者从函数中返回,但函数两点需要注意: 函数值类型不能作为map的key(可以作为value);
       函数值之间不能比较,只能和nil比较,函数类型的零值是nil
    b> 匿名函数赋值给变量,中间状态也被保存在变量中,那么在对匿名函数遍历的时候(创建了一个匿名函数slice或数组),匿名函数访问的外部变量是遍历结束后最终
       的变量值,因为多个闭包共享这个中间状态,示例如下:
       var nonSlice []func()
       strSlice := []string {"1","3","5","7"}
       for _, str := range strSlice {
           // temp := str
           nonSlice = append(nonSlice, func() { fmt.Println(str) })
           // nonSlice = append(nonSlice, func() { fmt.Println(temp) })
       }
       for _, val := range nonSlice {
           val()
       } 
       输出结果为:7 7 7 7
       有另外一种情况:
       var wg sync.WaitGroup
       for i := 0; i < 5; i++ {
         wg.Add(1)
         go func() {
            fmt.Println("loop i: ",i)
            wg.Done()
         }()
         time.Sleep(1 * time.Second)
       }
       wg.Wait()
       在注释掉time.Sleep(1 * time.Second)的场景下,输出为: 5 5 5 5 5; 加上time.Sleep()后输出为: 0 1 2 3 4 5
       对于这种并发中的闭包,由于其共享变量i,且主goroutine并不会等待子goroutine,在不等待的同时,由于goroutine的启动需要时间,导致子goroutine在
       运行时,主goroutine的for循环已结束,此时i=5,则子goroutine打印5,加上sleep之后,子goroutine启动,此时共享变量i=0,1,2,3,4,则打印出来的就是
       0 1 2 3 4,如果是上述非并发情况下,如何实现打印:0 1 2 3 4,可在:nonSlice = append(xxx) 之前加一个: temp := str , 用temp执行append操作
       的参数(最终结果如上注释部分),这样在每次循环的时候,都会申请一个临时变量,闭包绑定这个临时变量,再次打印就是: 1 3 5 7
    c> 数组长度是数组类型的一部分,比如[3]int和[4]int是两种不同类型的int数组
    d> 数组可以通过指定索引和对应的值来初始化,如: var array1 := [...]int {0:1,3:4} // 声明并初始化一个长度(len(arrayA)=4)为4的int数组
       数组的长度和最后一个索引的值相关,比如: var array2 := [...]int{99:-1} // 声明一个长度为100的int数组,最后一个元素为-1,其余为0
    e> 不能对map中的元素取地址操作,如: addr = &mapVar["wu"] // compile error: cannot take address of map element 
        原因可能是map类型随着元素的增加map可能会重新分配地址,这将导致原来的地址失效
    f> map为nil时不能添加值,如: var mapVar [string]string  mapVar["wu"] = "xxx" // panic: assignment to entry in nil map
       必须使用make或者将map初始化后才能添加元素, var mapVar map[string]string   mapVar = map[string]string {"wu": "xxx"} 
    g> &var.field和(&var).field的区别,&var.field相当于&(var.field),取的是field字段的内存地址,(&var).field取的是field字段的值,另:若直接
       对结构体变量取地址操作,发现打印出来的是:&{field1_val,field2_val,field3_val,...}
    h> go中函数返回值类型的时候,不能赋值,如下代码所示:
       type Employee struct { ID int Name string Addr string} 
       func EmployeeById(id int) Employee { return Employee{ID:id} }
       func main() {
           EmployeeById(1).Addr = "cdtf"
           // var xiaoming = EmployeeById(1)
           // var.Addr = "cdtf"
       }
       上述输出报错:cannot assign to EmployeeById(1).Addr,函数EmployeeById返回的是值类型(右值?),值类型不能被赋值,正确的做法是声明变量接收
       函数返回的值类型,而变量可以被赋值
   i> 在声明方法时,如果一个类型名称本身就是一个指针,则不允许出现在方法的接收器中,示例如下:
      type employee * Employee
      func(this employee) changeName(name string) { this.Name = name }
      编译报错:invalid reciever type employee(employee is a pointer type), Go中规定只有类型(Type)和指向类型的指针(* Type)才是合法的接收器
      为了避免歧义,如果一个类型本身就是一个指针的话,不允许出现在接收器声明
   j> 值类型变量赋值nil时,编译报错:cannot use nil as type Employee assignment
      对nil指针赋值,运行时报错:panic runtime error: invalid memory address or nil pointer dereference
   k> Go的时间格式化方法比较少见,与一般用YMD hms不同,示例如下:
      time := time.Now()
      time.Format("20060102") // YMD
      time.Format("2006-01-02") // Y-M-D
      time.Format("2006-01-02 15:04:05") // Y-M-D H-M-S
      time.Format("2006-01-02 00:00:00") // Y-M-D 00:00:00
      "2006-01-02 15:04:05"这个时间点,据说是Golang的诞生时间
   l> 给一个nil chan发送or接收数据,报:fatal error, all goroutines are asleep - deadlock,所以var a_chan chan type之后,必须调用make()函数
      对chan变量进行初始化,给一个已关闭的chan发送数据(在close(chan_var)之后再次发送数据),会报:panic, send on closed channel,而从一个已关闭
      的chan接收数据,不会报错,返回的是一个chan类型的零值
17. Go Command
    a> go build命令用于编译指定的源码或代码包及它们的依赖,如果执行go build命令时不跟代码包,则命令将会尝试编译当前目录所对应的代码包
       go中的源码文件有三大类,命令源码文件,库源码文件和测试源码文件,命令源码文件作为可执行的程序入口,库源码文件一般用于集中放置各种待被使用的
       程序实体(全局变量,全局常量,接口,结构体,函数等),测试源码文件则对上述两种文件的功能和性能进行测试
       go build命令在编译一个或多个包含库源码文件的代码包时,只会做检查性的编译,不会输出任何结果文件
       go build命令既不能同时编译多个含有命令源码文件的代码包,也不能同时编译多个命令源码文件,因为如果把多个命令源码文件看成一个整体,那么多个
       文件中的main函数就属于重名函数,编译器会抛出重复定义错误
       go build常见参数如下:
         -o xxx 指定编译输出文件名称
         -i 使得go build命令安装编译目标依赖的且还未被安装的代码包,默认放在当前工作区目录的pkg目录下的响应子目录中
         -a 强行对所有涉及到的代码包(包括标准库中的代码包)进行重新构建,即便已经是最新的
         -n 打印编译期间所用到的命令,但是并不真正执行它们,可用于模拟执行
         -x 打印编译期间用到的命令,但是会真正执行
         -w 打印编译时生成的临时工作目录,并在编译结束时保留,默认情况下,编译结束后会被删除
         -v 打印被编译的代码包的名称
         -buildmode=xxx 用于指定编译模式,可以控制编译器在编译完成后生成静态链接库(.a),动态链接库(.so),可执行文件(windows为.exe)
         -compiler=xxx 用于指定当前使用的编译器名称,其值可以是gc或gccgo,gc为go语言自带的编译器,gccgo为GCC提供的go语言编译器
         -gccgoflags 用于指定需要传递给gccgo编译器或链接器的标记列表
         -gcflags 用于指定需要传递给go tool compile命令的标记列表
         -ldflags 用于指定需要传递给go tool link命令的标记列表
         -linkshared 此标记与-buildmod=shared一同使用,后者会使得编译目标的非main代码包都被合并到一个动态链接库文件中,而前者会在此之上
         进行链接操作
         -pkgdir 指定一个目录,编译器只会从该目录中加载代码包的归档文件(.a),并会把编译可能生成的归档文件放在该目录下
         -tags 指定实际编译期间,需要受理的编译标签(编译约束)的列表,一般会作为源码文件开始处注释的一部分,如:
         // +build darwin dragonfly freebsd linux nacl netbsd ...
         -toolexec 此标记让我们自定义在编译期间使用一些go工具(vet,asm)的方式
    b> go install命令编译并安装指定的代码包及其依赖包,当指定的代码包的依赖还没有被编译和安装时,会先去处理依赖包,与go build命令一样,go install
       命令的代码包参数应该以导入路径的形式提供($GOROOT/src,$GOPATH/xxx_project/src),传给go build命令的大多数标记也可以传递给go install, go
       install只比go build多做了一件事:安装编译后的结果文件到指定目录,假如我们有以下项目目录结构:
       $GOPATH/home/go_proj/thd: ebin/ pkg/ src/ ../src/logging  ../src/helper ../src/utils
       如果go install命令后面跟的代码包中只包含库源码文件,那么go install命令会把编译后的文件保存在源码文件所在工作区的pkg/下,如果我们cd到
       thd/src/utils目录下,执行:go install -v -work, 那么将编译安装当前代码包,我们会在thd/pkg/下发现一个新的子目录,由当前平台的计算架构和操作
       系统命名的子目录,比如linux_386/,linux_386/目录下是代码包utils的归档文件utils.a,我们可以在utils/下编译安装logging代码包,以目录相对路径
       的方式(本地代码包路径),如: thd/src/utils/$ go install -a -v -work ../logging, 本地代码包路径以"./" or "../"的形式开始,本地代码包路径
       不支持绝对路径,go install命令会把标准库中的代码包的归档文件(.a)放到GO安装目录的pkg子目录中,而把指定代码包的第三方依赖的代码包的归档文件
       放置到当前项目的pkg子目录下,实现标准库代码包的归档文件和用户代码包的归档文件,以及不同项目的代码包的归档文件之间的分离
    c> go get命令可以根据实际情况从互联网上下载或更新指定的代码包及其依赖包,并对它们进行编译和安装,go get命令所做的动作也被称为代码包远程导入,
       传递给该命令的参数也被称为代码包远程导入路径,go get不仅可以从github上下载代码包,也可以从任何命令支持的某个或某些代码版本控制系统检出
       代码包,go get所做的就是从版本系统控制的远程仓库中检出代码包并对其编译安装
    d> go run命令可以编译并运行命令源码文件(main),它包含了编译动作所以也适用go build命令的标记,go run命令只接受源码文件作为参数,而不接收代码包
       也不接受测试源码文件,如果命令源码文件可以接受参数,可以在go run命令的文件参数后面跟上参数,如:
       thd/src/logging/$ go run main.go -pra_name xxx, 参数名以"-"开头,紧跟着的是参数名"pra_name",空格及参数值,如果需要多个参数,依次类推
    e> go test命令会自动测试每一个指定的代码包,前提是指定的代码包中存在测试源码文件,测试源码文件是以"_ test.go"为后缀的,内含若干测试函数的源码
       文件,这种测试以代码包为单位,go test在执行完代码包中的测试文件之后,会以代码包为单位打印出测试概要信息(测试时不同代码包以空格分隔),第一列
       为测试结果(是否通过),第三列为测试用时(以秒为单位)
18. tips about Beego框架之bee(类似erlang的rebar),bee工具是一个为了协助快速开发beego项目而创建的项目,通过bee可以快速创建项目,实现热编译,开发
    测试以及项目打包/发布,通过以下命令创建bee工具:go get github.com/beego/bee, 安装后默认在$GOPATH/bin,常用命令如下:
    a> new命令,新建一个web项目,在$GOPATH/src路径下,shell/cmd命令行执行:bee new thd 创建一个名为thd的新项目
    b> api命令,new新建web项目,api命令创建API应用,该命令还支持一些自定义参数自动连接数据库,创建相关的model和controller,如:
       bee api thd_api -table="xxx" -driver=mysql/tidb... -conn=root:xx@tcp(db_ip:port)/table_name
    c> run命令,监控beego项目,通过fsnotify监控文件系统,可以实时看到项目修改之后的效果,可以实时监控项目中controller相关的.go文件的更改,自动build
       及restart整个项目,只需刷新浏览器即可看到最新的效果,无需再手动编译运行
    d> test命令,基于go test封装的一个命令,执行beego项目test/下的测试用例
    e> pack命令,用来发布应用的时候打包,会把项目打包成zip文件
    f> generate命令自动生成代码,参数如下:
       bee generate model [model_name] [-fields="xxx"] 基于fields一键生成RESTful model, -fields是一系列table字段,
          格式为:field:type,如: -fields="id:int,name:string,title:string",-fields参数下同
       bee generate controller [controllerfile] 生成RESTful controller
       bee generate view [viewpath] 在指定viewpath生成CRUD view
       bee genearte scaffold [scaffold_name] [-fields="xxx"] [-driver=mysql] [-conn="root:@tcp(ip:port)/table_name] 
       bee generate migration [migrationfile] [-fields="xxx"] 生成数据库架构更新迁移文件
       bee generate test [routerfile] 生成router文件的测试用例
       bee generate appcode [-tables="xxx"] [-driver=mysql] [-conn="root:@tcp(ip:port)/table_name] [-level=3] 
          基于已创建的database自动生成appcode,-tables是一列以逗号分隔的表名,默认为空,表示database中的所有表, -driver [mysql|postgres|sqllist]
          默认是mysql, -conn driver使用的连接数据库信息,-level 1|2|3 1-models 2-models,controllers 3-models,controllers,router
    g> migrate命令,用于项目的数据库迁移,每次项目升级,降级的sql管理
       bee migrate [-driver=mysql] [-conn="xxx"] 运行all migrations
       bee migrate rollback [-driver=xxx] [-conn="xxx] 回滚上一次migration操作
       bee migrate reset [-driver=xxx] [-conn="xxx"] 回滚all migration操作
       bee migrate refresh [-driver=xx] [-conn="xxx"] 回滚all migration 并再次重新运行
    bee安装目录下有一个bee.json文件,为bee工具的配置文件,有几个常见的配置项
       "watch_ext" : [] 可用于监控其它类型的文件,默认只监控.go的文件
       "cmd_args" : [] 如果需要在每次启动时加入启动参数,可使用该配置项
       "envs" : [] 如果需要在每次启动时设置临时环境变量参数,可使用该配置项
