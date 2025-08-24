export function getCategoryPath(categories, targetId) {
    // 递归查找函数
    function findPath(items, currentPath = []) {
        for (const item of items) {
            // 复制当前路径并添加当前项名称
            const newPath = [...currentPath, item.name];

            // 如果找到目标ID，返回拼接后的路径
            if (item.id === targetId) {
                return newPath.join("-");
            }

            // 如果有子分类，递归查找
            if (item.children && item.children.length > 0) {
                const result = findPath(item.children, newPath);
                if (result) {
                    return result;
                }
            }
        }
        // 未找到返回null
        return null;
    }

    return findPath(categories);
}

// // 测试数据
// const categories = [
//     { id: 1001, name: "转账" },
//     { id: 1002, name: "消费" },
//     {
//         id: 1003,
//         name: "餐饮",
//         children: [
//             { id: 100301, name: "早餐" },
//             { id: 100302, name: "午餐" },
//             // 其他子分类...
//         ],
//     },
//     // 其他分类...
// ];

// // 测试示例
// console.log(getCategoryPath(categories, 100301)); // 输出: "餐饮-早餐"
// console.log(getCategoryPath(categories, 100403)); // 输出: "购物-穿搭美妆"
// console.log(getCategoryPath(categories, 1001)); // 输出: "转账"
