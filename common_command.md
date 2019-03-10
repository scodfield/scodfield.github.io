1. awk '{pattern + action}' 
   {} 非必需
   pattern 前后带斜杠的正则表达式（/^[abc]/）
   action 匹配后执行的命令
2. telnet (选项) (参数)
   telnet命令用于登录远程主机,采用明文传输,为应用层协议
   选项如下:
   -8: 允许使用8位字符资料,包括输入与输出
   -a: 尝试自动登入远端系统
   -b<主机别名>: 使用别名指定远端主机名称
   -c: 不读取用户专属目录里的.telnetrc文件
   -d: 启动排错模式
   -e<脱离字符>: 设置脱离字符
   -E: 滤除脱离字符
   -f: 此参数的效果和指定"-F"参数相同
   -F: 使用Kerberos V5认证时,加上此参数可把本地主机的认证数据上传到远端主机
   -k<域名>: 使用Kerberos认证时,加上此参数让远端主机采用指定的领域名,而非该主机的域名
   -K: 不自动登入远端主机
   -l<用户名称>: 指定要登入远端主机的用户名称
   -L: 允许输出8位字符资料
   -n<记录文件>: 指定文件记录相关信息
   -r: 使用类似rlogin指令的用户界面
   -S<服务类型>: 设置telnet连线所需的ip TOS信息
   -x: 假设主机有支持数据加密的功能,就使用它
   -X<认证形态>: 关闭指定的认证形态
   参数如下:
   远程主机: 指定要登录进行管理的远程主机
   端口: 指定TELNET协议使用的端口号
   示例:
      telnet 192.168.1.7
      telnet 192.168.1.8 8080 (可测试端口是否开启)
3. 检测端口是否开启的方式:
   telnet ip port
   ssh -v -p port username@ip (ssh -v -p 6789 root@192.168.1.7)
   wget ip:port
   连接成功一般都会有connected提示,失败一般是提示connection refused
4. cat /etc/redhat-release 查看centos发行版本
5. windows上进行压测的时候,统计已占用的tcp接口: netstat -ano | find "TCP" /c, /c 统计个数
6. chmod 管理文件或目录的访问权限,两种用法,一种是由字母和操作符组成的文字设定法,另一种是数字设定法
   文件或目录权限包括只读、只写和可执行,用户类型包括文件拥有者、同组用户、其它用户
   ls -l显示文件详细信息时,最左边的第一列显示了拥有者、同组用户、系统其它用户的权限(rwx),一共10个位置,第1个位置指定了文件类型
   d 表示一个目录, - 表示一个非目录文件,从第2到第10共9个位置,每3个一组,分别对应上述三组用户的权限
   用户标识: u 文件或目录所属用户; g 用户所在群组; o 除用户及其群组外的其它用户; a 用户、群组和其它用户
   权限标识: r 读权限,数字4表示; w 写权限,数字2表示; x 执行权限,数字1表示; - 不具有任何权限,数字0表示
   文字设定操作符: + 表示增加权限; - 表示取消权限; = 表示赋予限定的权限,并取消其它剩余的权限
7. svn commit的时候,提示:changing file 'xxxxx' is forbidden for the server,查了一下说是用户没有相应权限
8. netstat 查看网络状态,参数如下
   -a 显示所有选项
   -t 显示tcp协议的连接状态
   -u 显示udp协议的连接状态
   -n 不显示列名
   -l 显示处于listen状态的连接
   -s 按各协议(ip,icmp,tcp,udp等)显示统计信息
   -p 显示建立连接的程序名(PID/program_name)
   -r 显示本地路由信息
9. shell编程反引号表示执行命令
   if判断,[] 与变量之间加空格
   可由$dir1/$dir2 or $dir1"/"$dir2 拼接路径
   变量赋值没有空格,引用变量'$'
   shift左移变量,每执行一次变量个数('$#')减一,第一个变量被销毁,第二个变量变为第一个, shift n 左移n位(销毁前n位变量)
10. shell编程中,运算符和参数直接要用空格隔开,通配符(星)在用于乘法运算符时,要用反斜杠,或单双引号修饰
    expr 用于整数值计算,如: res=expr 3 + 4
    let 与expr类似,let可同时有多个表达式,运算符和参数之间也无需加空格,如: let res=8/3 res2=5+6, <>&|等特殊符号同样需要反斜杠,单双引号
