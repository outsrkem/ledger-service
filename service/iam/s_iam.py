# -*- coding=utf-8 -*-
from flask import session

from service import response_body
from models.m_permissions import Permissions
from settings import Logger
from service.utility import to_json, page_info

_log = Logger()


def roles_list(page=1, per_page=10):
    """分页查询权限列表"""
    payload = dict()
    count = Permissions().find_by_permission_count()
    payload["page_info"] = page_info(count, per_page, page)
    payload["items"] = Permissions().find_by_permission(page, per_page)
    return response_body(200, "", payload)
