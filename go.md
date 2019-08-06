1. 新的数据类型interface,channel,channel的运算符'<-',defer,recover与panic
   a> go的数据类型中,只有切片slice,字典map和信道channel是引用类型,其它都是值类型
   b> defer延迟调用函数,含有defer语句的函数,在函数返回前,调用另一个函数,main() {defer hello() xxx },在main()函数返回前调用hello()函数,不仅适用于
   函数,还可以延迟调用(结构体)方法,defer语句的实参取值是在执行defer语句的时候,而非在调用延迟函数的时候,main() { a := 5 defer printA(a)
   a = 10 fmt.Println("before defer a:" a)} func printA(a int) { fmt.Println("in defer a:",a) },上述延迟调用结果:"in defer a 5"
   defer栈,当在一个函数内多次调用defer语句时,Go会把defer调用放入一个栈中,再按照LIFO的顺序执行(可以实现字符串逆序输出)
   当一个函数应该在与当前代码流无关的环境下调用时,可以使用defer,比如用到sync.WaitGroup的地方,在协程函数内声明一个: defer wg.Done()
   c> mutex 是sync包中的一个结构体类型,它主要定义了Lock()和Unlock()这两个方法,用于在Go中提供了一种处理竟态条件(race confition)的加锁机制(locking 
   mechanism),可确保在某时刻只有一个协程在临界区(critical section)运行,防止出现竟态条件,如果有协程持有了锁(Lock()),当其它协程试图获得该锁时,这些
   协程会被阻塞,知道mutex变量解锁锁定(Unlock())为止,另外一种实现加锁的方法是使用非缓冲信道,在操作临界区前发送数据,操作后接收数据即可
   d> func panic(interface{})会终止程序的执行,并在延迟defer函数(如果定义了延迟函数)执行完之后,程序控制返回执行该panic()函数调用的调用方,退出过程
   一直进行直到当前协程的所有函数都退出,然后打印panic信息,接着打印堆栈跟踪(stack trace),发生panic时,recover可重新获得对该程序的控制,
   panic-recover与其它语言的try-catch-finally类似,调用panic的两种场景:发生不能恢复的错误,程序无法在继续运行下去(如,端口被其它程序占用,导致
   当前web应用绑定端口失败);发生了一个编程上的错误(如,一个非空的参数,在实际调用时为nil),上述退出过程会一直持续到main()(主协程)退出为止,跟踪堆栈
   则是从当前调用panic的协程函数开始直到main(),运行时错误(如,数组越界)也会导致panic,等价于调用了内置函数panic,其参数由接口类型runtime.Error给出
   e> func recover()interface{} 也是内建函数,用于重新获得panic协程的控制,只有在延迟函数内部,调用recover才有用,在延迟函数内调用recover,可以取到
   panic的错误信息,停止panic续发事件(panicking sequence),程序恢复正常运行,如果在defer外部调用recover,不能停止panic续发事件,只有在同一个go协程
   中调用recover才行,recover不能恢复一个不同协程的panic,recover恢复panic之后,堆栈跟踪就被释放,但是可以通过runtime/debug.PrintStack()函数打印
   堆栈跟踪,导入debug包,并在defer的recover函数里调用即可:import "runtime/debug" defer recov()   func recov(){if r := recover(); r != nil
     fmt.Println("recovered: ", r) debug.PrintStack()}
   在Printf()函数中,使用'%T'格式说明符,可以打印出变量的类型,unsafe包提供了一个Sizeof()函数,该函数接收变量并返回它的字节大小,不过unsafe包可能会带来
   可移植性问题,因此需要小心使用
   go有着非常严格的强类型特征,没有自动类型提升/转换,比如在C语言中,整型可以和浮点型变量相加,但go中会报错(invalid operation),操作符两侧的变量类型
   必须一致,类型转换为T(v),其中v为变量,T为系统变量类型(int,string,float),类型转换举例:i := int(56.78)
   go中字符串是字节的切片,'+'操作符可用于拼接字符串,go中的字符串兼容unicode编码,并使用utf-8编码
   rune是go内建类型,int32的别称,rune表示一个代码点,代码点无论占用多少个字节,都可以用一个rune来表示
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
    
