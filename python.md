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
   m> 假设我们要增加函数的功能,比如在函数调用前后打印调用日志,但又不希望改动函数原有的定义,这种在代码运行期间动态增加功能的方式,
      称之为装饰器Decorator,本质上Decorator就是一个返回函数的高阶函数,示例:
        def log(func):
            def wrapper(*args,**kw):
                print('call %s():' % func.__name__)   // __name__ 函数对象的属性,可以得到函数的名字
                return func(*args,**kw)
            return wrapper
      log()函数就是一个Decorator,它接受一个函数,并返回一个函数,我们需要借助Python的@语法,将Decorator置于函数的定义出
        @log
        def now():
            print('2019-8-1')
      调用now()函数时,不仅会运行now本身,还会在运行now函数前打印一行调用日志:
      >>> now()   // call now():  2019-8-1
      把@log放到now()函数的定义处,相当于执行了:now = log(now), 由于log是装饰器,所以它返回一个函数,因此原来的now()函数还在,新的now()函数
      指向返回的新函数,于是调用now()将执行新函数,也就是返回的wrapper()函数,wrapper函数参数为(*args,**kw),它可以接受任意参数的调用,在
      wrapper函数内,首先打印日志,紧接着调用原始的now函数
   n> functools.partial帮助我们创建一个偏函数,示例:
        >>> import functools
        >>> int2 = functools.partial(int,base=2)
        >>> int2('1000000')  // 64
      functools.partial的作用是,把一个函数的某些参数固定住(设置默认值),返回一个新函数,方便新函数的创建和调用
4. some tips to remeber:
   a> Python以下划线开头的标识符有特殊意义,以单下划线开头的表示不能直接访问的类属性(私有属性,只能在模块内引用),需通过类提供的接口访问,也不能
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
   g> python的包就是包含各个.py模块的文件夹,通过包来组织模块,避免模块冲突,要注意的是每一个包目录下面都一个__init__.py的文件,这个文件是必须
      存在的否则Python会把这个目录当成普通目录,而不是一个包,__init__.py可以是空文件,也可以有Python代码,__init__.py本身就是一个模块,模块名
      就是包名,类似的,也可以有多级目录,组成多层次的包结构
      注:当我们在命令行运行模块时,Python解释器会把一个特殊变量'__name__'置为'__main__',常见用法如下:
          if __name__ == '__main__':
              test_code()
      如果在其它地方导入该模块时,if 判断将失败,因此上述用法可以让模块通过命令行运行时,执行一些代码,常用于运行测试
   python进阶: https://docs.pythontab.com/interpy/
5. asyncio编程模型就是一个消息循环,类似于erlang中的gen_server:loop函数,不断的循环处理可以执行的coroutine
   具体说就是从asyncio模块获取一个EventLoop的引用,然后把需要执行的coroutine扔到EventLoop循环中,这样就可以实现异步IO
   python的协程是由yield定义的generator实现的,@asyncio.coroutine可以把一个generator标记为coroutine类型,使用示例:
   import threading,asyncio
   @asyncio.coroutine
   def hello():
       print("hello, %s" % threading.currentThread())
       yield from asyncio.sleep(1)
       print("hello again, %s" % threading.currentThread())
   loop = asyncio.get_event_loop()  // 获取EventLoop
   tasks = [hello(), hello()]
   loop.run_until_complete(asyncio.wait(tasks)) // 等待所有的task执行完毕
   loop.close()
   从上述演示代码可知,通过@asyncio.coroutine把一个generator标记为coroutine,然后再coroutine内部通过yield from来调用另一个coroutine实现
   异步操作,asyncio是python3.4引入的对异步IO的支持,为了更好的标识异步IO,python3.5引入了新的语法async/await,可以让coroutine更简洁易读,
   使用新语法,只需要做两个简单的替换即可: @asyncio.coroutine替换为async; yield from 替换为await; 新语法代码如下示:
   async def hello():
       print("hello, %s" % threading.currentThread())
       await asyncio.sleep(1)
       print("hello again, %s" % threading.currentThread())
   注:yield from 语句可以方便的调用另一个generator,并拿到返回的值
