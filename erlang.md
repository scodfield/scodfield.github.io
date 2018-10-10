1. centos6.9安装erlang 20.3
   安装依赖: yum install -y gcc gcc-c++ glibc-devel make automake autoconf ncurses-devel openssl-devel m4 kernel-devel
   下载源码: wget http://erlang.org/download/otp_src_20.3.tar.gz
   解压缩: tar -zxf otp_src_20.3.tar.gz
   编译安装: 
     ./otp_build autoconf
     ./configure && make && sudo make install
   erlang默认安装在/usr/local/(bin,lib/erlang)
   ./configure --prefix 指定安装路径,如 --prefix /opt/erlang/20.3
   ./configure --without-javac 不适用java编译器
2. sublime3依据文件、文件夹创建项目
3. 项目脚本gencmd,protoc,package等,shift命令左移参数,xcopy复制文件夹的文件及子文件夹(/e参数复制所有子文件夹,/s复制非空子文件夹)
4. 项目管理rebar,编译前协议tag与id(crc32循环冗余校验)映射,生成error_msg
5. game_server.spec文件必须为标准json格式,之前第一行出现多余的‘,’,导致报错,另外,xls2erl被编译成了exe文件
6. erlang escript, 第一行#!xxx 指定解释器,若以escript xxx调用escript脚本,则第一行不会起作用,
   第二行%% -*- erlang -*- 可选指令,使Emacs编辑器进入erlang源码模式
   escript脚本可以通过%%! xxx 模拟器传递参数 -smp -sname -mnesia等等
7. 关于配置文件,erlang的配置文件是为app服务的,是app的环境变量,也即xx.app文件中的env选项
   erlang有三种配置,优先级依次如下:
   a> erl -AppName1 Par1 Val1 -AppName2 Par2 Val2
   b> erl -config xxx.config 该配置文件中罗列各app的配置项, [{AppName1, [{Par1,Vla1}]}, {AppName2, [{Par2, Val2}]}].
   c> xxx.app 中的env选项
   项目中用到的其它所谓的“配置”,均由py脚本生成.config文件,再在服务启动的时候,编译生成对应的.beam文件
8. 配置导表时,若报错KeyError则输入的config name在game_server.spec里没找到json对应的key值
   若报错UTF BOM则是在修改game_server.spec文件后保存的时候默认的编码格式不是UTF-8(记事本默认编码为ANSI)
   sublime text3 -> preferences -> settings -> 在左边栏找到"default_encoding":"xxx",将该行复制到右边栏的user
   文件,值改为"UTF-8",sublime text3默认不显示当前文件的编码格式,可在左边栏找到"show_encoding"和"show_line_endings" 复制到右边栏user文件中,
   将值改为true即可,保存关闭,可在sublime的最右下角看到编码格式
9. rebar2 get-deps的时候,重定向deps的下载路径,rebar_deps:download_source/2,  
   第一个参数AppDir是deps下载后的保存路径
10.rebar2 create-app appid=AppName 在当前目录下创建app,生成src/xx.app.src,xx_app.erl,xx_sup.erl
