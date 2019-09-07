1. pyinstaller用来打包Python应用程序,打包完的程序可以在没有安装Python解释权的机器上运行,支持的操作系统包括:windows,linux,mac os
   solaris,AIX等,常用参数如下:
   -F 生成单个可执行文件(同--onefile)
   -D /path/to/dist 指定生成可执行文件的路径,默认在当前路径,一般会生成两个文件夹build,dist,可执行文件在dist中
   -p 添加自定义搜索路径,一般用不上
   --add-data [src;dest or src:dest]指定资源文件,前面是源,后面是目的,中间是分隔符,windows分号,linux冒号,可多次使用添加多个资源文件
      例:--add-data xxx:yyy --add-data xyz:zyx,该参数的使用效果和下述spec文件中的Analysis.datas相同 
   -i 指定可执行文件图标
   -c 使用控制台,无窗口
   -w 使用窗口,无控制台
   pyinstaller分析Python程序,找到所有的依赖,然后将Python解释器和依赖项放在一个文件夹下或一个可执行文件中
   pyinstaller默认打包成一个文件夹,文件夹下包括依赖项及可执行文件,打包成文件夹的好处是在debug阶段即可看到依赖项有没有被包含进来,
   打包成文件夹时的工作   流程,pyinstaller的引导程序是一个二进制可执行程序,启动程序的时候,引导程序开始运行,首先创建一个临时的Python环境,
   然后通过Python解释器导入依赖,运行
   pyinstaller可通过-F参数将所有文件打包到一个可执行文件,当程序运行时,引导程序创建一个临时文件夹,解压缩依赖到临时文件夹,后续执行流程同上
   pyinstaller默认在当前路径生成xxx.spec文件,spyinstaller通过执行spec文件中的内容来生成app,有点像makefile,一般无需管spec文件,
   不过以下情况需要修改spec文件:需要打包资源文件,需要include一些pyinstaller不知道的run-time库,为可执行文件添加run-time选项,多程序打包
   可通过:pyi-mkspec options xxx.py [yyy.py ...] 命令生成spec文件
   spec文件主要有4个class:
   Analysis 以本地py文件为输入,分析py文件的依赖项,并生成相应的信息
   PYZ 是一个.pyz的压缩包,包含程序运行所需的所有依赖项
   EXE 根据上面两项生成
   COLLECT 生成其它部分的输出文件夹,也可以没有
   默认生成的spec文件不满足需求,最常见的情况就是我们的程序依赖一些本地文件,此时就需要编辑spec文件来添加本地数据文件,上面Analysis中的datas就是
   要添加到项目中的数据文件,datas是一个列表,每个元素是一个二元组,二元组的第一个元素是本地文件索引,第二个元素是copy到项目中之后文件的名字
   例: a = Analysis(... datas=[('/path/to/local/file', 'name_in_project'), ...] ...)
   一般情况下,pyinstaller [option] xxx.py, 如果需要修改spec文件,可以通过上述命令先生成spec文件
   修改之后,运行:pyinstaller [option] xxx.spec 即可
   'pyinstaller [option] xxx.py'命令的执行过程也是先生成spec文件,再按spec文件进行打包,上述先生成spec文件只不过是将这个过程拆开了而已
2. 安装pyInstaller时,下载某些依赖时提示time out,只好直接下载whl文件
   whl文件的安装方法很简单,下载之后,在whl文件目录开启gitbash,执行:pip install xxx.whl即可
   自动安装到python安装目录的Lib/site-packages/xxx
