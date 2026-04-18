<template>
    <div>
        <el-card>
            <template #header>
                <div class="my_refresh">
                    <el-row>
                        <span>记账管理</span>
                        <span style="padding-left: 5px; padding-right: 5px"></span>
                    </el-row>
                    <el-row>
                        <el-button type="primary" @click="onAddTransactions" style="margin-left: 10px">记一笔</el-button>
                        <el-button type="primary" :icon="Refresh" @click="onRefresh" :loading="loading" style="margin-left: 10px">刷新</el-button>
                    </el-row>
                </div>
            </template>

            <!-- 只替换这里：原生el-table → MyTable，其余完全不变 -->
            <MyTable :data="transactions" :columns="columns" v-loading="loading">
                <!-- 交易时间插槽 -->
                <template #occ_time="{ row }">
                    {{ formatDate(row.occ_time) }}
                </template>
                <!-- 类别插槽 -->
                <template #category_title="{ row }">
                    {{ row.category_title }}
                </template>
                <!-- 金额插槽 -->
                <template #amount="{ row }">
                    <span
                        :style="{
                            color: row.amount > 0 ? 'red' : row.amount < 0 ? 'green' : 'black',
                        }"
                    >
                        {{ row.amount }}
                    </span>
                </template>

                <!-- 备注插槽 -->
                <template #remark="{ row }">
                    {{ row.remark }}
                </template>
                <!-- 操作插槽 -->
                <template #action="{ row }">
                    <el-button link type="primary" @click="onDetail(row)">查看详情</el-button>
                    <el-button link type="primary" @click="onUpdate(row)">修改</el-button>
                    <el-popconfirm class="box-item" :title="`删除：${row.amount}`" placement="left-end" @confirm="onDeleteTransactions(row)">
                        <template #reference>
                            <el-button link type="primary">删除</el-button>
                        </template>
                    </el-popconfirm>
                </template>
            </MyTable>

            <div class="pagination">
                <pagination :pageTotal="pageTotal" :pageSize="pageSize" @CurrentChange="onCurrentChange" @SizeChange="onSizeChange" />
            </div>
        </el-card>
        <add-transactions ref="AddTransactions" />
    </div>
</template>

<script>
import MyTable from "../../components/MyTable/MyTable.vue";
import { msgcon } from "../../utils/message.js";
import { GetTransactions, Getcategory, DelTransactions } from "../../api/basic.js";
import { withDelay, convertToLimitOffset } from "../../utils/common.js";
import { getCategoryPath } from "../../utils/category.js";
import { Refresh } from "@element-plus/icons-vue";
import { formatTime } from "../../utils/date.js";
export default {
    name: "TransactionsIndex",
    components: {
        MyTable,
    },
    setup() {
        return {
            Refresh,
        };
    },
    data() {
        return {
            loading: true,
            pageTotal: 0,
            pageSize: 10,
            page: 1,
            transactions: [{}, {}, {}],
            category: [],
            // 新增列配置，对应MyTable
            columns: [
                { label: "交易时间", slot: "occ_time" },
                { label: "类别", slot: "category_title" },
                { label: "金额(¥)", slot: "amount" },
                { label: "备注", slot: "remark" },
                { label: "操作", slot: "action" },
            ],
        };
    },
    methods: {
        formatDate(time) {
            return formatTime(time);
        },
        onCurrentChange(p) {
            this.page = p;
            this.loadGetPerson(this.pageSize, p);
        },
        onSizeChange(s) {
            this.pageSize = s;
            this.page = 1;
            this.loadGetPerson(s, 1);
        },
        loadGetPerson: async function (pageSize, page) {
            this.loading = true;
            const params = convertToLimitOffset(page, pageSize);
            withDelay(() => GetTransactions(params))
                .then((res) => {
                    let resp = res.payload.items || [];
                    let category = window.localStorage.getItem("category");
                    resp.map((item) => (item.category_title = getCategoryPath(JSON.parse(category), item.category_id) || "未知分类"));
                    this.transactions = resp;
                    this.pageTotal = res.payload.page_info.total || 0;
                })
                .finally(() => {
                    this.loading = false;
                });
        },
        loadGetCategory: async function () {
            Getcategory({ direction: 2 }).then((res) => {
                window.localStorage.setItem("category", JSON.stringify(res.payload.items));
            });
        },
        onRefresh() {
            this.loading = true;
            this.loadGetCategory();
            this.loadGetPerson(this.pageSize, this.page);
        },
        onAddTransactions() {
            this.$refs.AddTransactions.onOpenDialog();
        },
        // 查看详情
        onDetail(val) {
            console.log(val);
        },
        onUpdate(val) {
            console.log(val);
        },
        onDeleteTransactions(val) {
            DelTransactions(val.id)
                .then(() => {
                    this.$message.success(msgcon("删除成功"));
                    this.onRefresh();
                })
                .catch((err) => {
                    let msg = err.data.metadata.message;
                    this.$message.error(msgcon("删除失败 " + msg));
                });
        },
    },
    created() {
        this.$globalBus.emit("updateActivePath", "/transactions");
        this.onRefresh();
        this.$globalBus.on("onRefresh", () => {
            this.onRefresh();
        });
    },
};
</script>

<style scoped lang="less"></style>
