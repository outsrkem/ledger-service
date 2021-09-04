# ledger-service

#### ledger: 账单

### api
```
# 注册
POST /api/v1/common/user/register
{
    "username": "admin",
    "password": "123456",
    "mobile": "12345214",
    "describes":""
}

# 登录
POST /api/v1/common/user/login
{
    "account": "admin",
    "password": "123456"
}

# 登出
POST /api/v1/common/user/logout

# 日统计
GET /api/v1/recedisbu/statistical/daily?y=2021&m=10&d=2
# 月统计
GET /api/v1/recedisbu/statistical/monthly?y=2021&m=10
# 年统计
GET /api/v1/recedisbu/statistical/annual?y=2021

# 添加流水记录
POST /api/v1/recedisbu/journal/account
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

################
# 查看交易类型
GET /api/v1/common/deal/type

# 添加交易类型
POST   /api/v1/common/deal/type
{
    "raw_data": [
        {
            "title": "香烟",
            "category": "1011",
            "describes": "暂无"
        },
        {
            "title": "酒水",
            "category": "1011",
            "describes": "暂无"
        }
    ]
}

# 修改交易类型
PATCH  /api/v1/common/deal/type
{
    "type_id": 2004,
    "title": "香烟",
    "describes": "暂无ssssssss",
    "status": 2
}

# 删除交易类型
DELETE /api/v1/common/deal/type?type_id=2004
```