3. tips about syntax:
   a> python中的字符串类型是str,在内存中以unicode表示,一个字符对应若干个字节,但是如果要在网络中传输或者保存到磁盘上,就需要把str变为以字节为
      单位的bytes(类似于go中的字节数组?),bytes类型的数据以带'b'前缀的单引号或双引号表示(x=b'ABC'),要注意区分'ABC'和b'ABC',两者在显示上是一
      样的,但是bytes每个字符只占用一个字节,在内存中以unicode表示的str可以通过encode()方法,编码为指定的bytes,如:'ABC'.encode('ascii')  
      '中文'.encode('utf-8'), 反过来,从磁盘或网络读取到的是bytes字节流,Python通过decode()方法把bytes转换为str,如:b'ABC'.decode('ascii')
      如果bytes中含有无法解码的字节,decode会报错,如:b'\xe4\xb8\xad\xff'.decode('utf-8')  // UnicodeDecodeError:cannot decode '0xff' in 
      position 3 invalid start byte, 如果bytes中只有一小部分的无效字节,可以通过传入errors='ignore'忽略错误的字节,如: 
      b'\xe4\xb8\xad\xff'.decode('utf-8',errors='ignore') // '中' 
      注:len()函数计算的是str的字符个数,如需要计算字节,转为bytes后再调用len()就是计算的字节数
   b> python的两种有序列表类型列表list&元组tuple,list可通过+/-1/2等下标直接访问or更改or赋值指定下标的元素,还可以通过var_list.append(xx)
      追加元素到末尾,通过var_list.pop()方法删除末尾元素,或var_list.pop(index)删除指定下标的元素,通过var_list.insert(i,xx)把元素插到指定位置
      tuple则不同,一旦初始化就不能再更改,它没有append(),insert(),pop()等方法,也不能用下标来赋值,所以在定义的时候tuple的元素就必须确定下来,空
      tuple是(),如果tuple只有一个元素var_t1 = (1),则会发现var_t1为1,这是因为()还用来表示数学公式中的小括号,这就产生歧义了,因此Python规定,定义
      只有一个元素的tuple时,带上一个逗号,以消除歧义,如: var_t1 = (1,),  tuple不可变,但是它可以包含其他类型的变量(比如list),如:var_t2 = (1,2,
      [23,'2er3','g34',67]),这就定义了一个"可变的"tuple,这个时候变的只是tuple中可变类型变量的值,tuple中的每个元素的指向是不变的,指向了一个list,
      就不能改成指向其它的对象,但是指向的这个list本身是可以变的
   c> python有两种循环方式,一种是for xx in xxx,依次把list/tuple中的每个元素迭代出来,另一种是while循环,条件满足就不断循环,不满足时退出循环
      break语句主动退出循环,continue语句结束当前循环,直接开始下一次循环
   d> python内置了字典类型dict,其它语言也有类似的结构,如erlang中的maps,go中的map,dict使用键-值存储,访问速度快
      列表list是中括号[],元组tuple是小括号(),字典dict是大括号{},如: var_dict1 = {"name":"thd","score":95}
      dict也是可变类型,可通过key直接赋值,如:d["df"] = 83, 访问一个不存在的key会报错,避免访问不存在的key,有两种方法,其一是in关键字判断key是否
      存在,如: 'thd2' in d1 // False, 其二是通过dict提供的get()方法,如果key不存在返回None或者指定的值,这种方法和erlang中的maps:get/3类似
      示例: d1.get('thd2',-1)
      要删除dict中的某个key,调用pop()方法,如: d1.pop('thd')
      注: dict中的key必须是不可变类型,比如字符串和整数等,list可变,不能作为key,key不可变是主要还是因为要对key执行hash算法,计算value的存放地址
   e> set与dict类似,不过它保存的是一组key的集合,并没有value,且set中key不能重复,创建一个set,需要一个list作为输入集合,如:var_set = set([1,2,3])
      通过var_set.add(key4)方法添加元素到set,允许重复添加相同的key,但不会有任何效果,通过var_set.delete(key4)方法删除指定key
      set可以看做是数学意义上的无序且无重复元素的集合,因此可以对两个set做并集(|),交集(&)操作,示例如下:
        s1 = set([1,2,3])
        s2 = set([3,4,5])
        s1 & s2 // [3]
        s1 | s2 // [1,2,3,4,5]
      set和dict的唯一区别就是没有存储对应的value,其它均一样,也就是说set同样不可以存放可变对象
      注:关于不可变对象,先看一个示例:
        a = 'abc'
        b = a.replace('a','A')
        a // abc
        b // Abc
        注意区分变量和字符串对象,变量a在调用replace方法的时候,方法是作用在字符串对象'abc'上的,字符串对象不可变,所以replace方法并不是修改
        字符串对象'abc'的内容,而是创建一个新的字符串对象'Abc',并返回,所以对于不变对象来说,调用自身任意的方法都不会改变自身的内容,而是创建新
        的对象并返回,这样就保证了不可变对象本身是永远不可变的
   f> 单引号'',双引号"",三引号'''...''',"""...""" 均可表示字符串,使用时没有任何区别
      不过单双引号字符串通常都写在一行,如需换行,需在行尾加换行符"\",同时单双引号相互嵌套时,需在转义字符
      而三引号字符串可以由多行组成,无需显式使用换行符,同时三引号可以字符串内可以直接使用单双引号,而无需转义
      最后,三引号内还可以包含注释
   g> 递归函数定义简单,逻辑清晰,但是要防止出现栈溢出,由于函数调用是借助栈帧这种数据结构实现的,每次调用都会增加一个栈帧,返回时再减一个,所以在栈
      大小一定的情况下,递归调用的次数过多时,会导致栈溢出,解决栈溢出的方法是尾递归优化,尾递归指的是函数返回(return)时调用自身,它要求返回语句不能
      包含表达式,这样编译器或者解释器就可以做尾递归优化,无论递归调用多少次,都只占用一个栈帧,不会出现栈溢出的情况,示例如下:
      def fact_no_wei(n):
          if n == 1:
              return 1
          return n * fact_no_wei(n-1)
      // 尾递归优化
      def fact(n):
          return fact_tail_recursion(n,1)
      def fact_tail_recursion(n,res):
          if n == 1:
              return res
          return fact_tail_recursion(n-1,n * res)
      可以看到: return fact_tail_recursion(n-1,n * res)只返回了递归函数本身, n-1 & n * res会在函数调用前就计算好,不影响函数调用
      注:上述代码的优化并不一定能实现尾递归调用,最终还是要看该语言的编译器or解释器是否针对尾递归做优化,python标准解释器貌似没有做相关的优化,
      不过erlang虚拟机已对尾递归进行了优化,所以erlang的尾递归不会出现栈溢出
  h> Python中的切片操作下标是左闭右开,切片的[start:end:step],第一个':'前后是起始&结束下标,第二个':'的右边是切片时的步长
     list,tuple,str都可以应用切片操作,操作结果还是对应的数据类型
  i> 任何可迭代对象都可以用for...in来迭代
     对于dict来说,默认迭代的是key, for k in var_dict: , 如果迭代values,可以用for v in var_dict.values(): ,如果同时迭代key和values,
     可以用for k,v in var_dict.items(): ,
     那么如何判断一个对象是不是可迭代对象,可以通过collections模块的Iterable类型来判断,示例:
     from collections import Iterable
     isinstance('abc',Iterable) // True
     isinstance([2,3,4,6],Iterable) // True
     isinstance(100,Iterable) // False
     注:如果想对list进行下标循环,Python内置的enumerate()函数把list变成索引-元素对,这样就可以在for循环找中同时迭代索引和元素本身,示例:
        for i,v in enumerate(['a','b','c','d']):
            print(i,v)
  j> Python的列表生成式语法,[var_express for var in var_list],生成元素的表达式放到最前面,后面跟for...in循环,for循环后面还可以加上if判断
     [var_express for var in var_list if var_bool_express], 还可以使用双层循环,构造全排列,示例:
     [a+b for a in 'abc' for b in 'efg']
     erlang也有列表生成式,语法稍有不同([express || var <- var_list,var_bool_express])
  k> 通过列表生成式可以很方便的生成一个list,但是创建一个很大的列表,不仅浪费空间,有时候也并不会访问全部的元素,因此,如果列表的元素可以按照
     某种算法推算出来,那么我们就可以在循环中不断推算出后续元素,这样就不必创建完整的list,从而节省大量的空间,Python中这种一边循环一边计算
     的机制,称之为生成器generator
     创建生成器的方法有多种,第一种方法最简单,只有把一个列表生成式的'[]',改成'()'即可,示例:
     L = [x * x for x in [1,2,3,4]]  // 1 4 9 16
     G = (x * x for x in [1,2,3,4])  // <generator object <genexpr> at 0x00001B5D56D3>
     生成器创建之后,可以通过next()函数获得generator的下一个返回值,每次调用next(G)就计算G的下一个元素,直到抛出StopIteration错误,由于生成器
     也是可迭代对象,所以也可以用for循环,事实上大多数情况下都是用for循环来迭代生成器
     需要注意的是,调用next()或for循环迭代生成器G之后,再次调用时next()会抛出StopIteration错误,for循环则不再打印任何元素,可知迭代结束后,生
     成器会保持结束状态
     生成器定义的另一种方法是:函数定义中包含yield关键字,那么这个函数就不再是一个普通函数,而是一个generator,示例:
       def fibo(tar):
           n,a,b = 0,0,1
           while n < tar:
               yield b
               a,b = b, a+b
               n = n + 1
           return 'done'
       f = fibo(6) // <generator object fibo at 0x00001D4C6BD5>
       next(f)  // 1
       next(f)  // 1
       next(f)  // 2
       next(f)  // 3
     generator函数每次调用next()时执行,遇到yield语句返回,再次调用next()时,从上次返回的yield处继续执行,同样的,多数情况下也不会用next()函数
     返回下一个值,而是直接使用for...in循环,for n in fibo(6): print(n)  // 1 1 2 3 5 8
   l> 可以被next()函数调用并不断返回下一个值的对象称为迭代器Iterator,通过isinstance()和collections模块的Iterator类型,可以判断一个对象是不是
      迭代器对象,生成器都是Iterator对象,list,dict,tuple,str虽然是Iterable,却不是Iterator,list,str等Iterable可以变成Iterator,通过iter()函数
      isinstance([],Iterator) // False
      isinstance(iter([]),Iterator) // True
      注:为什么list,str不是Iterator,因为Iterator对象表示的是一个数据流,一个有序序列,但是我们不能提前知道序列的长度,只能不断通过next()函数
      按需计算下一个值,Iterator的计算是惰性的,只有在需要返回下一个值时才会进行计算,所以Iterator可以看做是一个惰性计算序列
