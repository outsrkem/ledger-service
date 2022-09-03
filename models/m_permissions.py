# -*- coding=utf-8 -*-

from models import dbconnect
from sqlalchemy import Table, func
from service.utility import now_timestamp

dbsession, dbmodel, metadata = dbconnect()


class Permissions(dbmodel):
    __table__ = Table('ledger_permissions', metadata, autoload=True)

    @staticmethod
    def find_by_permission_count():
        count = dbsession.query(func.count(Permissions.id)).scalar()
        return count

    def find_by_permission(self, page, per_page):
        """查询所有权限列表"""
        results = dbsession.query(Permissions).paginate(page=page, per_page=per_page, error_out=False).items
        results_list = list()
        if len(results) == 0:
            return False
        for i in results:
            results_list.append({c.name: getattr(i, c.name) for c in self.__table__.columns})
        # 删除结果无用的信息
        for i in results_list:
            i.pop("create_time")
            i.pop("update_time")
        return results_list
