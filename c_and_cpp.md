1. c++中的虚函数用来实现多态,多态简单来说就是用父类型指针指向子类对象,进而调用实际的子类成员函数,虚函数的动态绑定是一种泛型技术(试图以
   不变的代码实现不同的算法),每个包含虚函数的类都有一个虚函数表,表中主要是该类的虚函数的地址,也即虚函数表是一个函数指针数组
2. GCC和G++的区别,手册:http://gcc.gnu.org/onlinedocs/gcc/G_002b_002b-and-GCC.html#G_002b_002b-and-GCC
   GCC最开始是GNU C Compiler,是开源免费的C语言编译器,后来又加入了对C++,Pascal,Objective-C等其它语言的支持,名字就变成了 GNU Compiler Collection
   g++是GCC编译器集合中的C++编译器, gcc则是GCC中的C编译器 
   GCC提供了C语言预处理器:C Preprocessor 简称CPP
   典型的编译过程:调用预处理器,比如CPP;调用实际的编译器,比如cc,cc1;调用汇编器(assembler),比如as;调用链接器,比如ld
   编译器是可以更换的,所以gcc调用的是C Compiler,而g++调用的是C++ Compiler
   gcc与g++的区别包括:对于.c/.cpp文件,gcc分别当做C和cpp文件编译(c/c++语法强度不一样,c++更严谨一些),而g++则统一当做cpp文件处理;
   编译时g++会调用gcc,但是gcc不能自动连接c++的标准库文件(STL),所以gcc编译c++文件时,需要加上-lstdc++参数,而g++可以自动连接库文件
   传统编译器的工作原理基本上是三段式的,分为前端(Front end),优化器(Optimizer),后端(Back end),前端负责解析源代码,语义检查,生成抽象语法树(Abstract
   Sytax Tree, AST),优化器对中间代码进行优化,后端负责生成机器代码,这一过程后端会最大化利用目标机器的特殊指令,以提高代码的性能
   GCC实现了很多前端,支持多种语言,但它是一个完整的可执行文件,没有给其它语言的开发提供重用中间代码的接口,也就是说GCC太重,非模块化
   Low Level Virtual Machine(LLVM) 是一个定位较底层的虚拟机,它的有点正是GCC的缺点,用来解决编译器代码重用的问题,LLVM与其它编译器最大的区别是
   它不仅仅是Compiler Collection,也是Libraries Collection,也就是说LLVM是一个编译器,也是一个SDK
   Clang是一个支持C,C++,Objective-C和Objective-C++等编程语言的编译器前端,Clang采用LLVM作为其后端,Clang主要有C++编写,Clang相比GCC
   编译速度快,内存占用小,基于库的模块化设计,易于IDE集成,出错提示更友好,而GCC的优势在于,除了支持C系列,还支持Java,Fortran,Go等语言,其次是使用更广泛
   支持更多平台
3.  gcc常用编译参数:
    -c 只激活预处理,编译和汇编,生成.o的目标文件
    -S 只激活预处理和编译,生成.s的汇编代码
    -E 只激活预处理,不会生成文件,不过可以通过重定向输出到另一个文件,如: gcc -E fight.c > fight_pre.txt 
    -o 指定输出目标,默认为.out
    -L 指定编译时路径,参考第60条
    -fPIC 生成位置无关的代码,不适用此选项时,编译后的代码是位置相关的,所以动态载入时通过代码拷贝的方式满足不同进程的需要,不能真正达到代码共享
    -shared 生成共享目标文件(动态库)
    -share 尽量使用动态库
    -static 禁止使用动态库
    -g 指示编译器在编译的时候,产生调试信息
    -OX X = 0,1,2,3 编译器优化选项的4个级别,0 没有优化,1 为缺省, 3 优化级别最高
    -M 生成文件的关联信息,包含目标文件所有依赖
    -MM 同-M, 但它会忽略由#include<file>造成的依赖
    -MD 同-M, 生成的关联信息将会输出到.d文件里
    -MMD 同-MM, 输出到.d文件
    LD_LIBRARY_PATH 该环境变量指示了链接器可以加载的动态库的路径,有root权限时,可通过修改/etc/ld.so.conf文件,然后/sbin/ldconfig来达到指定
    动态库的的目的,如果没有root权限,可使用该环境变量
    编译执行的4个步骤:
    预处理 预处理器cpp; 
    编译 将预处理后的文件转换成汇编代码,生成.s文件 编译器egcs; 
    汇编 将汇编文件转换为目标代码(机器码),生成.o文件 汇编器as;
    链接 连接目标代码,生成可执行程序 链接器ld
