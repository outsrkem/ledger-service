# -*- coding=utf-8 -*-
from flask import session
from models.m_money_type import MoneyType
from service import response_body
from service.utility import now_timestamp
from settings import Logger

_log = Logger()


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
        payload = {"items": response}
        _log.logger.debug("Querying transaction Type. %s" % response)
        _log.logger.info("Querying transaction Type.")
        return response_body(200, '', payload)
    _log.logger.warning("Querying transaction Type. The query is empty.")
    return response_body(404, 'The query is empty!')


def add_deal_type(raw):
    for i in raw["raw_data"]:
        deal_types = MoneyType().find_by_money_type(i["category"])
        if not deal_types:  # 查询大类是否存在
            return response_body(406, f'The category({i["category"]}) was not found!')
        elif deal_types["category"] != 1:
            '''category 为 1 ，则表示为大类'''
            return response_body(406, f'A category id ({i["category"]}) is a subclass.')

    raw["uid"] = session['user_id']
    raw["create_time"] = now_timestamp()
    raw["update_time"] = now_timestamp()
    if MoneyType().inst_money_type(raw):
        return response_body(201)

    return response_body(500, 'Error')


def del_deal_type(raws):
    if MoneyType().del_money_type(raws):
        return response_body(200)
    return response_body(500, 'Error')


def update_deal_type(raw):
    raw["update_time"] = now_timestamp()
    if MoneyType().update_money_type(raw):
        return response_body(201)
    return response_body(500, 'Error')


def deal_title_typ():
    # 获取类别（只获取大类）
    res_money_type = MoneyType().find_by_all_money_type()
    money_type = list()
    for i in res_money_type:
        if i["category_id"] != 1:
            '''只获取大类'''
            continue
        money_type.append({"id": str(i["id"]), "title": i["category"]})
    payload = {"items": money_type}
    return response_body(200, None, payload)
