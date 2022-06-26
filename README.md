# 公会仓库管理系统
## 技术栈介绍
 该系统采用golang作为后端开发语言，golang作为前端的简易客户端的实现语言（仅仅是一个命令行），以及使用mongodb作为数据库并且使用docker作为系统的运行环境
## 设计思路
- 由于题目要求的仅仅是个人对公会数据的修改，公会成员之间并无交互，所以这里我采用的方案是：主协程监听连接+新协程处理新连接 的方式，如下图所示
![图片](/images/20220626-151157.jpg)
- mongodb使用两个collection，user和alliance，分别用来存储用户信息和公会信息，后端服务会在内存中维护一个已连接mongo的全局handler用来提供与mongo的交互
- 客户端不同的命令会分发到不同的分支里做相应的处理，然后把操作结果返回到客户端，客户端也可以输入指定的指令查看公会的数据
## 运行步骤
最好是在有VPN可以直连海外的情况下进行：
1.在命令行输入以下命令拉取项目文件
```bash
git clone https://github.com/ShaneChan/alliance.git
```
2.输入以下命令生成docker镜像文件
```bash
cd ./alliance
docker-compose build
```
此时会生成三个镜像文件，名字分别为mongo、alliance_alliance_server和alliance_alliance_client，如下图所示
![图片](/images/20220626-154836.jpg)
3.命令行开启三个标签，按顺序启动三个容器，具体顺序为mongo->alliance_alliance_server->alliance_alliance_client（注：该顺序不能乱，因为后一个容器的启动依赖于前一个容器的运行），三个容器的启动命令分别为：
```bash
docker run -d --net=host mongo的容器id
```
```bash
docker run -i -t --net=host alliance_alliance_server的容器id
```
（alliance_alliance_server容器启动时间可能会比较长，因为服务器在启动的过程中需要拉一些golang的依赖）
```bash
docker run -t -i --net=host alliance_alliance_client的容器id
```
以上三条命令中都有一个参数 --net=host ，这个参数的意思是以宿主机模式启动容器服务，容器可以使用宿主机的网络ip地址和端口，这样可以实现多容器之间的互联
## 具体操作
在经过上面几个步骤之后所有服务都成功启动，接下来可以在客户端输入命令开始操作了，客户端有以下命令可用:
1.注册并登录: register accountName password
 
2.登录: login accountName password
 
3.查看自己的公会信息: whichAlliance
 
4.查看当前所有公会: allianceList 
 
5.创建公会: createAlliance allianceName
 
6.加入公会: joinAlliance allianceName 
 
7.解散公会: dismissAlliance
 
8.查看公会物品: allianceItems
 
9.提交公会物品: storeItem itemType itemNum 

10.扩容公会仓库: increaseCapacity 
 
11.销毁公会物品: destroyItem index 

上面几个命令会检测参数个数，如果数量不对，会返回失败