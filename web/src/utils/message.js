// 提示消息内容
export function msgcon(msg) {
    return {
        showClose: true,
        message: msg,
        duration: 1500,
        grouping: true,
        plain: true,
    };
}
