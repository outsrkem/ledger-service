# -*- coding=utf-8 -*-
from app_flask import db


def create_ledger_options():
    # 配置表
    db.session.execute('''
    CREATE TABLE IF NOT EXISTS `ledger_options` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `uid` int(11) NOT NULL DEFAULT 0 COMMENT '0: 系统内置',
      `name` varchar(32) NOT NULL ,
      `title` varchar(32) DEFAULT '',
      `value` text,
      `options_type` varchar(32) DEFAULT '' COMMENT '用于标识一组资源',
      `describes` varchar(255) DEFAULT NULL,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1001 DEFAULT CHARSET=utf8 COMMENT='配置表';
    ''')


def create_ledger_user():
    # 用户表
    db.session.execute('''
    CREATE TABLE IF NOT EXISTS `ledger_user` (
      `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '0: 系统管理员，最高操作权限',
      `account` varchar(50) NOT NULL DEFAULT '',
      `username` varchar(50) NOT NULL DEFAULT '',
      `password` varchar(300) DEFAULT NULL,
      `mobile` varchar(15) NOT NULL DEFAULT '',
      `status` int(11) NOT NULL DEFAULT 1,
      `describes` varchar(255) DEFAULT '',
      `create_time` bigint(19) DEFAULT 1000000000000,
      `update_time` bigint(19) DEFAULT 1000000000000,
      PRIMARY KEY (`id`),
      UNIQUE KEY `account` (`account`)
    ) ENGINE=InnoDB AUTO_INCREMENT=2001 DEFAULT CHARSET=utf8 COMMENT='用户表';
    ''')


def create_ledger_role():
    # 角色表
    db.session.execute('''
    CREATE TABLE IF NOT EXISTS `ledger_role` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `uid` int(11) NOT NULL DEFAULT 1 COMMENT '0: 系统内置，除管理员外不可编辑',
      `role_name` varchar(255) NOT NULL,
      `role_state` int(11) NOT NULL,
      `role_type` int(11) DEFAULT NULL COMMENT '1,系统角色，不可删除；2,自定义角色',
      `status` int(11) NOT NULL DEFAULT 1,
      `describes` varchar(255) DEFAULT NULL,
      `create_time` bigint(19) DEFAULT 1000000000000,
      `update_time` bigint(19) DEFAULT 1000000000000,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=2001 DEFAULT CHARSET=utf8 COMMENT='角色表';
    ''')


def create_ledger_permissions():
    # 权限码表
    db.session.execute('''
    CREATE TABLE IF NOT EXISTS `ledger_permissions` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `uid` int(11) NOT NULL DEFAULT 0 COMMENT '0: 系统内置，除管理员外不可编辑',
      `permission_code` varchar(255) NOT NULL COMMENT '权限码',
      `status` int(11) NOT NULL DEFAULT 1 COMMENT '1: 启用，正常；0：禁用',
      `describes` varchar(255) DEFAULT NULL COMMENT '描述，说明',
      `create_time` bigint(19) DEFAULT 1000000000000,
      `update_time` bigint(19) DEFAULT 1000000000000,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1001 DEFAULT CHARSET=utf8 COMMENT='权限码表';
    ''')


def create_ledger_role_permissions():
    # 角色权限关联表
    db.session.execute('''
    CREATE TABLE IF NOT EXISTS `ledger_role_permission` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `uid` int(11) NOT NULL DEFAULT 0 COMMENT '0: 系统内置，除管理员外不可编辑',
      `role_id` int(11) NOT NULL COMMENT '',
      `perm_id` int(11) NOT NULL COMMENT '',
      `create_time` bigint(19) DEFAULT 1000000000000,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1001 DEFAULT CHARSET=utf8 COMMENT='角色权限关联表';
    ''')


