-- ----------------------------
-- Table structure for ledger_category
-- ----------------------------
CREATE TABLE IF NOT EXISTS `ledger_category`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '大类名称',
  `create_time` bigint(19) NULL DEFAULT 1000000000000,
  `update_time` bigint(19) NULL DEFAULT 1000000000000,
  PRIMARY KEY (`kid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1001 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '收支大类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ledger_subcategory
-- ----------------------------
CREATE TABLE IF NOT EXISTS `ledger_subcategory`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '子类名称',
  `category_id` int(11) NOT NULL DEFAULT 1 COMMENT '外键，关联大类表的 id',
  `create_time` bigint(19) NULL DEFAULT 1000000000000,
  `update_time` bigint(19) NULL DEFAULT 1000000000000,
  PRIMARY KEY (`kid`) USING BTREE,
  INDEX `category_id`(`category_id`) USING BTREE,
  CONSTRAINT `ledger_subcategory_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `ledger_category` (`kid`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1001 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '收支子类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ledger_detail
-- ----------------------------
CREATE TABLE IF NOT EXISTS `ledger_detail`  (
  `kid` int(11) NOT NULL AUTO_INCREMENT,
  `subcategory_id` int(11) NOT NULL COMMENT '关联到子类(间接关联大类)',
  `occ_time` int(11) NOT NULL COMMENT '支出或收入发生的年，（occurrence year）',
  `amount` decimal(19, 4) NOT NULL COMMENT '金额',
  `total` int(11) NOT NULL DEFAULT 1 COMMENT '是否计入本月收支，1：计入；0：不计入',
  `direction` int(11) NOT NULL DEFAULT 1 COMMENT '金钱流向，1：支出；2：收入',
  `status` int(11) NOT NULL DEFAULT 2 COMMENT '状态，1：封存，不可编辑；2：可编辑',
  `remark` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '详细备注',
  `create_time` bigint(19) NULL DEFAULT 1000000000000,
  `update_time` bigint(19) NULL DEFAULT 1000000000000,
  PRIMARY KEY (`kid`) USING BTREE,
  INDEX `category_id`(`subcategory_id`) USING BTREE,
  CONSTRAINT `ledger_detail_ibfk_1` FOREIGN KEY (`subcategory_id`) REFERENCES `ledger_subcategory` (`kid`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1001 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '收支明细' ROW_FORMAT = DYNAMIC;
