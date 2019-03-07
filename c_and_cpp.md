1. c++中的虚函数用来实现多态,多态简单来说就是用父类型指针指向子类对象,进而调用实际的子类成员函数,虚函数的动态绑定是一种泛型技术(试图以
   不变的代码实现不同的算法),每个包含虚函数的类都有一个虚函数表,表中主要是该类的虚函数的地址,也即虚函数表是一个函数指针数组
2. GCC和G++的区别,手册:http://gcc.gnu.org/onlinedocs/gcc/G_002b_002b-and-GCC.html#G_002b_002b-and-GCC
   GCC最开始是GNU C Compiler,是开源免费的C语言编译器,后来又加入了对C++,Pascal,Objective-C等其它语言的支持,名字就变成了 GNU Compiler Collection
   G++是GCC编译器集合中,专门用来处理C++语言的
   GCC提供了C语言预处理器:C Preprocessor 简称CPP
   传统编译器的工作原理基本上是三段式的