11. gcc常用编译参数:
    -o 输出到指定文件
    -E 仅作预处理,不进行编译,汇编和链接
    -S 仅编译到汇编语言,不进行汇编和链接
    -c 编译,汇编到目标代码,不进行链接
    -Idir 将dir加入搜索头文件的路径列表中
    -Ldir 将dir加入搜索库文件的路径列表中
    -std=xxx 指定使用某个标准编译程序,例: -std=c99 or -std=gnu99
    -shared 指定生成动态链接库
    -fPIC 表示编译为位置独立的代码,不用此选项的话编译后的代码是位置相关的,所以动态载入时,是通过代码拷贝的方式来满足不同进程的需要
          而不能达到真正代码段共享的目的
    例: gcc -fPIC -shared -o test.so test.c -I /usr/local/lib
12. 查看系统内核版本的几个方式:
    uname -a or -r
    cat /proc/version
    查看发行版本的方式:
    lsb_release -a
    cat /etc/redhat-release
    cat /etc/issue
13. zip/unzip
    zip -r xxx.zip yyy 将当前路径下的yyy目录压缩为xxx.zip
    zip -r xxx.zip yyy zzz.txt 压缩yyy目录和zzz.txt到xxx.zip
    unzip xxx.zip 将xxx.zip解压到当前路径
    unzip xxx.zip -d yyy 将xxx.zip解压到yyy目录
14. 设置系统环境变量
    export PATH=xxx/yyy:$PATH 只对当前shell有效,退出则失效,为临时设置
    编辑/etc/profile文件,添加:export PATH=$PATH:xxx/yyy 保存退出后,下次进入shell生效,或者执行:source /etc/profile,对所有用户永久生效
    编辑~/.bash_profile文件,后续操作同/etc/profile,对当前用户永久生效
    source /etc/profile的原理是再一次执行/etc/profile shell脚本,使用sh /etc/profile是不行的,因为sh是在子shell进程中执行的,即使环境改变了
    也不会反应到当前环境中,而source是在当前shell进程中执行的,所以能看到环境的改变
    环境设置文件包括:系统环境设置文件和个人环境设置文件
    系统环境设置文件包括:登录环境设置文件/etc/profile,非登录环境设置文件/etc/bashrc
    个人环境设置文件包括:登录环境设置文件~/.bash_profile,非登录环境设置文件~/.bashrc
    登录环境指的是用户登录系统后的工作环境
    非登录环境指的是用户调用子shell时使用的工作环境
    env命令显示所有的环境变量
15. source,sh,bash,./
    source xxx/yyy 在当前shell读取并执行文件中的命令,文件无需可执行权限,source命令可简写为".", source xxx/yyy 等同于 . xxx/yyy
    sh/bash xxx/yyy 使用sh/bash解释器在子shell中读取并执行命令,文件无需可执行权限
    ./xxx 在子shell中读取并执行命令,需要文件有可执行权限
16. vi编辑模式下,按Esc进入命令行模式,"u"是撤销,相当windows下的Ctrl+Z,撤销了多次,ctrl+r(重做)来反转撤销的的动作,也即它是撤销的撤销
    命令行模式,移动光标到某一个字符上,按"x",删除一个字符
    命令行模式,移动光标到某一行,按"dd",删除一整行
17. 测试服升级了配置,想查看一下CPU及内存相关信息
    cat /proc/cpuinfo | grep "physical id" | sort | uniq | wc -l 查看物理cpu个数
    cat /proc/cpuinfo | grep "cpu cores" | uniq 查看每个物理cpu中的core数
    cat /proc/cpuinfo | grep "processor" | wc -l 查看逻辑cpu个数,逻辑cpu个数 = 物理cpu数 X 每个物理cpu的core数 X 超线程数
    cat /proc/cpuinfo | grep name | cut -f2 -d: | uniq -c 查看cpu型号
    查看内存的使用情况一般用free命令,free的输出可以看做是一个二维数据,包含了3行6列
    6列字段值分别如下
    total 总内存大小, used 实际已使用内存, free 空闲内存, shared 进程共享的内存总内存, buffers 缓冲区内存, cached 缓存内存
    3行字段如下
    Mem 操作系统角度的内存使用情况, -/+ buffers/cached 应用程序角度, Swap 交换区角度
    free命令常用的几个参数
    -b 以byte为单位显示内存使用情况, -k -m -g 以K/M/GB为单位显示, -o 不显示缓冲区调节列, -s<间隔秒数> 持续观察内存使用情况
    cat /proc/meminfo 同样可查看内存使用信息
