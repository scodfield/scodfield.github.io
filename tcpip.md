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

tips to remember:
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
   
2. TCP连接建立（三次握手）由内核协议栈实现,连接建立后socket状态转为established,并被放入icsk_accept_queue，accept()被唤醒,返回socket
3. listen()开启监听队列,客户端SYN包到来,创建新sock,sock为状态TCP_SYN_RECV,并被存入半连接队列syn_table中
4. SYN攻击:客户端伪造大量IP地址,不间断的向服务器发送SYN包,塞满服务端半连接队列,导致正常的SYN请求被丢弃,SYN攻击是DDos攻击的一种,检测SYN攻击
   netstat + awk '/^tcp/' 查看SYN_RECV状态的tcp连接即可