# -*- coding=utf-8 -*-
from flask import session

from models.m_menus import Menus
from models.m_options import Options
from service import response_body


def query_rest_options():
    options = dict()
    user_id = session.get('user_id')
    res = Options().find_by_options(user_id)
    for opt in res:
        options[opt["name"]] = opt["value"]
    return response_body(200, '', options)


def query_layout_menus():
    # 查询菜单，构建返回体
    menus = list()
    row_level_1 = Menus().find_by_menus_level1()
    rem_list = ['seq_sort', 'update_time', 'create_time', 'describes', 'm_code', 'm_level', 'parent_menu_id']
    if row_level_1:
        for level1_menu in row_level_1:
            level_2 = list()

            """删除字典中多余的key：https://www.imangodoc.com/91882797.html"""
            [level1_menu.pop(key) for key in rem_list]
            level1_menu["leaf_node"] = []
            row_leve_2 = Menus().find_by_menus_level2(level1_menu["id"])
            """没有子菜单则不处理子菜单"""
            if row_leve_2:
                for level2_menu in row_leve_2:
                    """删除字典中多余的key"""
                    [level2_menu.pop(key) for key in rem_list]
                    level2_menu["paths"] = level1_menu["paths"] + level2_menu["paths"]
                    level_2.append(level2_menu)

                level1_menu["leaf_node"] = level_2
            menus.append(level1_menu)

        response = dict()
        response["items"] = menus
        if response:
            return response_body(200, '', response)
    return response_body(404, 'The query menu is empty. Contact the administrator to configure the menu.')
