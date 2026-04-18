<template>
    <div class="table-container">
        <!-- 表格主体 -->
        <table class="custom-table">
            <!-- 表头 -->
            <thead>
                <tr>
                    <th v-for="(column, index) in columns" :key="index">
                        <el-text>{{ column.label }}</el-text>
                    </th>
                </tr>
            </thead>

            <!-- 表体 -->
            <tbody>
                <tr v-for="(row, rowIndex) in data" :key="rowIndex" :style="getRowStyle(row, rowIndex)">
                    <td v-for="(column, colIndex) in columns" :key="colIndex">
                        <!-- 支持自定义单元格内容（插槽） -->
                        <template v-if="column.slot">
                            <el-text>
                                <slot :name="column.slot" :row="row" :index="rowIndex"></slot>
                            </el-text>
                        </template>
                        <!-- 默认显示字段值（支持嵌套属性，如 order.serial） -->
                        <template v-else>
                            <el-text>
                                {{ getNestedValue(row, column.prop) }}
                            </el-text>
                        </template>
                    </td>
                </tr>
                <!-- 空数据提示 -->
                <tr v-if="data.length === 0">
                    <td :colspan="columns.length" class="empty-row">
                        {{ emptyText }}
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script>
export default {
    name: "MyTable",
    props: {
        // 表格数据
        data: {
            type: Array,
            default: () => [],
        },
        // 列配置
        columns: {
            type: Array,
            required: true,
            validator: (value) => {
                return value.every((col) => {
                    if (!col.label) return false;
                    if (!col.slot && !col.prop) return false;
                    return true;
                });
            },
        },
        // 空数据提示文本
        emptyText: {
            type: String,
            default: "暂无数据",
        },
        // 行样式回调函数：接收行数据和索引，返回样式对象
        rowStyle: {
            type: Function,
            default: () => ({}), // 默认返回空对象（无特殊样式）
        },
    },
    methods: {
        // 获取行样式
        getRowStyle(row, index) {
            const customStyle = this.rowStyle(row, index);
            return typeof customStyle === "object" && customStyle !== null ? customStyle : {};
        },
        // 解析嵌套属性（如 order.serial -> row.order.serial）
        getNestedValue(obj, path) {
            // 将路径按 '.' 分割成数组（如 "order.serial" -> ["order", "serial"]）
            const keys = path.split(".");
            // 递归遍历对象获取值
            return keys.reduce((current, key) => {
                // 如果当前值是对象且包含该键，继续深入；否则返回空
                return current && typeof current === "object" ? current[key] : "";
            }, obj);
        },
    },
};
</script>

<style scoped>
table {
    width: 100%;
    border-collapse: collapse;
}

th {
    font-weight: bold;
    white-space: nowrap;
}

th,
td {
    padding: 8px 12px;
    text-align: left;
    border: none;
    border-bottom: 1px solid #ddd;
}

tr:hover {
    background-color: #f5f7fa;
}

/* 优先级调整：确保自定义样式能覆盖hover样式 */
tr[style] {
    transition:
        background-color 0.2s,
        color 0.2s;
}

.empty-row {
    text-align: center;
    color: #9ca3af;
    padding: 32px;
}
/* 关键调整：仅当行有自定义color时，才让el-text继承颜色 */
tr[style*="color"] td :deep(.el-text),
tr[style*="color"] td :deep(.el-text__inner) {
    color: inherit !important;
}
</style>
