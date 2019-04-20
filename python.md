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
3. 单引号'',双引号"",三引号'''...''',"""...""" 均可表示字符串,使用时没有任何区别
   不过单双引号字符串通常都写在一行,如需换行,需在行尾加换行符"\",同时单双引号相互嵌套时,需在转义字符
   而三引号字符串可以由多行组成,无需显式使用换行符,同时三引号可以字符串内可以直接使用单双引号,而无需转义
   最后,三引号内还可以包含注释
4. Python以下划线开头的标识符有特殊意义,以单下划线开头的标识不能直接访问的类属性,需通过类提供的接口访问,也不能通过from xximport * 导入;
   以双下划线开头的代表类的私有成员;以双下划线开头和结尾的是Python里特殊方法专用的标识,比如__init__() 标识类的构造函数
