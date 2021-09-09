# -*- coding=utf-8 -*-
from flask import Blueprint, request, json
from service.transaction.s_recedisbus import daily_statistical, add_journal_account, monthly_summarizing
from service.transaction.s_recedisbus import annual_summarizing
from service.transaction.s_recedisbus import monthly_statistical
from service.transaction.s_recedisbus import annual_statistical

recedisbu = Blueprint('recedisbu', __name__)


@recedisbu.route("/statistical/daily", methods=['GET'])
def r_daily_statistical():
    # 按天统计
    # /api/v1/recedisbu/statistical/daily?y=2021&m=10&d=2
    y = request.args.get('y')
    m = request.args.get('m')
    d = request.args.get('d')
    row = daily_statistical(y, m, d)
    return row, row["meta_info"]["res_code"]


@recedisbu.route("/statistical/monthly", methods=['GET'])
def r_monthly_statistical():
    # 按月统计
    y = request.args.get('y')
    m = request.args.get('m')
    row = monthly_statistical(y, m)
    return row, row["meta_info"]["res_code"]


@recedisbu.route("/statistical/annual", methods=['GET'])
def r_annual_statistical():
    # 按年统计
    y = request.args.get('y')
    row = annual_statistical(y)
    return row, row["meta_info"]["res_code"]


@recedisbu.route("/journal/account", methods=['POST'])
def r_add_journal_account():
    row = add_journal_account(json.loads(request.get_data()))
    # 添加一条或多条记账
    return row, row["meta_info"]["res_code"]


@recedisbu.route("/summarizing/monthly", methods=['GET'])
def r_monthly_summarizing():
    # 查询月度汇总
    y = request.args.get('y')
    m = request.args.get('m')
    row = monthly_summarizing(y, m)
    return row, row["meta_info"]["res_code"]


@recedisbu.route("/summarizing/annual", methods=['GET'])
def r_annual_summarizing():
    # 查询年度汇总
    y = request.args.get('y')
    row = annual_summarizing(y)
    return row, row["meta_info"]["res_code"]
