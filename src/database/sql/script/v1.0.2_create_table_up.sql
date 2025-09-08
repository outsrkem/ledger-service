-- 收入类一级分类（direction=1），ID段：5001-9999（与支出类1001-5000隔离）
INSERT INTO `ledger_category` 
    (`kid`, `instance_id`, `name`, `direction`, `layer`, `create_time`) 
VALUES 
    (5001, NULL, '薪资', 1, 1, 1630664740000),
    (5002, NULL, '理财', 1, 1, 1630664740000),
    (5003, NULL, '红包', 1, 1, 1630664740000),
    (5004, NULL, '转账', 1, 1, 1630664740000),
    (5005, NULL, '借入', 1, 1, 1630664740000),
    (5006, NULL, '收款', 1, 1, 1630664740000),
    (5007, NULL, '出售', 1, 1, 1630664740000),
    (5008, NULL, '其他', 1, 1, 1630664740000);

-- 3. 出售子分类（layer=2）
INSERT INTO `ledger_category` 
    (`kid`, `instance_id`, `name`, `direction`, `layer`, `create_time`) 
VALUES 
    (500701, NULL, '个人创作', 1, 2, 1630664740000),
    (500702, NULL, '二手', 1, 2, 1630664740000),
    (500702, NULL, '废品', 1, 2, 1630664740000);

INSERT INTO `ledger_caterela` 
    (`parent_id`, `child_id`, `create_time`) 
VALUES 
    (5007, 500701, 1630664740000),
    (5007, 500702, 1630664740000),
    (5007, 500703, 1630664740000);
