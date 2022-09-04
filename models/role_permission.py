# -*- coding=utf-8 -*-
from flask import session

from models import dbconnect
from sqlalchemy import Table, func
from service.utility import now_timestamp

dbsession, dbmodel, metadata = dbconnect()


class RolePermission(dbmodel):
    __table__ = Table('ledger_role_permission', metadata, autoload=True)

    @staticmethod
    def find_by_permission_for_role_id_count(role_id):
        """查询用户总数目"""
        count = dbsession.query(RolePermission).filter(RolePermission.role_id == role_id).count()
        return count