4. some tips to remeber:
   a> Python以下划线开头的标识符有特殊意义,以单下划线开头的表示不能直接访问的类属性,需通过类提供的接口访问,也不能
      通过from xximport * 导入; 以双下划线开头的代表类的私有成员;
      以双下划线开头和结尾的是Python里特殊方法专用的标识,比如__init__()表示类的构造函数
   b> 当函数的参数不确定时,可以使用* args和 ** kwargs,* args 没有key值,** kwargs有key值,主要用于函数定义,用于传递不定数量
      的参数,其实并不是必须写成args和kwargs,也可以用var和kwvar,只有* 才是必须的
      * args用来发送一个非键值对的可变数量的位置参数列表给函数,实际存储类型是一个tuple,如下示例:
      def var_args(fixed_arg,* argv):
          print("fixed arg: ", fixedZ_arg)
          for arg in argv :
              print("another args through * argv: ", arg)
      var_args("apple","banana","orange") 
      ** kwargs 允许将不定长度的键值对,作为关键字参数传递给函数,实际存储类型是一个dictionary,示例如下:
      def kw_args(** kwargs):
          for key,value in kwargs.items():
              print("{0} --> {1}".format(key,value))
      kw_args(name="thd") // name --> thd
      * args 和 ** kwargs的使用示例:
      def var_args_kwargs(arg1,arg2,arg3):
          print("arg1: ", arg1)
          print("arg3: ", arg3)
          print("arg3: ", arg3)
      // * args
      args ("f23",23,"f3r") // arg1: f23  arg2: 23 arg3:f3r
      var_args_kwargs(* args)
      // ** kwargs
      kwargs = {"arg3":"34t","arg2":"2f3","arg1":67}
      var_args_kwargs(** kwargs) // arg1: 67 arg2: 2f3 arg3: 34t
      如果想在函数中同时使用这三种参数,定义顺序为: def fun_name(stand_args, * args, ** kwargs)
   c> global语句声明全局唯一变量,一般在函数内为函数外定义的变量赋值时,需要再次声明一下,表明这个变量是在语句块以外定义的
   d> async/await是python在3.5版本中引入的关于协程的语法糖,主要用于异步编程
   python进阶: https://docs.pythontab.com/interpy/
