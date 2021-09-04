# -*- coding=utf-8 -*-
from flask import session

from models.m_money_type import MoneyType
from models.m_recedisbu import RecedisbuStatement
from service.common import response_body
from service.transaction import statistical_trim
from service.utility import now_time_timestamp


def daily_statistical(y, m, d):
    result = RecedisbuStatement().find_by_income_and_daily_statistical(y, m, d)
    response = statistical_trim(result)
    if response:
        return response_body(200, '', response)
    return response_body(404, 'The query is empty!')


def monthly_statistical(y, m):
    result = RecedisbuStatement().find_by_income_and_monthly_statistical(y, m)
    response = statistical_trim(result)
    if response:
        return response_body(200, '', response)
    return response_body(404, 'The query is empty!')


def annual_statistical(y):
    result = RecedisbuStatement().find_by_income_and_annual_statistical(y)
    response = statistical_trim(result)
    if response:
        return response_body(200, '', response)
    return response_body(404, 'The query is empty!')


def add_journal_account(raw):
    for i in raw["raw_data"]:
        if not MoneyType().find_by_money_type(i["mtid"]):  # 查询类别是否存在
            return response_body(406, f'The mtid({i["mtid"]}) was not found!')

    raw["uid"] = session['user_id']
    raw["create_time"] = now_time_timestamp()
    raw["update_time"] = now_time_timestamp()

    if RecedisbuStatement().inst_journal_account(raw):
        return response_body(201)

    return response_body(500, 'Error!')
