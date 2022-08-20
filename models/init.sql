-- -
-- 请手工创建数据库
-- -
CREATE
    DATABASE /*!32312 IF NOT EXISTS */ ledgerdb /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

GRANT ALL
    ON ledgerdb.* TO ledger@'%' IDENTIFIED BY '123456';

GRANT ALL
    ON ledgerdb.* TO ledger@'localhost' IDENTIFIED BY '123456';

-- -
-- 请手工插入初始化数据
-- -
-- -
-- ledger_user
-- -
INSERT INTO `ledgerdb`.`ledger_user` (`id`, `account`, `username`, `password`, `mobile`, `status`, `describes`)
VALUES ('1000', 'admin', 'Administrator',
        'pbkdf2:sha256:260000$kqOMyzAytNwtWUCH$1dda32f737a990ce2b9f61cca29c355a663d00887d7858d7cfbed201be17def3',
        'admin', '1', '');
-- -
-- ledger_menus
-- -
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1001', '用户管理', '/users', NULL, '0', '1', '2');
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1002', '用户总览', '/preview', NULL, '1001', '2', '1');
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1003', '用户注册', '/create', NULL, '1001', '2', '2');
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1004', '用户管里', '/aregulator', NULL, '1001', '2', '3');
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1005', '记账管理', '/bill', NULL, '0', '1', '1');
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1006', '添加记账', '/journal/add', NULL, '1005', '2', '1');
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1007', '每日明细', '/particulars', NULL, '1005', '2', '2');
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1008', '年度汇总', '/statistical/annual', NULL, '1005', '2', '3');
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1009', '月度汇总', '/statistical/monthly', NULL, '1005', '2', '4');
INSERT INTO `ledgerdb`.`ledger_menus` (`id`, `name`, `paths`, `m_code`, `parent_menu_id`, `m_level`, `seq_sort`)
VALUES ('1010', '账单浏览', '/list', NULL, '1005', '2', '5');
--
-- ledger_money_type
--
INSERT INTO `ledgerdb`.`ledger_money_type` (`id`, `uid`, `title`, `category`, `status`, `describes`)
VALUES ('1001', '0', '餐饮', '1', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`id`, `uid`, `title`, `category`, `status`, `describes`)
VALUES ('1002', '0', '医疗', '1', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`id`, `uid`, `title`, `category`, `status`, `describes`)
VALUES ('1003', '0', '文化', '1', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`id`, `uid`, `title`, `category`, `status`, `describes`)
VALUES ('1004', '0', '服饰', '1', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`id`, `uid`, `title`, `category`, `status`, `describes`)
VALUES ('1005', '0', '物品', '1', '1', '');
-- 子类
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('早饭', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('午饭', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('午饭', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('晚饭', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('夜宵', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('聚餐', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('酒水', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('香烟', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('茶叶', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('蔬菜', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('肉类', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('大米', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('零食', '1001', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('医保', '1002', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('保健药物', '1002', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('医院', '1002', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('体检', '1002', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('门诊', '1002', '1', '');
INSERT INTO `ledgerdb`.`ledger_money_type` (`title`, `category`, `status`, `describes`)
VALUES ('衣服', '1004', '1', '');