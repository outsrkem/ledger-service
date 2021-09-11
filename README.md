# ledger-service

### ledger: 账单

### Api

```
  POST /api/v1/ledger/common/user/register
  POST /api/v1/ledger/common/user/login
  POST /api/v1/ledger/common/user/logout
   GET /api/v1/ledger/common/deal/type
  POST /api/v1/ledger/common/deal/type
 PATCH /api/v1/ledger/common/deal/type
DELETE /api/v1/ledger/common/deal/type
   GET /api/v1/ledger/common/deal/title/type
   GET /api/v1/ledger/recedisbu/statistical/daily
   GET /api/v1/ledger/recedisbu/statistical/monthly
   GET /api/v1/ledger/recedisbu/statistical/annual
  POST /api/v1/ledger/recedisbu/journal/account
   GET /api/v1/ledger/recedisbu/summarizing/monthly
   GET /api/v1/ledger/recedisbu/summarizing/annual
```

### 注册用户

```
POST /api/v1/common/user/register
```

- Body参数

| 名称      | 类型   | 必选 | 描述                 |
| --------- | ------ | ---- | -------------------- |
| username  | string | 是   | 用户名               |
| password  | string | 是   | 密码                 |
| mobile    | string | 是   | 手机，会做为登录账号 |
| describes | string | 否   | 描述                 |

- 请求示例

```
POST https://api.ledger.com/api/v1/common/user/register
{
    "username": "张三",
    "password": "123456",
    "mobile": "1234521114",
    "describes": ""
}
```

### 用户登录

```
POST /api/v1/ledger/common/user/login
```

- Body参数

| 名称     | 类型   | 必选 | 描述     |
| -------- | ------ | ---- | -------- |
| account  | string | 是   | 登录账号 |
| password | string | 是   | 登录密码 |

- 请求示例

```
POST https://api.ledger.com/api/v1/ledger/common/user/login
{
    "account": "admin",
    "password": "123456"
}
```

### 用户登出

```
POST /api/v1/ledger/common/user/logout
```

### 统计

```
# 日统计
GET /api/v1/ledger/recedisbu/statistical/daily?y=2021&m=10&d=2
# 月统计
GET /api/v1/ledger/recedisbu/statistical/monthly?y=2021&m=10
# 年统计
GET /api/v1/ledger/recedisbu/statistical/annual?y=2021
```

### 添加流水记录

```
POST /api/v1/ledger/recedisbu/journal/account
```

- 请求示例

```
POST https://api.ledger.com/api/v1/ledger/recedisbu/journal/account
{
    "raw_data": [
        {
            "mtid": 1002,
            "occ_date": "2021-9-12",
            "amount": "200.54",
            "total": 1,
            "describes": "火车"
        },
        {
            "mtid": 1202,
            "occ_date": "2021-9-6",
            "amount": "1200",
            "total": 1,
            "describes": "飞机"
        }
}
```

### 查看交易类型

```
GET /api/v1/ledger/common/deal/type
```

- 请求参数

| 名称    | 类型 | IN    | 必选                 | 描述                       |
| ------- | ---- | ----- | -------------------- | -------------------------- |
| type_id | int  | query | 否（无参数查询所有） | 类型id，获取单个类型的详情 |

- 请求示例

```
GET https://api.ledger.com/api/v1/ledger/common/deal/type
GET https://api.ledger.com/api/v1/ledger/common/deal/type?type_id=1002
```

### 添加交易类型

```
POST /api/v1/ledger/common/deal/type
```

- Body参数

| 名称     | 类型  | 必选 | 描述     |
| -------- | ----- | ---- | -------- |
| raw_data | Array | 是   | 类型参数 |

- raw_data 参数

| 名称      | 类型   | 必选 | 描述     |
| --------- | ------ | ---- | -------- |
| title     | string | 是   | 名称     |
| category  | int    | 是   | 父类型id |
| describes | string | 否   | 描述     |

- 请求示例

```
POST https://api.ledger.com/api/v1/ledger/common/deal/type
{
    "raw_data": [
        {
            "title": "香烟",
            "category": 1011,
            "describes": "暂无"
        },
        {
            "title": "酒水",
            "category": 1011,
            "describes": "暂无"
        }
    ]
}
```

### 修改交易类型

```
PATCH /api/v1/ledger/common/deal/type
{
    "type_id": 2004,
    "title": "香烟",
    "describes": "暂无ssssssss",
    "status": 2
}
```

### 删除交易类型

```
DELETE /api/v1/ledger/common/deal/type
{
    "type_id": [
        1035,
        1036,
        1037
    ]
}
```

### 只获取大类别的名称和id，用于表格的渲染

```
GET /api/v1/ledger/common/deal/title/type
```

### 获取月账单

```
GET /api/v1/ledger/recedisbu/summarizing/monthly?y=2021&m=8
```

### 获取年账单

```
GET /api/v1/ledger/recedisbu/summarizing/annual?y=2021
```