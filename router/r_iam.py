# -*- coding=utf-8 -*-
from flask import Blueprint, request, json, session

from service.iam.s_iam import roles_list
from service.common.s_money_type import deal_title_type

iam = Blueprint('iam', __name__)


@iam.route("/roles/list", methods=['GET'])
def r_roles_list():
    page, per_page = request.args.get('page', type=int), request.args.get('per_page', type=int)
    row = roles_list(page, per_page)
    return row
