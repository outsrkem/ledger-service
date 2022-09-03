# -*- coding=utf-8 -*-
from flask import session, request

from models.m_iam import Iam
from models.m_rols import Role
from service import response_body
from models.m_permissions import Permissions
from settings import Logger
from service.utility import to_json, page_info

_log = Logger()


def permission_list(page=1, per_page=10):
    """分页查询权限列表"""
    is_have_permission, msg = check_permission(request.path, request.method)
    if is_have_permission:
        payload = dict()
        count = Permissions().find_by_permission_count()
        payload["page_info"] = page_info(count, per_page, page)
        payload["items"] = Permissions().find_by_permission(page, per_page)
        return response_body(200, "", payload)
    else:
        return response_body(403, msg)


def roles_list(page=1, per_page=10):
    """分页查询权角色（用户组）"""
    is_have_permission, msg = check_permission(request.path, request.method)
    if is_have_permission:
        payload = dict()
        count = Role().find_by_role_count()
        payload["page_info"] = page_info(count, per_page, page)
        payload["items"] = Role().find_by_role(page, per_page)
        return response_body(200, "", payload)
    else:
        return response_body(403, msg)


def check_permission(path, method):
    user_id = session.get('user_id')
    per_code_list = Iam().find_by_user_for_permission_code(user_id)
    permission_code = Permissions().find_by_api_permission(path, method)
    """判断是否具有权限"""
    if permission_code in per_code_list:
        _log.logger.info("Permission to check successfully. %s %s" % (method, path))
        return True, None
    elif '*:*:*' in per_code_list:
        _log.logger.info("Permission to check successfully. %s %s" % (method, path))
        return True, None
    else:
        _log.logger.info("Permission to check error. %s %s" % (method, path))
        return False, {"message": "No Permission.", "permission_code": permission_code}
