-- ----------------------------
-- Table structure for ledger_instance
-- ----------------------------
CREATE TABLE `ledger_instance`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT COMMENT '实例ID（主键）',
  `instance_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '所属用户ID',
  `create_time` bigint(19) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`kid`) USING BTREE,
  UNIQUE INDEX `idx_unique_user_id`(`user_id`) USING BTREE,
  UNIQUE INDEX `idx_unique_instance_id`(`instance_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1000000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '记账实例表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ledger_budget
-- ----------------------------
CREATE TABLE `ledger_budget`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT COMMENT '预算ID（主键）',
  `instance_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '所属实例ID',
  `category_id` int(11) NOT NULL COMMENT '关联分类ID',
  `amount` decimal(10, 4) NOT NULL COMMENT '预算金额',
  `cycle` tinyint(4) NOT NULL COMMENT '周期（1=月度，2=季度，3=年度）',
  `start_time` bigint(19) NOT NULL COMMENT '预算开始时间',
  `end_time` bigint(19) NULL DEFAULT NULL COMMENT '预算结束时间（null=长期有效）',
  `create_time` bigint(19) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`kid`) USING BTREE,
  UNIQUE INDEX `uk_instance_budget`(`instance_id`, `category_id`, `cycle`, `start_time`) USING BTREE,
  INDEX `fk_budget_instance`(`instance_id`) USING BTREE,
  CONSTRAINT `ledger_budget_ibfk_1` FOREIGN KEY (`instance_id`) REFERENCES `ledger_instance` (`instance_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1000000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '预算计划表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ledger_category
-- ----------------------------
CREATE TABLE `ledger_category`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT COMMENT '分类ID（主键）',
  `instance_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '所属实例ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '分类名称（如“餐饮”“工资”）',
  `direction` tinyint(4) NOT NULL COMMENT '类型（1=收入，2=支出）',
  `pid` int(11) NULL DEFAULT NULL COMMENT '父分类ID（二级分类用）',
  `create_time` bigint(19) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`kid`) USING BTREE,
  UNIQUE INDEX `uk_instance_category_direction`(`instance_id`, `name`, `direction`) USING BTREE,
  INDEX `fk_category_instance`(`instance_id`) USING BTREE,
  INDEX `fk_category_pid`(`pid`) USING BTREE,
  CONSTRAINT `fk_category_instance` FOREIGN KEY (`instance_id`) REFERENCES `ledger_instance` (`instance_id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_category_pid` FOREIGN KEY (`pid`) REFERENCES `ledger_category` (`kid`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1000000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '收支分类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ledger_transaction
-- ----------------------------
CREATE TABLE `ledger_transaction`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT COMMENT '交易ID（主键）',
  `instance_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '所属实例ID',
  `category_id` int(11) NOT NULL COMMENT '关联分类ID',
  `amount` decimal(10, 4) NOT NULL COMMENT '金额（正数=收入，负数=支出）',
  `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '交易备注',
  `occ_time` char(24) NOT NULL COMMENT '交易发生时间(ISO8601)',
  `create_time` bigint(19) NOT NULL COMMENT '记录创建时间',
  `update_time` bigint(19) NOT NULL COMMENT '记录更新时间',
  PRIMARY KEY (`kid`) USING BTREE,
  INDEX `fk_transaction_instance`(`instance_id`) USING BTREE,
  INDEX `fk_transaction_category`(`category_id`) USING BTREE,
  INDEX `idx_instance_occ_time`(`instance_id`, `occ_time`) USING BTREE,
  CONSTRAINT `ledger_transaction_ibfk_1` FOREIGN KEY (`instance_id`) REFERENCES `ledger_instance` (`instance_id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_trans_category` FOREIGN KEY (`category_id`) REFERENCES `ledger_category` (`kid`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1000000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '收支交易主表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ledger_detail
-- ----------------------------
CREATE TABLE `ledger_detail`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT COMMENT '明细ID（主键）',
  `transaction_id` int(11) NOT NULL COMMENT '关联交易ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '物品名称（如“牛奶”）',
  `quantity` decimal(10, 4) NOT NULL COMMENT '数量（支持小数）',
  `price` decimal(10, 4) NOT NULL COMMENT '单价',
  `total` decimal(10, 4) NOT NULL COMMENT '该项总价',
  `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '明细备注（如“促销装”）',
  PRIMARY KEY (`kid`) USING BTREE,
  INDEX `fk_detail_transaction`(`transaction_id`) USING BTREE,
  CONSTRAINT `fk_detail_transaction` FOREIGN KEY (`transaction_id`) REFERENCES `ledger_transaction` (`kid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1000000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '交易明细表' ROW_FORMAT = DYNAMIC;

CREATE TABLE `ledger_tag`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT COMMENT '标签ID（主键）',
  `instance_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '所属实例ID',
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标签名称（如“日常用品”“工作餐”）',
  `color` varchar(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '#333333' COMMENT '标签颜色（十六进制，如#FF0000）',
  `create_time` bigint(19) NOT NULL COMMENT '创建时间',
  `update_time` bigint(19) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`kid`) USING BTREE,
  UNIQUE INDEX `uk_instance_tag_name`(`instance_id`, `name`) USING BTREE COMMENT '同一实例下标签名称唯一',
  INDEX `fk_tag_instance`(`instance_id`) USING BTREE,
  CONSTRAINT `fk_tag_instance` FOREIGN KEY (`instance_id`) REFERENCES `ledger_instance` (`instance_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1000000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '标签表' ROW_FORMAT = DYNAMIC;

CREATE TABLE `ledger_trantag`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT COMMENT '关联ID（主键）',
  `transaction_id` int(11) NOT NULL COMMENT '关联交易ID',
  `tag_id` int(11) NOT NULL COMMENT '关联标签ID',
  `create_time` bigint(19) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`kid`) USING BTREE,
  UNIQUE INDEX `uk_transaction_tag`(`transaction_id`, `tag_id`) USING BTREE COMMENT '同一交易不能重复关联同一标签',
  INDEX `fk_transaction_tag_tag`(`tag_id`) USING BTREE,
  CONSTRAINT `fk_transaction_tag_tag` FOREIGN KEY (`tag_id`) REFERENCES `ledger_tag` (`kid`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_transaction_tag_transaction` FOREIGN KEY (`transaction_id`) REFERENCES `ledger_transaction` (`kid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1000000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '账单标签关联表' ROW_FORMAT = DYNAMIC;
