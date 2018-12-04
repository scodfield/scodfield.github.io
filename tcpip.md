Windows下几个影响TCP并发连接的参数：
0. Ctrl+R --> regedit --> HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\services\Tcpip\Parameters
1. TcpNumConnections
   maximum number of connections of that TCP can have open simultaneously
   0 or 未设置该参数表示可以建立任意的连接(不过其它条件,如tcb,可用端口数or系统内存会限制连接上限)
2. MaxUserPort
   highest available user port when requests from the system
   typically, ephemeral ports(those used briefly) are allocated to port numbers 1024 through 5000
   服务器监听固定的端口,所以并不会存在端口不够用的情况,而客户端建立socket连接时,默认并未bind本地端口,因此需要向
   系统申请一个可用的端口,当并发连接足够多时,会出现端口不够用的场景
3. MaxFreeTcbs
   TCP control blocks(TCBs), each connection requires a control block
   TCB是一种内存驻留结构,包含套接字编码,传入和传出数据缓冲位置,已收到或未确认的字节及其它信息
   为了快速检索这些信息, windows server将TCB存储在一个哈希表中
4. MaxHashTableSize
   TCP control blocks(TCB) stored in a hash table, the value must be power of 2
   
Linux下TCP并发连接数量限制：
1. 本地端口范围
2. 单进程同时打开的文件句柄上限,默认为1024,可通过ulimit -n xxx 命令进行调整
3. 系统同时打开的文件句柄上限,可通过/etc/sysctl.conf调整
4. 系统内核的IP_TABLES防火墙对最大跟踪的TCP连接数有限制

Tips to remember:
1. TCP报文有两部分组成：头部和数据部分,头部的各字段体现了TCP的使用和功能,通信的另一端收到tcp报文后,去掉头部,组装接收到的消息
   头部的前20个字节是固定的,这也是头部的最小大小,头部最大60个字节,结构如下:
   a> 源端口 2byte
   b> 目的端口 2btye
      源端口告知数据来自于哪个应用程序,目的端口指明数据发送给哪个应用程序,应用程序绑定端口,结构大致如下：
      程序1 <----> Sock1(源&目的ip,源&目的port) <----> Sock2(源&目的ip,源&目的port) <----> 程序2
   c> 序列号 seq number 4byte
      tcp面向字节流,序列号就是所发送报文的首字节在字节流中的位置或序号
   d> 确认号 ack number 4byte
      确认号是接收方希望接收到的下一个报文的首字节的序号
      tcp双向通信,所以在三次握手建立连接的时候,双方同步seq number,更新各自的ack number
      这里想到个问题,在第三次握手的时候,客户端发送ACK包,其中 ACKbit=1,ack_number=服务端seq number+1,那此时客户端的seq number=??
      又或者,服务端ack number已在第一次握手时同步了客户端seq number,在连接建立之前(此时为syn_recv状态)or接收到客户端ACK包时
      并不会更改自身ack number,也即第三次握手时,客户端无需关心seq number
   e> Data offset数据偏移(首部长度) 4bit
      数据偏移指的是报文的数据部分起始处到报文起始处的距离(刚好度量了首部的长度),一个偏移量是4byte,4bit最大是15,所以TCP头部
      的最大值是60byte
   f> 保留位 6bit 
      以后使用,目前为0
   g> 标志位 6bit
      6bit标志位,每一个bit都有特定用途：
         1> 紧急 URG(Urgent) URG=1表示紧急指针(urgent pointer)有效,告诉系统该报文有紧急数据,需尽快发送,与下面的紧急指针配合使用
         2> 确认 ACK(Acknowledgment) ACK=1时确认号(ack number)有效,ACK=1的报文一般成为'确认报文'(这解决了第三次握手时客户端序列号
            的问题,第三次握手客户端发送的是'确认报文',此时ACK=1,SYN=0),连接建立后发送报文时,ACK=1
         3> 推送 PSH(Push) PSH=1表示该报文是高优先级,接收方(内核)应尽快将该报文交付给后续的应用程序,而无需等整个TCP缓冲满后再交付
         4> 复位 RST(Reset) RST=1表示tcp连接出现严重错误,需释放并重新连接,一般称该报文为'复位报文'
         5> 同步 SYN(Synchronization) SYN=1表示该报文是一个请求连接的报文,一般称为'同步报文',三次握手的第一个报文
         6> 终止 FIN(Finish) FIN=1表示该报文发送方的发送数据已发送完毕,请求释放连接,一般称为'结束报文',四次挥手断开连接时用到此标志
   h> 窗口大小 window size 2btye
      告诉对方本段tcp接收缓冲区还可以接收多少字节的数据,对方据此控制发送的速度
      窗口大小指的是从本报文段头部ack number起,可接收的字节数
   i> 校验和 TCP Checksum 2byte
      校验和由发送端填充,接收端对整个报文进行CRC计算,目的是检验报文段在传输过程中是否损坏,损坏则丢弃
   j> 紧急指针 Urgent Pointer 2byte
      仅在标志位URG=1时有效,指定了本报文段中紧急数据的字节数
      紧急数据在报文段数据部分的首部,其后是普通数据,即当URG=1时,TCP报文:头部+数据(紧急+普通)
   
