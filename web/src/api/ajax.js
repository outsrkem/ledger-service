import request from "../utils/request.js";

/**
 * 请求函数模块 ajax 封装
 * 向外部暴露一个函数 ajax
 * 共有如下4个参数
 * @param {*} url 请求路径，默认为空
 * @param {*} method 请求方法，默认为GET
 * @param {*} params 请求参数，默认为空
 * @param {*} data 请求参数，默认为空
 */
export default function ajax(url = "", method = "GET", params, data) {
    /** 返回值 Promise 对象，异步返回的数据是response.data，而不是response */
    return new Promise(function (resolve, reject) {
        /**
         * 利用axios异步执行ajax请求
         * 发送请求 这个promise用来保存axios的返回值(promise对象)
         * resolve() : 成功回调
         * reject() : 失败回调
         */
        const promise = request({
            method,
            url,
            params,
            data,
        });
        promise
            .then((response) => {
                /** 成功回调resolve() */
                resolve(response.data);
            })
            .catch((error) => {
                /** 失败回调reject() */
                reject(error);
            });
    });
}