18. top命令显示系统正在运行的进程及相关资源的使用情况,是常用的性能分析工具
    top命令前5行是系统整体情况的统计
    第1行是系统运行时间和平均负载,包括当前时间,系统已运行时间,当前登录用户数量及过去三个统计周期内的平均负载
    第2行是任务(进程)信息,包括总进程,及处于running,sleeping,stopped,zombie等状态的进程数量
    第3行是CPU状态信息,包括运行用户(us)/内核(sy)进程的cpu时间百分比,运行已调整(ni)优先级的用户进程的cpu时间百分比,
    cpu空闲(id)时间百分比,等待IO完成(wa)的cpu时间百分比,处理硬件(hi)/软件(si)中断的cpu时间百分比,虚拟机被hypervisor偷去的cpu时间百分比
    第4/5行是物理内存/交换分区的内存使用信息,包括内存总量,已使用,空闲,缓冲和缓存等,与free命令类似
    系统信息下是各进程的统计信息,包括以下字段
    PID 进程ID, USER 进程所有者, PR 进程调度优先级, NI 进程的nice值(优先级),值越小,优先级越高,可为负值, VIRT 进程使用的虚拟内存(SWAP+RES),
    RES 进程使用的,未被换出的物理内存大小,也即驻留内存, SHR 进程使用的共享内存大小, S 进程状态 D=不可中断的睡眠状态 R=运行中 S=睡眠中 
    T=被跟踪或已停止 Z=僵尸进程, %CPU 自上次更新到现在的cpu占用时间百分比, %MEM 进程使用的内存百分比, TIME+ 进程启动后到现在使用的全部cpu时间,
    精确到百分之一秒(1/100s), COMMAND 运行进程使用的命令
    top命令默认以%CPU的降序排列,可通过:shift + '<' 或者 '>' 向左、向右改变排序列
    top命令执行中可输入交互命令:M 按驻留内存(RES)排序, P 按cpu使用百分比排序, T 按时间/累计时间排序
19. vi查找,在命令行模式下输入"/"进入搜索,输入要匹配的字符串,回车
    n 继续查找下一个匹配的关键字, N 在当前位置,反向查找匹配的字符串
    vi中有很多设置项,比如显示行号,在命令行模式,输入: set number 显示行号, set nonumber 关闭行号, set nu/nonu 是上述两个命令的简写
    比如在长文本中查找时,现在命令行模式下打开行号,再输入"/xxx",即可清楚的看到匹配的字符串在哪一行
20. 项目打包,之前用的是WinRAR,后边有同学想直接用zip命令,在此记录一下
    https://sourceforge.net/projects/gnuwin32/files/zip/3.0/zip-3.0-bin.zip/download 解压缩,取zip.exe
    https://sourceforge.net/projects/gnuwin32/files/bzip2/1.0.5/bzip2-1.0.5-bin.zip/download 解压缩,取bzip2.dll
    将上述两个文件放在 X:\Program files\Git\usr\bin 
21. 为了便于回溯项目改动,TortoiseSVN可以强制在提交时写日志,项目文件夹右键 --> Properties --> New --> Log size 设定日志的大小即可
22. linux系统万物皆文件,关于程序的配置文件可参考:https://www.ibm.com/developerworks/cn/linux/management/configuration/index.html
23. 阮老师的科普文用来重温基础知识真真是极好的,这次是make命令:http://www.ruanyifeng.com/blog/2015/02/make.html
    安装使用GNU的autoconf,automake生成的程序,最常用的就是下面三个命令: ./configure,make,make install
    ./configure命令检测平台的目标特征,生成Makefile文件,./configure命令可加上参数,对安装的程序进行控制,./configure --help参考参数及其含义
    make命令读取Makefile文件中的第一个目标,读取目标下的commands,编译程序
    make install读取Makefile中的install目标,将程序安装到制定的位置
    make命令接收目标参数读取Makefile文件,或:-f xxx or --file=xxx制定的规则文件,Makefile文件由一系列规则构成,规则格式如下
    targe:prerequisites
    <tab> commands
    规则描述了目标,目标的依赖文件及如何构建目标的命令,make命令在执行时先扫描依赖文件,若依赖文件不存在或者last_modification时间戳比target
    的时间戳新,则会重新构架目标文件,如果依赖文件时间戳都比目标文件的时间戳晚(依赖文件自上次构建之后,再未变化),则不会再次构建目标
    commands构建命令,是一行或多行shell命令,每一行commands都在不同的子shell中执行,所以并不能跨行使用在其它子shell中声明的变量or环境变量
    解决方法可以放到同一行,用分号隔开,或者在换行符前加反斜杠转移,再或者是使用.ONESHELL内置变量
    make命令默认不编译上次编译之后就没有更改的文件,若要更改该默认行为,可使用-B选项,例: make -B xxx
