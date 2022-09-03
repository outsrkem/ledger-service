# -*- coding=utf-8 -*-
from flask import session

from models.m_rols import Role
from service import response_body
from models.m_permissions import Permissions
from settings import Logger
from service.utility import to_json, page_info

_log = Logger()


def permission_list(page=1, per_page=10):
    """分页查询权限列表"""
    payload = dict()
    count = Permissions().find_by_permission_count()
    payload["page_info"] = page_info(count, per_page, page)
    payload["items"] = Permissions().find_by_permission(page, per_page)
    return response_body(200, "", payload)


def roles_list(page=1, per_page=10):
    """分页查询权角色（用户组）"""
    payload = dict()
    count = Role().find_by_role_count()
    payload["page_info"] = page_info(count, per_page, page)
    payload["items"] = Role().find_by_role(page, per_page)
    return response_body(200, "", payload)
