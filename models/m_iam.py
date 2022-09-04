# -*- coding=utf-8 -*-

from models import dbconnect
from sqlalchemy import Table
from service.utility import now_timestamp

dbsession, dbmodel, metadata = dbconnect()


class Iam(dbmodel):
    __table__ = Table('ledger_permissions', metadata, autoload=True)

    @staticmethod
    def find_by_user_for_permission_code(user_id):
        """根据用户ID查询该用户的权限码信息"""
        sql = '''SELECT
                        permissions.permission_code
                    FROM
                        ledger_role_user ru,
                        ledger_role_permission rm,
                        ledger_permissions permissions
                    WHERE
                        ru.role_id = rm.role_id
                    AND rm.perm_id = permissions.id
                    AND ru.user_id = ''' + str(user_id) + ''' AND permissions.permission_code IS NOT NULL;'''
        results = dbsession.execute(sql)
        result = [dict(zip(result.keys(), result)) for result in results]
        permission_list = list()
        if len(result) > 0:
            # 删除结果中的无用信息
            for i in result:
                permission_list.append(i["permission_code"])
            return permission_list
        return list()

    @staticmethod
    def find_by_permission_for_role_id(role_id):
        sql = '''SELECT
                    ledger_permissions.id,
                    ledger_permissions.permission_code,
                    ledger_permissions.service_group,
                    ledger_permissions.service_group_title,
                    ledger_permissions.permission_type
                FROM
                    ledger_role_permission,
                    ledger_role,
                    ledger_permissions
                WHERE
                    ledger_role_permission.role_id = ledger_role.id
                AND ledger_role_permission.perm_id = ledger_permissions.id
                AND ledger_role.id = ''' + str(role_id) + '''
                ORDER BY
                    ledger_permissions.id ASC'''
        results = dbsession.execute(sql)
        result = [dict(zip(result.keys(), result)) for result in results]
        return result