4. 字节对齐
   对齐跟数据在内存中的存储位置有关,如果变量地址正好位于它长度的整数倍,那么它就是自然对齐,比如在32位cpu下的一个4byte的int型变量,
   如果它在内存中的地址为:0x00000004,字节对齐的根本原因在于cpu访问数据的效率问题,比如上述整型变量地址在0x00000002,cpu取值会访问两次内存
   第一次读取两个字节0x00000002-0x00000003的一个short,第二次读取0x00000004-0x00000005的一个short,然后组合得到一个4字节的int,而如果该变量
   的地址在0x00000003,那么cpu就需要读三次内存,第一次0x00000003的一个char,第二次0x00000004-0x00000005的一个short,第三次0x00000006的一个char
   然后再组装成int,而如果变量自然对齐,则只需要读一次内存即可
   对于标准数据类型,它的地址只需要是长度(sizeof(type_name))的整数倍就行了,对于自定义数据类型或数组,对齐规则如下:
   数组 按照存储的元素数据类型对齐即可; 联合体(union) 按其包含的长度最大的数据类型对齐; 结构体 结构体中每个数据类型都要对齐
   GCC默认是4字节对齐,定义如下结构体: struct role{ char sex; int age; char name[10];}; struct role william;
   sizeof(william) = 20byte, GCC会在sex和name后面分别补齐3byte和2byte
   我们可以使用__attribute__ 属性来自定义字节对齐方式,使用方式如下:
   struct role{ char sex; int age; char name[10];}__attribute__ ((aligned (1))); struct role william;
   上述表示对role结构体1字节对齐,sizeof(william) = 15byte
   另外一种1字节对齐方法(也就是取消字节对齐)是:__attribute__ ((packed)) 
   第三种对齐or取消对齐的方法是使用伪指令#pragma pack ([n]),含义如下:
   #pragma pack (n) 按n字节对齐
   #pragma pack () 设定对齐方式为上一次对齐方式
   #pragma pack (push [n]) 保存当前对齐方式,并设置对齐方式为n字节
   #pragma pack (pop) 恢复最近一次保存的对齐方式
   #pragma 使用方式如下:
   #pragma pack(1)
   struct role{ char sex; int age; char name[10];}william;
   #pragma pack () or #pragma pack (pop)
   struct role2{ char sex; int age; char name[10];}george;
   sizeof(william) = 15byte  sizeof(george) = 20byte
   关于__attribute__ 可以参考简书:https://www.jianshu.com/p/29eb7b5c8b2d
   gcc编译参数-fpack-struct可以在编译时指定结构体的对齐方式,不带"=n"时,默认为4byte对齐,gcc -fpack-struct=1 -o echo -c echo.c
5. c++野指针vs悬挂(悬垂)指针:
   野指针 wild pointer,是指未初始化指针变量,此时指针指向的是内存中的任意地址,直接使用会造成严重后果
   悬挂指针 dangling pointer,是指指针变量最初指向的内存已被释放,指针被free/delete掉,或者指针指向的栈上分配的临时变量已被系统回收
   注:free只是释放指针指向的内存,而不是指针,指针是一个变量,只有程序结束时,才会被销毁,释放后指针指向的内存不再有效,
   会被系统当做垃圾内存进行回收,但是并没有改变指针变量的指向,正确的做法是释放内存之后,将指针指向NULL,防止指针后边不小心又被解引用
   https://blog.csdn.net/eszrdxtfcygv/article/details/38523659
   https://blog.csdn.net/wj3319/article/details/6871957
   https://www.cnblogs.com/idorax/p/6475941.html