2. TCP连接建立（三次握手）由内核协议栈实现,连接建立后socket状态转为established,并被放入icsk_accept_queue，accept()被唤醒,返回socket
3. listen()开启监听队列,客户端SYN包到来,创建新sock,sock为状态TCP_SYN_RECV,并被存入半连接队列syn_table中
4. SYN攻击:客户端伪造大量IP地址,不间断的向服务器发送SYN包,塞满服务端半连接队列,导致正常的SYN请求被丢弃,SYN攻击是DDos攻击的一种,检测SYN攻击
   netstat + awk '/^tcp/' 查看SYN_RECV状态的tcp连接即可
5. 在进行本地压测的时候,出现一个情形,首次启动客户机没有任何问题,加机器人再次启动客户机时,日志报:connrefused
   第一反应是上次的客户机tcp连接没有释放,由于客户机tcp连接设置了保活机制(keepalive),而默认的保活检测时间过长(默认是7200s),以及保活重发次数
   另外一个会影响的参数是系统释放连接资源之前的等待时间(TcpTimedWaitDelay),系统不释放资源导致本地动态端口,TCB,TCB hashtable等系统资源不足,
   从而拒绝连接,注册表调整参数
   a. Ctrl+R --> regedit --> HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\services\Tcpip\Parameters 
   b. 右键 --> 新建 --> DWORD --> 输入键值 --> 选中键值&右键 --> 修改 --> 基数栏选择十进制&填入具体参数值
   c. MaxUserPort 65534
   d. MaxFreeTcbs 16000
   e. MaxHashTableSize 65536
   f. TcpTimedWaitDelay 30 (s)
   g. KeepAliveTime 120000 (ms)
   h. KeepAliveInterval 1000 (系统未收到响应而重发保活信号的间隔,ms)
   调整参数,关闭注册表,重启电脑即可
6. windows下查看tcp、udp及端口等统计情况:netstat -an | find "ESTABLISHED" /c 统计活跃状态的tcp连接,状态与linux类似,包括LISTENING,CLOSE_WAIT,
   ESTABLISHED,TIME_WAIT
7. netstat参数 
   -a 显示所有连接和监听端口
   -n 以数字形式显示地址和端口号
   -o 显示与每个连接相关的所属进程
   -p proto 显示proto指定的协议连接,包括TCP,UDP,TCPv6,UDPv6
   比如显示所有tcp连接: netstat -an -p tcp
8. windows下与linux类似的为find,常用参数如下: 
   /v 指定不包含指定所有行
   /c 对指定的行进行技术
   /i 不区分大小写
   "" 指定要搜索的字符串
   比如显示端口为20001的所有连接: netstat -ano | find "20001" ,类似的还有一个findstr,不需要(""),如: netstat -ano | findstr 20001
9. HTTP/1.0 短连接,每次request都会建立一个单独的连接,因此请求较多时,连接的建立和释放会占用大量的系统资源
   HTTP/1.1 支持长连接,管线处理,在一个连接上可以传送多个请求和响应,并且客户端可以在上一个请求未返回前再次发送请求,不过服务器则需要保证
   按客户端请求的顺序,返回响应,1.1还新增了一些请求/响应头域来扩展功能,比如:status code, request method(options,put,delete..),host域
10. 服务器在httpc:request/4时报错:inet eaddrinuse, 同时:erl -sname test  也会报同样的错误,POSIX Error Codes说是address被占用
    通过对比其他资料,这个address应该就是需要占用的端口,httpc申请的是动态端口,那么很有可能是动态端口不够用
    /proc/sys/net/ipv4/ip_local_port_range 查看本地TCP/UDP端口范围
    可由: echo '32768 60999' > /proc/sys/net/ipv4/ip_local_port_range  更改端口范围
11. TCP/UDP端口类型:
    0 - 1023 固定端口,与常见应用紧密绑定的端口
    1024 - 49151 注册端口,与应用松散绑定
    49452 - 65535 动态端口,可与任意应用绑定
12. 做个试验,在浏览器直接输入:server_ip:port/anything, 查看日志
    11:03:34.564 gw进程init ok
    11:03:34.564 gw进程do_terminate,报head_error,这个是因为在解析http请求时候,<<Len:32>> = <<"GET ">>,解析协议头部发送的body长度时
    Len解析出来的是1195725856,过滤异常包,直接断掉连接即可
    11:03:34.567 gw进程init ok (浏览器在上一次断开后,紧接着就再一次的发起连接请求)
    11:03:35.597 gw进程报相同的错误
    11:05:15.616 11:05:15.617 发起两次请求,#Port<0.7882> #Port<0.7883>,有意思的是7882有打印报错信息,7883没有打印任何信息,且与11:10
    发起的连接相比,#Port<>不在连续
    11:10:16.597 11:10:16.598 11:10:18.629 发起三次请求,#Port<0.9078> #Port<0.9079> #Port<0.9080>,三次请求连续,且均有打印错误信息
    后续直到11:37,浏览器未再发起连接请求,试验用的是chrome浏览器,莫非chrome浏览器在实现上会多次发起请求,待解
