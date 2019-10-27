# GomokuGame
GomokuGame是自动对弈平台的后端业务层，管理用户注册登陆，请求开始比赛，查看比赛结果相关数据等业务操作。

目前暂时没有实现前端，后端基本业务功能如注册登陆，开始比赛，获取比赛结果的基本功能已经完成，通过HTTP调用API接口可以进行相关操作。

要想进行比赛，首先需要fork GomokuGameImpl（https://github.com/Diode222/GomokuGameImpl）代码，实现其中impl_server.go文件中的Init()
和MakePiece()接口，注册时会用到该fork项目的地址。


### API

若为post请求，应区分post body中的参数与query参数。

##### 注册
http://127.0.0.1:8080/register ([post]; postparams: 1. user_name, 2. password, 3. warehouse_addr)

##### 登陆
http://127.0.0.1:8080/login ([post]; postparams: 1. user_name, 2. password; return: token (jwt生成的token，过期时间为10天))

##### 开启一局对战
http://127.0.0.1:8080/game/start ([get]; header: 1. token, 2. Content-type: application/x-www-form-urlencoded; query: 1. player1_first_hand (true/false, 选择先手还是后手), 2. max_thinking_time (本局游戏单步最大思考时间), 3. enemy_user_name (optional，挑战选手的用户名)

##### 获取一局对战的游戏结果
http://127.0.0.1:8080/game/result ([get]; header: 1. token, 2. Content-type: application/x-www-form-urlencoded; query 1. game_id (表示一局对战的唯一id))

......

### 相关项目传送门

GomokuGameImpl (https://github.com/Diode222/GomokuGameImpl, 对弈选手需要实现接口的项目)

GomokuGameReferee (https://github.com/Diode222/GomokuGameReferee, 对弈“裁判”， 负责一次对局的调度和比赛相关数据的传输)
