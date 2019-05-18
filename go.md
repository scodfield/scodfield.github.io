1. 新的数据类型interface,chan,chan的运算符'<-'
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