def create_ledger_role_user():
    # 角色用户关联表
    db.session.execute('''
    CREATE TABLE IF NOT EXISTS `ledger_role_user` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `uid` int(11) NOT NULL DEFAULT 0 COMMENT '0: 系统内置，除管理员外不可编辑',
      `role_id` int(11) NOT NULL COMMENT '',
      `user_id` int(11) NOT NULL COMMENT '',
      `create_time` bigint(19) DEFAULT 1000000000000,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1001 DEFAULT CHARSET=utf8 COMMENT='角色用户关联表';
    ''')


def create_ledger_money_type():
    db.session.execute('''
    CREATE TABLE IF NOT EXISTS `ledger_money_type` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `uid` int(11) NOT NULL DEFAULT 0,
      `title` varchar(255) NOT NULL COMMENT '资金的支出或收入类型',
      `category` int(11) DEFAULT 1 COMMENT '类别，1：代表大类，不可变，子类的值为该大类的id',
      `status` int(11) NOT NULL DEFAULT 1,
      `describes` varchar(255) DEFAULT '',
      `create_time` bigint(19) DEFAULT 1000000000000,
      `update_time` bigint(19) DEFAULT 1000000000000,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=2001 DEFAULT CHARSET=utf8 COMMENT='资金来往分类表';
    ''')


def create_ledger_recedisbu_statement():
    # receipt and disbursement statement
    db.session.execute('''
    CREATE TABLE IF NOT EXISTS `ledger_recedisbu_statement` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `uid` int(11) NOT NULL DEFAULT 0 COMMENT '记录所有者，必选',
      `mtid` int(11) NOT NULL COMMENT 'ledger_money_type 的id',
      `occ_year` int(11) NOT NULL COMMENT '支出或收入发生的年，（occurrence year）',
      `occ_month` int(11) NOT NULL COMMENT '支出或收入发生的月',
      `occ_day` int(11) NOT NULL COMMENT '支出或收入发生的日',
      `amount` DECIMAL(19,4) NOT NULL COMMENT '金额',
      `total` int(11) NOT NULL DEFAULT 1 COMMENT '是否计入本月收支，1：计入；0：不计入',
      `direction` int(11) NOT NULL DEFAULT 1 COMMENT '金钱流向，1：支出；2：收入',
      `status` int(11) NOT NULL DEFAULT 2 COMMENT '状态，1：封存，不可编辑；2：可编辑',
      `describes` varchar(255) DEFAULT '' COMMENT '详细备注',
      `create_time` bigint(19) DEFAULT 1000000000000,
      `update_time` bigint(19) DEFAULT 1000000000000,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=2001 DEFAULT CHARSET=utf8 COMMENT='收支明细';
    ''')


def create_ledger_menus():
    db.session.execute('''
    CREATE TABLE IF NOT EXISTS `ledger_menus` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `name` char(255) NOT NULL,
      `paths` varchar(255) DEFAULT NULL COMMENT '请求路径',
      `m_code` varchar(255) DEFAULT NULL COMMENT '权限码',
      `parent_menu_id` int(11) NOT NULL DEFAULT 0 COMMENT '父菜单id',
      `m_level` int(255) NOT NULL DEFAULT 1 COMMENT '菜单级别。1、2、3',
      `seq_sort` int(11) NOT NULL COMMENT '菜单顺序',
      `icon_name` varchar(255) DEFAULT '',
      `describes` varchar(255) DEFAULT '',
      `create_time` bigint(19) DEFAULT 1000000000000,
      `update_time` bigint(19) DEFAULT 1000000000000,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=2001 DEFAULT CHARSET=utf8 COMMENT='菜单表';
    ''')


def initialize_sql():
    create_ledger_options()
    create_ledger_user()
    create_ledger_role()
    create_ledger_permissions()
    create_ledger_role_permissions()
    create_ledger_role_user()
    create_ledger_money_type()
    create_ledger_recedisbu_statement()
    create_ledger_menus()
    return True


if __name__ == '__main__':
    if initialize_sql():
        print("The database is successfully initialized. Procedure")
