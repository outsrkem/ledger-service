CREATE TABLE IF NOT EXISTS `ledger_category`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '资金的支出或收入类型',
  `pid` int(11) NOT NULL DEFAULT 1 COMMENT '类别，1：代表大类，不可变，子类的值为该大类的id',
  `type` int(11) NOT NULL DEFAULT 1,
  `sort` int(11) NOT NULL,
  `create_time` bigint(19) NULL DEFAULT 1000000000000,
  `update_time` bigint(19) NULL DEFAULT 1000000000000,
  PRIMARY KEY (`kid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1001 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '资金来往分类表' ROW_FORMAT = DYNAMIC;


--
CREATE TABLE IF NOT EXISTS `ledger_detail`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) NOT NULL COMMENT 'ledger_money_type 的id',
  `occ_time` int(11) NOT NULL COMMENT '支出或收入发生的年，（occurrence year）',
  `amount` decimal(19, 4) NOT NULL COMMENT '金额',
  `total` int(11) NOT NULL DEFAULT 1 COMMENT '是否计入本月收支，1：计入；0：不计入',
  `direction` int(11) NOT NULL DEFAULT 1 COMMENT '金钱流向，1：支出；2：收入',
  `status` int(11) NOT NULL DEFAULT 2 COMMENT '状态，1：封存，不可编辑；2：可编辑',
  `remark` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '详细备注',
  `create_time` bigint(19) NULL DEFAULT 1000000000000,
  `update_time` bigint(19) NULL DEFAULT 1000000000000,
  PRIMARY KEY (`kid`) USING BTREE,
  INDEX `category_id`(`category_id`) USING BTREE,
  CONSTRAINT `ledger_detail_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `ledger_category` (`kid`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1001 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '收支明细' ROW_FORMAT = DYNAMIC;
