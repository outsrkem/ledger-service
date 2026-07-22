/*
 * 与后台交互模块 （依赖已封装的ajax函数）
 * 包含n个接口请求函数的模块，函数的返回值是promise对象
 */
import ajax from "./ajax";
/**
 * ajax 有如下4个参数
 * @param {*} url 请求路径，默认为空
 * @param {*} method 请求方法，默认为GET
 * @param {*} params 请求参数，默认为空对象
 * @param {*} data 请求参数，默认为空对象
 */

/** *********************************用户相关**************************************************** */
export const logout = () => ajax("/api/uias/v1/user/logout", "POST");
export const basicInfo = () => ajax("/api/uias/v1/uias/user/basicInfo", "GET", null, null);

// 需要添加前缀: /api/ledger
// PATCH /v1/transactions/:id   修改记账
export const GetTransactions = (params) => ajax("/api/ledger/v1/transactions", "GET", params, null);
export const AddTransactions = (data) => ajax("/api/ledger/v1/transactions", "POST", null, data);
export const EditTransactions = (id, data) => ajax(`/api/ledger/v1/transactions/${id}`, "PATCH", null, data);
export const DelTransactions = (id) => ajax(`/api/ledger/v1/transactions/${id}`, "DELETE", null, null);
export const GetTranDetails = (id) => ajax(`/api/ledger/v1/transactions/${id}`, "GET", null, null);

export const GetStatistic = (params) => ajax("/api/ledger/v1/bill/statistic", "GET", params, null);
export const GetBillDetails = (id) => ajax(`/api/ledger/v1/bill/${id}`, "GET", null, null);

// 数据总览
export const GetCycleSummary = (params) => ajax(`/api/ledger/v1/stat/cycle-summary`, "GET", params, null);
export const GetChartDay = (params) => ajax(`/api/ledger/v1/stat/chart/day`, "GET", params, null);

// 分类管理
//    GET /v1/category                    查询记账单
// DELETE /v1/category/:category_id       删除分类
//   POST /v1/category/:category_id/sub   添加子分类 body:{"name":"加油"}
//   POST /v1/category/main               添加主分类 body:{"name":"小车","direction":2}  direction => 1:收入,2:支出
export const Getcategory = (params) => ajax("/api/ledger/v1/category", "GET", params, null);
export const DelCategory = (id) => ajax(`/api/ledger/v1/category/${id}`, "DELETE", null, null);
export const AddSubCategory = (paths, data) => ajax(`/api/ledger/v1/category/${paths.id}/sub`, "POST", null, data);
export const AddMainCategory = (data) => ajax(`/api/ledger/v1/category/main`, "POST", null, data);
