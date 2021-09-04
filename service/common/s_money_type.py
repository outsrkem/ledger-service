# -*- coding=utf-8 -*-
from flask import session

from models.m_money_type import MoneyType
from service.common import response_body
from service.utility import now_time_timestamp


def query_all_money_type():
    res = MoneyType().find_by_all_money_type()
    response = list()
    for i in res:
        if i["category_id"] == 1:
            i["alias"] = i["category"]
        else:
            i["alias"] = i["category"] + '-' + i["title"]
        response.append(i)
    if response:
        return response_body(200, '', response)
    return response_body(404, 'The query is empty!')


def add_deal_type(raw):
    for i in raw["raw_data"]:
        deal_types = MoneyType().find_by_money_type(i["category"])
        if not deal_types:  # 查询大类是否存在
            return response_body(406, f'The category({i["category"]}) was not found!')
        elif deal_types["category"] != 1:
            return response_body(406, f'The category({i["category"]}) was not found!')

    raw["uid"] = session['user_id']
    raw["create_time"] = now_time_timestamp()
    raw["update_time"] = now_time_timestamp()
    if MoneyType().inst_money_type(raw):
        return response_body(201)

    return response_body(500, 'Error')


def del_deal_type(type_id):
    if MoneyType().del_money_type(type_id):
        return response_body(200)
    return response_body(500, 'Error')


def update_deal_type(raw):
    raw["update_time"] = now_time_timestamp()
    if MoneyType().update_money_type(raw):
        return response_body(201)
    return response_body(500, 'Error')
