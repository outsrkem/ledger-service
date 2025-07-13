# ledger-service

### ledger: 账单

### Api

### 添加流水记录

```
POST /api/v1/finances/transactions
```

- 请求示例

```
s{
    "transactions": [
        {
            "mtid": 1002,
            "occ_date": "2021-9-12",
            "amount": "200.54",
            "total": 1,
            "remark": "火车"
        },
        {
            "mtid": 1202,
            "occ_date": "2021-9-6",
            "amount": "1200",
            "total": 1,
            "remark": "飞机"
        }
}
```