24. 子项目增多,每次启动都要在多个文件夹中来回切,搞个启动脚本,记录一下在shell脚本中调用另一个脚本的方式,资料来自网络
    fork:调用方式为sh path/to/script.sh or ./script.sh,fork在执行时,新开子shell执行脚本,子shell的环境变量继承自父shell,执行完毕返回
    父shell,子shell的环境变量不会带回父shell
    exec:被调用脚本在当前shell执行,使用exec调用脚本后,父脚本中exec之后的内容不再执行
    source:被调用脚本在当前shell执行,当前脚本可使用被调用脚本声明的变量和环境变量,相当于将多个脚本合并在一起执行
25. 之前都是查看文件后几行,现在有一个erl_crash.dump文件,比较大(4.15G),编辑器打开比较慢,在gitbash里面用head命令查看前几行即可
    head -n N xxx 或者 head -N xxx,例:head -n 5 erl_crash.dump, head -5 erl_crash.dump
26. 临时搭建一个http服务器
    yum install httpd 安装
    service httpd start 启动
    配置文件 /etc/httpd/conf/httpd.conf,常见修改的选项包括:ServerName,DocumentRoot,Listen
   去掉apache的欢迎页,/etc/httpd/conf.d/welcome.conf文件删掉或改为其它名字,此时将会显示/var/www/html/的目录结构
   更改后,重启服务,service httpd restart
27. 利用scp命令来通过ssh上传下载文件
    从远程服务器下载文件or文件夹:scp [-r] user_name@server_addr:/path/to/file_or_dir /path/to/local_file_or_dir
    上传文件or文件夹到远程服务器:scp [-r] /path/to/local_file_or_dir user_name@server_add:/path/to/remote_dir
28. 查找文件所有安装路径:whereis xxx (whereis erl)
    查找运行文件所在路径:which xxx (which erl)
    查找某一个具体文件:find /target/path -name xxx.y (find / -name erl_nif.h)
29. 为了方便测试,搞个虚拟机,安装及配置网络桥接模式参考以下文章
    https://blog.csdn.net/collection4u/article/details/14127671
    https://blog.csdn.net/fuguangruomeng/article/details/79244055
    官网下载VMware,再百度一下VMware的注册机或者激活码
    其次是要明白WMware三种网络模式的区别,目前的需求是必须得用桥接模式,配置桥接模式的时候
    最开始没找到网上说的/etc/sysconfig/network-scripts/ifcfg-eth0,最后发现是修改一下目录下的ifcfg-ensxx即可(xx依据版本不同)
    新增的几个配置项主要是IPADDR,NETMASK,GATEWAY,DNS1,DNS2 (MTU,NM_CONTROLLED可改可不改), 更改了ONBOOT=yes
    保存退出之后,service nework restart 重启网络服务
    安装VMware Tools需要手动挂载镜像:https://www.cnblogs.com/liwanliangblog/p/9193880.html
30. VMware安装后,启动节点,然而并不能访问,本地和VMware能ping通,telent server_ip port_number不同,熟悉的一幕又来了,防火墙....
    centos7默认是firewall防火墙,查看防火墙状态:firewall-cmd --state  // running
    关闭防火墙:systemctl stop firewalld.service
    禁止开机启动:systemctl disable firewalld.service
    重启防火墙:systemctl restart firewalld.service, firewall-cmd reload
    查看已开发的端口:firewall-cmd --list-ports
    开启端口:firewall-cmd --zone=public --add-port=80/tcp --permanent 
    参数说明:zone 作用域, add-port 添加端口,格式为端口/通信协议, permanent 永久生效,没有此参数重启后会失效
31. centos7安装svn
    yum -y install subversion
    svn help [sub_command] 查看checkout,update(up),revert等命令的用法
    通过svn命令行设置文件or文件夹的externals对应的subcommand是propxxx,可查看详细的用法
32. what happens when start a process on linux?
    首先进程有很多属性,包括:打开的文件或网络连接,环境变量,信号处理器,内存,寄存器,命名空间,当前工作路径等
    首先fork(),克隆当前进程,fork()返回-1为报错,0为子进程ID
    其次execve(),改变内存,寄存器,和执行的程序,不过环境变量,信号处理器等都不变
    关于fork()函数的内存copy,实际上是"copy on write",只有当父子任一进程有写内存操作时,才会发生实际的copy操作
