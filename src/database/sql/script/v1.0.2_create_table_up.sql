-- 收入类一级分类（direction=1），ID段：5001-9999（与支出类1001-5000隔离）
INSERT INTO `ledger_category` 
    (`kid`, `instance_id`, `name`, `direction`, `pid`, `create_time`) 
VALUES 
    (5001, NULL, '薪资', 1, NULL, 1630664740000),
    (5002, NULL, '理财', 1, NULL, 1630664740000),
    (5003, NULL, '红包', 1, NULL, 1630664740000),
    (5004, NULL, '转账', 1, NULL, 1630664740000),
    (5005, NULL, '借入', 1, NULL, 1630664740000),
    (5006, NULL, '收款', 1, NULL, 1630664740000),
    (5007, NULL, '其他', 1, NULL, 1630664740000);
