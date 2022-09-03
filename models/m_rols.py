# -*- coding=utf-8 -*-

from models import dbconnect
from sqlalchemy import Table, func
from service.utility import now_timestamp

dbsession, dbmodel, metadata = dbconnect()


class Role(dbmodel):
    __table__ = Table('ledger_role', metadata, autoload=True)

    @staticmethod
    def find_by_role_count():
        count = dbsession.query(func.count(Role.id)).scalar()
        return count

    def find_by_role(self, page, per_page):
        """查询所有角色"""
        results = dbsession.query(Role).paginate(page=page, per_page=per_page, error_out=False).items
        results_list = list()
        if len(results) == 0:
            return False
        for i in results:
            results_list.append({c.name: getattr(i, c.name) for c in self.__table__.columns})
        return results_list
