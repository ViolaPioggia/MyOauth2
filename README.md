# MyOauth2

小白尝试基于[golang/oauth2](https://github.com/golang/oauth2 )实现Oauth2自定义登录

- 利用 viper 实现 mysql 和 Oauth2 的相关配置（实现了 redis 的结构体，但未实现 authorization 的登录函数）

- 此为 server 端代码，client 端代码需要另外敲



## yaml 配置文件

db：

```yaml
mysql:
  username: root
  password: root
  addr: 127.0.0.1
  port: 3306
  db: chatroom_users
  charset: utf8mb4
  connMaxIdleTime: 1h
  connMaxLifeTime: 1h
  place: Asia%2FShanghai
  maxIdleConns: 10
  maxOpenConns: 5
```

oauth2：

```yaml
session:
  name: session_id
  secret_key: "asdfsdaf"
  max_age: 1200

oauth2:
  accessTokenExp: 1
  refreshTokenExp: 1
  isGenerateRefresh: true
  jwt_signed_key: "qwefsdfa1"
  client:

    - id: test_client_1
      secret: test_secret_1
      name: 测试应用1
      addr: http://localhost
      port: 9093


    - id: test_client_2
      secret: test_secret_2
      name: 测试应用2
      addr: http://localhost
      port: 9094
```



## Mysql的表设计

连接到的Mysql的表必须名为 users，至少含有 username，password，id 三个键，不然 gorm 无法连接到数据库

```go
package model

type User struct {
	Username int64  `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	ID       int64  `form:"id" json:"id"`
}
```



## Api接口说明

1，`GET`       `http://localhost:9096/oauth/authorize` 鉴权

![image-20230501223606954](docs/image-20230501223606954.png)

![image-20230501223842802](docs/image-20230501223842802.png)

获得code=MZYXNJK0ZGMTNZBHMS0ZZWZJLTKWNMUTMGQ4MZBLOTU1NGI4



2，`POST`       `http://localhost:9096/oauth/token`  code 获取 token

![image-20230501224012964](docs/image-20230501224012964.png)

响应得到

```json
{
	"access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2NsaWVudF8xIiwiZXhwIjoxNjgzMjExMjM1LCJzdWIiOiIyMDIyMjE0NzQwIn0.KBu9FqCG7nL9IP8t0Qkm_PAqUBcnpV8HMincug0IrZaRjtORZIW55U5SmC9BzTdVNIkdMqmL4V-ZAobGedQLLQ",
	"expires_in": 259200,
	"refresh_token": "OTZHMZM0ZJCTODRLZC01ZMZLLTHJOTETYWQ2ZMIXZJCYZWM0",
	"scope": "all",
	"token_type": "Bearer"
}
```



3,`POST`       `http://localhost:9096/oauth/token`  password 直接获取 token，安全系数相对较低

![image-20230501224337936](docs/image-20230501224337936.png)

```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2NsaWVudF8xIiwiZXhwIjoxNjgyOTU5NDIzLCJzdWIiOiLvv70ifQ.NqUoQuDR1DJ0fLcMigDN9-KyAL5bY8PZHrOcZP_xvJvsGHGcv8qr5JLix_c8qdIVhngCsX1B7Bw3j7X35XpXyg",
	"expires_in": 7200,
	"refresh_token": "NJQ4MDAZM2YTZGI0MI01OGJILWIZNMMTMMRHMDQ0NTJMMZMY",
	"scope": "all",
	"token_type": "Bearer"
}
```



4,`POST`       `http://localhost:9096/oauth/token`  利用 `refresh_token` 刷新 `access_token`



![image-20230501224558504](docs/image-20230501224558504.png)

```json
{
	"access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2NsaWVudF8xIiwiZXhwIjoxNjgyOTU5NTk0LCJzdWIiOiLvv70ifQ.VFXxO3AZluthDZxkkzBxCa1qi3jvxmTbEbb8dy95BvGNvbyeNYNgfTxvG1FgDYTDDpTfcJ9XSoOSMGzUOomDcA",
	"expires_in": 7200,
	"refresh_token": "NWZLYMY1YTQTYWEXNS01MJZHLTGXMMQTODQ3NDIWZWEWZDZM",
	"scope": "all",
	"token_type": "Bearer"
}
```

5,`GET`       `http://localhost:9096/verify`  检验获得的 `access_token` 是否有效



1. ![image-20230501224925887](docs/image-20230501224925887.png)
2. ![image-20230501224948242](docs/image-20230501224948242.png)

```json
{
	"client_id": "test_client_1",
	"domain": "http://localhost:9093",
	"expires_in": 6977,
	"user_id": ""
}
```



------

此项目多为个人使用，若有侵权，请联系本人。