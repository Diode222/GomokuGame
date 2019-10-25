# GomokuGame
GomokuGame是自动对弈平台的后端业务层，管理用户注册登陆，请求开始比赛，查看比赛结果相关数据等业务操作。

目前暂时没有实现前端，后端基本业务功能如注册登陆，开始比赛，获取比赛结果的基本功能已经完成，通过HTTP调用API接口可以进行相关操作。

要想进行比赛，首先需要fork GomokuGameImpl（https://github.com/Diode222/GomokuGameImpl）代码，实现其中impl_server.go文件中的Init()
和MakePiece()接口，注册时会用到该fork项目的地址。


### API

若为post请求，所有参数均为body的k/v对。

http://127.0.0.1:8080/register ([post]; params: 1. user_name, 2. password, 3. warehouse_addr)

http://127.0.0.1:8080/login ([post]; params: 1. user_name, 2. password; return: token (jwt生成的token，过期时间为10天))

......

### 相关项目传送门

GomokuGameImpl (https://github.com/Diode222/GomokuGameImpl, 对弈选手需要实现接口的项目)

GomokuGameReferee (https://github.com/Diode222/GomokuGameReferee, 对弈“裁判”， 负责一次对局的调度和比赛相关数据的传输)
