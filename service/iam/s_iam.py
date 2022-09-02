# -*- coding=utf-8 -*-
from flask import session

from service import response_body
from models.m_permissions import Permissions
from settings import Logger
from service.utility import to_json

_log = Logger()


def roles_list(page=1, per_page=10):
    """分页查询权限列表"""
    # _log.logger.info("Querying the User List; page:%s,per_page:%s" % (page, per_page))
    payload = Permissions().find_by_permission(page, per_page)
    # payload = Iam().find_by_user_for_permission_code(2001)
    return response_body(200, "", payload)
