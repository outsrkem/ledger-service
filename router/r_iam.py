# -*- coding=utf-8 -*-
from flask import Blueprint, request, json, session

from service.iam.s_iam import permission_list, roles_list
from service.common.s_money_type import deal_title_type

iam = Blueprint('iam', __name__)


@iam.route("/permission/list", methods=['GET'])
def r_permission_list():
    page, per_page = request.args.get('page', type=int), request.args.get('pageSize', type=int)
    row = permission_list(page, per_page)
    return row


@iam.route("/roles/list", methods=['GET'])
def r_role_list():
    page, per_page = request.args.get('page', type=int), request.args.get('pageSize', type=int)
    row = roles_list(page, per_page)
    return row
