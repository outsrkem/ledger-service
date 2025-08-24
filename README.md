# ledger-service

### ledger: 账单

### Api

### 添加流水记录

```
POST /api/v1/finances/transactions
```

-   请求示例

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

```sql
-- 添加触发器
-- 插入前校验
DELIMITER //
CREATE TRIGGER `check_detail_total_before_insert`
BEFORE INSERT ON `ledger_detail`
FOR EACH ROW
BEGIN
    IF ROUND(NEW.quantity * NEW.unit, 4) != NEW.total THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '明细总价必须等于数量×单价（保留4位小数）';
    END IF;
END //
DELIMITER ;

-- 更新前校验
DELIMITER //
CREATE TRIGGER `check_detail_total_before_update`
BEFORE UPDATE ON `ledger_detail`
FOR EACH ROW
BEGIN
    IF ROUND(NEW.quantity * NEW.unit, 4) != NEW.total THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '明细总价必须等于数量×单价（保留4位小数）';
    END IF;
END //
DELIMITER ;

-- 查询
SHOW TRIGGERS;
-- 删除
DROP TRIGGER IF EXISTS check_detail_total_before_insert;
DROP TRIGGER IF EXISTS check_detail_total_before_update;
```