33. lsof(list open files) 列出当前系统打开的文件,输出各列的含义如下:
    command 进程名称
    pid 进程标识
    user 进程所有者
    fd 文件描述符
    type 文件类型,如DIR(目录),CHR(字符类型),IPv4/6(ip协议套接字),UNIX(unix域套接字)
    device 磁盘名称
    size 文件大小
    node 索引节点(文件在磁盘上的标识)
    name 打开文件的确切名称
    lsof的-i选项可以列出符合条件的进程情况,语法:lsof -i[46][protocol][@hostname|hostaddr][:service|port]
    46 --> IPv4 or IPv6
    protocol --> tcp or udp
    hostname --> 网络主机名
    hostaddr --> IPv4地址
    service --> /etc/services中的service-name,可以有多个,以逗号分隔,如:lsof -i:rje,echo
    port --> 端口号,可以有多个,同样以逗号分隔,如:lsof -i:80,8088
34. PAM(Pluggable Authentication Modules) 可插拔认证模块
    PAM机制采用模块化设计和插件功能,使用户可以轻易地在应用程序中插入新的认证模块或替换原先的组件,同时不必对应用程序做任何修改
    PAM为了实现插件化和易用性,采用了分层设计思想,让各认证鉴别模块和应用程序独立,然后通过PAM API作为二者联系的纽带
    PAM的体系: 应用程序 <---> PAM API(PAM配置文件) <---> 各PAM模块
    pam.limit.so模块, 主要功能是限制用户会话过程中对各种系统资源的使用情况,该模块的配置文件默认在/etc/security/limits.conf
    配置文件由4个字段组成:用户/用户组; 类型(soft/hard); 资源(resource); 值
    可配置资源如下:
    core -- 内核文件(core file)大小(KB); data -- 最大数据大小(KB); fsize -- 最大文件大小(KB); memlock -- 最大锁定内存地址(KB);
    nofile -- 最大可打开文件句柄数; rss -- 最大驻留空间(KB); stack -- 最大堆栈空间(KB); CPU -- 最大cpu使用时间(Min);
    nproc -- 最大运行进程数; as -- 地址空间限制(KB); maxlogins -- 该用户可最多登录系统次数; maxsyslogins -- 最多可登录系统次数;
    priority -- 用户进程优先级; locks -- 用户最大锁定文件数; sigpending -- 最大挂起信号数量; nice -- 最大nice值,默认为[-20,19]
    rtprio -- 最大实时优先级
    比如针对某用户的最大文件数,先修改/etc/security/limits.conf xx soft/hard/- nofile 65535 
    配置文件修改后,修改/etc/pam.d/login文件,添加: session required /lib64/securiy/pam_limits.so
    可通过ulimit命令做进一步调整,之后重启机器
    注: nofile 是基于用户层面的限制, 系统层面的限制需要修改/etc/sysctl.conf fs.max-file参数,调整之后, systcl -p 使之生效
35. 内核根据进程的nice值决定进程需要多少处理器时间,nice值的取值范围是[-20,19],-20优先级最高,19最低,ps axl 命令可查看进程nice值(NI字段)
    nice -n adjustment -adjustment --adjustment[=]Value command/process 调整/指定应用程序的优先级
    renice [-n] priority [-p|--pid] pid 调整正在运行的进程的优先级
    进程优先级越高,获得的cpu时间越多
36. file命令用来辨识给定文件的类型,file命令检验文件类型按以下顺序来完成:
    检验文件系统中支持的文件类型; 检验magic file文件规则; 检验文件内容的语言和字符集
    使用格式:file [-bcLvz][-f <文件名称>][-m <魔法数字文件>][文件或目录],参数如下:
    -b 列出辨识结果时,不显示文件名
    -c 详细显示指令执行过程,便于排错或分析程序执行的情形
    -f <文件名称> 指定文件名称,有多个文件名称时,file依序辨识这些文件
    -i 显示MIME类型
    -L 直接显示符号连接所指向的文件类别
    -m <魔法数字文件> 指定魔法数字文件
    -v 显示版本信息
    -z 尝试解读压缩文件的内容
    [文件或文件列表] 要确定类型的文件列表,多个文件直接使用空格分开,可以使用shell通配符匹配多个文件
    注:魔法数字文件(magicfile),以.mgc为扩展名的文件
