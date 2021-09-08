# -*- coding=utf-8 -*-
import copy
from decimal import Decimal
from flask import session
from models.m_money_type import MoneyType
from models.m_recedisbu import RecedisbuStatement
from service import response_body
from service.transaction import statistical_trim, formatting_type_decimal, statistical_type_decimal_to_str
from service.utility import now_timestamp


def daily_statistical(y, m, d):
    result = RecedisbuStatement().find_by_income_and_daily_statistical(y, m, d)
    res = statistical_type_decimal_to_str(statistical_trim(result))
    response = dict()
    response["items"] = res
    if response:
        return response_body(200, '', response)
    return response_body(404, 'The query is empty!')


def monthly_statistical(y, m):
    result = RecedisbuStatement().find_by_income_and_monthly_statistical(y, m)
    res = statistical_type_decimal_to_str(statistical_trim(result))
    response = dict()
    response["items"] = res
    if response:
        return response_body(200, '', response)
    return response_body(404, 'The query is empty!')


def annual_statistical(y):
    result = RecedisbuStatement().find_by_income_and_annual_statistical(y)
    res = statistical_type_decimal_to_str(statistical_trim(result))
    response = dict()
    response["items"] = res
    if response:
        return response_body(200, '', response)
    return response_body(404, 'The query is empty!')


def add_journal_account(raw):
    for i in raw["raw_data"]:
        if not MoneyType().find_by_money_type(i["mtid"]):  # 查询类别是否存在
            return response_body(406, f'The mtid({i["mtid"]}) was not found!')

    raw["uid"] = session['user_id']
    raw["create_time"] = now_timestamp()
    raw["update_time"] = now_timestamp()

    if RecedisbuStatement().inst_journal_account(raw):
        return response_body(201)

    return response_body(500, 'Error!')


def monthly_summarizing(y, m):
    # 月汇总，含每天的统计
    res = RecedisbuStatement().find_by_all_month_record(y, m)
    res_money_type = MoneyType().find_by_all_money_type()
    items = dict()
    period = y + "-" + m
    day_list = []
    rows = list()
    money_type = list()
    _d_temporary = dict()
    amount_sum_list = list()
    amount_sum = Decimal()
    for i in res_money_type:
        if i["category_id"] != 1:
            '''只获取大类'''
            continue
        money_type.append({"id": str(i["id"]), "title": i["category"]})
        _d_temporary[str(i["id"])] = ""

    if res is False:
        '''当月没有记录'''
        items["period"] = period
        items["amount_sum"] = str(amount_sum)
        items["money_type"] = money_type
        payload = {"items": items}
        msg = "There is no record of income and expenditure this month"
        return response_body(200, msg, payload)

    for day in res:
        '''获取月的天,处理每天的数据'''
        day_list.append(day["occ_day"])
        d = day["occ_day"]
        aa = statistical_trim(RecedisbuStatement().find_by_income_and_daily_statistical(y, m, d))
        _d = copy.deepcopy(_d_temporary)
        _d["date"] = aa["period"]
        _d["amount"] = aa["amount_sum"]
        for i in aa["details"]:
            _d[str(i["money_type"])] = str(i["amount_total"])
        amount_sum_list.append(_d["amount"])

        rows.append(_d)

    '''汇总'''
    for i in amount_sum_list:
        amount_sum = amount_sum + i
    amount_sum = amount_sum.quantize(Decimal('0.000'))

    items["period"] = period
    items["amount_sum"] = str(amount_sum)
    items["money_type"] = money_type
    '''转换数据类型'''
    items["rows"] = formatting_type_decimal(rows, "amount")

    payload = {"items": items}
    return response_body(200, None, payload)
