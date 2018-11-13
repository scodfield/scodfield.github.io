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
