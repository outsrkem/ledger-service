# -*- coding=utf-8 -*-
from flask import Blueprint, request, json, session

from service.common.s_menus import query_layout_menus
from service.common.s_money_type import add_deal_type, del_deal_type, update_deal_type, query_all_money_type

common = Blueprint('common', __name__)


@common.route("/layout/menus", methods=['GET'])
def r_layout_menus():
    # 首页菜单
    row = query_layout_menus()
    return row, row["meta_info"]["res_code"]


@common.route("/user/register", methods=['POST'])
def r_user_register():
    from service.user.s_user import user_register
    row = user_register(json.loads(request.get_data()))
    return row, row["meta_info"]["res_code"]


@common.route("/user/login", methods=['POST'])
def r_user_login():
    from service.user.s_user import user_login
    row = user_login(json.loads(request.get_data()))
    return row, row["meta_info"]["res_code"]


@common.route("/user/logout", methods=['POST'])
def r_user_logout():
    if session.get('islogin'):
        session.clear()
        return ''
    return '', 401


@common.route("/deal/type", methods=['GET'])
def r_query_all_money_type():
    row = query_all_money_type()
    return row, row["meta_info"]["res_code"]


@common.route("/deal/type", methods=['POST'])
def r_add_deal_type():
    # 添加交易类型
    row = add_deal_type(json.loads(request.get_data()))
    return row, row["meta_info"]["res_code"]


@common.route("/deal/type", methods=['PATCH'])
def r_patch_deal_type():
    # 修改交易类型
    row = update_deal_type(json.loads(request.get_data()))
    return row, row["meta_info"]["res_code"]


@common.route("/deal/type", methods=['DELETE'])
def r_delet_deal_type():
    # 删除交易类型，接收如下格式数据
    # {"type_id":[1035,1036,1037]}
    row = del_deal_type(request.get_json())
    return row, row["meta_info"]["res_code"]
