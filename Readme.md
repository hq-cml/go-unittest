#Golang单元测试

Golang的单测用到的框架比较多，各有各的用途和场景， 梳理汇总一下常见的单测框架及使用方法。

##云时代： 
* 容器化
* CI/CD

##名词： 
* mock
* 打桩
* 表驱动

##常用框架： 
###Convey：
 * GoConvey是一个单元测试框架，可以将单测case进行整合，框架兼容Golang原生的单元测试，所以可以使用go test -v来运行测试。
 * 因为Goland是可以自动生成单测代码的，所以个人感觉用处不是特别的大

###GoStub：
主要是为了打桩：为一个全局变量、函数等等打桩

###GoMock：
GoMock是由Golang官方开发维护的测试框架，主要用于interface的Mock功能，实现依赖注入之类

###Monkey：
Monkey是Golang的一个猴子补丁（monkeypatching）框架  
原理是在运行时通过汇编语句重写可执行文件，将待打桩函数或方法的实现跳转到桩实现，原理和热补丁类似

###GoMonkey：
可以认为是monkey的升级版，支持函数、方法、全局变量打桩（关闭内联-gcflags=all=-l）

###GoSqlMock：
sqlmock包，用于单测中mock db操作

###Httptest：
Golang官方自带，于生成一个模拟的http server.   
主要使用的单测场景是：已经约定了接口，但是服务端还没实现

###Httpexpect：
严格意义来说，感觉他并不算一个单测的框架  
主要用能与场景：  
主要用于http server已经完成实现，就绪启动的场景下，构造各类参数，对http server的回值进行检验探测