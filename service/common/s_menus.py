# -*- coding=utf-8 -*-
from flask import session

from models.m_menus import Menus
from service import response_body


def query_layout_menus():
    # 查询菜单，构建返回体
    menus = list()
    level_1 = Menus().find_by_menus_level1()
    for i in level_1:
        level_2 = list()
        del i['seq_sort']
        del i['update_time']
        del i['create_time']
        del i['describes']
        del i['m_code']
        del i['m_level']
        del i['parent_menu_id']

        for m in Menus().find_by_menus_level2(i["id"]):
            del m['seq_sort']
            del m['update_time']
            del m['create_time']
            del m['describes']
            del m['m_code']
            del m['m_level']
            del m['parent_menu_id']
            level_2.append(m)

        i["leaf_node"] = level_2
        menus.append(i)

    response = dict()
    response["items"] = menus
    if response:
        return response_body(200, '', response)
    return response_body(404, 'The query is empty!')