6. tips about Flask:
   a> Flask()构造函数使用当前模块(__name__)的名称作为参数,创建一个Flask类对象,一个对象就是一个WSGI应用程序
      from flask import Flask
      app = Flask(__name__)
   b> flask类的route()函数是一个装饰器,告诉应用程序,哪个url由那个函数处理,示例:
      @app.route('/hello')
      def hello():
          return 'Hello Flask!'
      application对象的add_url_rule()方法也可以用于将url与函数绑定,如:app.add_url_rule('/','hello',hello)
   c> Flask类的run()方法,在本地开发服务器上运行应用程序,app.run(host,port,debug,options),所有参数均可选
      host 默认为127.0.0.1,设置为'0.0.0.0'可使服务器外部可用
      port 默认为5000
      debug 默认为False,设置为True,则提供调试信息
      options 要转发到底层的Werkzeug服务器
   d> 通过向url rule参数添加变量部分,可用动态构建url,变量标记为:<var_name>,变量名作为关键字参数传递给与url规则绑定的函数,示例:
      @app.route('/hello/<name>')
      def hello(name):
          return 'Hello %s' % name
      如果在浏览器中输入: http://localhost:5000/hello/thd, 则'thd'将作为参数传递给hell()函数,除了默认字符串变量外,还可以添加
      int,float,path等类型变量,示例:
      app.route('/thd/ver/11')   app.route('/thd/sub_ver/1.1')  app.route('/thd/flask/')
      定义'/thd/flask/'和'/thd/flask', 这两个url规则有些不同,第一个规则是规范的,使用'/thd/flask/'和'/thd/flask'返回相同的输出,但是
      如果定义的是第二个规则,使用'/thd/flask/'访问,会返回404错误
   e> url_for()函数用于动态构建特定函数的url,需要从flask导入redirect和url_for,示例:
      @app.route('/hello_admin')
      def hello_admin(): 
          return 'Hello admin'
      @app.route('/guest/')
      def hello_guest(name):
          return 'Hello guest, %s!' % name
      @app.route('/user/')
      def hello_user(user):
          if user == 'admin':
              return redirect(url_for('hello_admin'))
          else:
              return redirect(url_for('hello_guest',name = user))
      http://localhost:5000/user/admin // Hello admin
      http://localhost:5000/user/thd // Hello guest, thd!
      hello_user()函数接受来自url的参数的值,赋给user变量,通过与'admin'等比较,使用url_for()函数重定向到hello_amin或hello_guest函数
   f> Flask HTTP方法:GET 以未加密方式将数据发送到服务器; HEAD 与GET方法相同,但没有响应体; POST 将html表单数据发送到服务器; PUT 用上传
      的内容,替换目标资源的所有当前表示; DELETE 删除由url给出的目标资源的所有当前表示
      Flask路由默认响应GET请求,但是可以通过为route()方法提供参数来更改默认选项,示例:
      <form action = "/login", method = "post"> xxx </form> 
      @app.route('/regist', methods = ['GET','POST'])
      def regist():
          if request.method == 'POST':
              name = request.form['name']
              pwd  = request.form['pwd']
          else:
              name = request.args.get('name')
              pwd  = request.args.get('pwd')
          return redirect(url_for('login',name=name,pwd=pwd))
   g> Flask模板,绑定url的函数可以以html的形式返回,但是从Python代码生成具体甚至是负责的html内容将会非常麻烦,可以利用Flask基于的jinjia2模板
      引擎,返回html,需导入render_template(),示例:
      from flask import render_template
      @app.route('/hello')
      def hello():
          return render_template('hello.html')
      Flask将尝试从应用程序的templates/文件夹中找到hello.html
      web模板系统(web templating system)指的是设计一系列html脚本,其中可以动态插入变量数据,模板系统包括模板引擎,数据源和模板处理器
      jinja2引擎使用以下分隔符从html转义:{% ... %} 用于if,for等语句,结束语句则是:{% endif %}, {% endfor %}; {{ ... }} 用于表达式输出到
      模板,比如:{{ var_name }}, 用于显示在render_template()函数中指定的关键字参数; {# ... #} 未包含在模板输出中的注释; # ... # 行语句
   h> Flask静态文件,在应用程序的static/文件夹下,包括支持html显示的css,js文件以及图片等,html中的script标签指定js脚本文件,示例:
      <script type = "text/javascript" src = " {{ url_for('static', filename='xxx.js') }}"></script>
   i> Flask Request对象,来自客户端网页的数据作为全局请求对象(上述示例中的request变量)发送到服务器,Request对象包含以下属性: form 一个字典
      对象,包含表单参数及其值; args 解析查询字符串,url中'?'之后的那部分; cookies 字典对象,包含cookies的键值; files 与上传文件有关的数据;
      method 当前请求方法
   j> cookie以文本文件的形式存储在客户端的计算机上,目的是记住和跟踪与客户使用相关的数据,Flask Request对象包含cookie属性,如果对Respond
      响应对象设置cookie,需要先通过resp = make_response('resp_xxx.html')函数的返回值来获取响应对象,再通过响应对象的set_cookie函数来
      存储cookie,resp.set_cookie('userID',1001),最后再:return resp,返回响应,上述不带cookie的返回是直接:return(render_template('xx'))
      与cookie不同,session(会话)数据保存在服务器上,会话是客户端登录服务器,并注销的时间间隔,需要在会话期间保存的数据存储在服务器上的临时目录
      为每个客户端的会话分配一个会话ID,会话数据存储在cookie的顶部,服务器以加密方式对其进行签名,对于此加密,Flask需要定义一个SECRET_KEY,session
      对象也是一个全局的字典对象,保存会话变量及其值,比如,提示玩家登录:
      @app.route('/')
      def index():
          if 'username' in session:
              username = session['username']
              return render_template('succ_login.html',username=username)
          return render_template('login.html')
      当然我们需要在'/login'规则绑定的login()函数中,设置该session变量: session['username'] = request.form['username']
      以及在'/logou'规则绑定的logout()函数中,删除该session变量: session.pop('username',None)
      实现上述session的前提是包含以下操作:
      from flask import session
      app.secret_key = 'my_random_key'
   k> Flask消息闪现,类似桌面系统的消息框,或js使用的警报,Flask模块提供flash(message, category)方法,它将消息传递给下一个请求,通常是一个模板
      message参数是要闪现的消息,category参数可选,可以是'error','info'或'warning',为了从会话中删除消息,模板调用get_flashed_messages()
      以下示例在模板中接收消息:
      {% with messages = get_flashed_messages() %}
        {% if messages %}
          {% for message in messages %}
            {{ message }}
          {% endfor %}
        {% endif %}
      {% endwith %}
   l> Flask处理文件上传非常简单,只需要一个html表单,并将其enctype属性设置为'multipart/form-data',url规则对应的处理函数从request.files[]
      对象中提取文件,并保存到指定位置,可在Flask应用的配置文件中设置默认的上传路径及文件的最大大小(字节),字段为'UPLOAD_PATH','MAX_CONTENT_PATH'
      比如,form表单中的<input type='file' name='up_file' />, 则提取文件为: file = request.files['up_file']
