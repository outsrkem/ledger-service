<template>
    <div>
        <el-card>
            <template #header>
                <div class="my_refresh">
                    <span>记账管理</span>
                    <el-space>
                        <el-button type="success" @click="onAddBill">记一笔</el-button>
                        <el-button type="primary" :icon="Refresh" @click="onRefresh" :loading="loading">刷新</el-button>
                    </el-space>
                </div>
            </template>

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
                            color: row.amount > 0 ? '#f53f3f' : row.amount < 0 ? 'green' : 'black',
                        }">
                        {{ row.amount }}
                    </span>
                </template>

                <!-- 备注插槽 -->
                <template #remark="{ row }">
                    {{ row.remark }}
                </template>
                <!-- 操作插槽 -->
                <template #action="{ row }">
                    <el-button link type="primary" @click="onDetail(row)">详情</el-button>
                    <el-button link type="primary" @click="onUpdate(row)">修改</el-button>
                    <el-button link type="primary" @click="onDeleteTransactions(row)">删除</el-button>
                </template>
            </MyTable>

            <div class="pagination">
                <pagination :pageTotal="pageTotal" :pageSize="pageSize" @syncsize="onSyncsize" @CurrentChange="onCurrentChange" @SizeChange="onSizeChange" />
            </div>
        </el-card>
        <add-transactions ref="AddTransactions" />
        <UpdateTransactions ref="UpdateTransactions" />
        <TranDetail ref="TranDetail" />
    </div>
</template>

<script>
import { Refresh } from "@element-plus/icons-vue";
import { ElMessageBox } from "element-plus";
import MyTable from "../../components/MyTable/MyTable.vue";
import AddTransactions from "./AddTransactions.vue";
import UpdateTransactions from "./update.vue";
import TranDetail from "./detail.vue";
import { GetTransactions, Getcategory, DelTransactions } from "../../api/basic.js";
import { msgcon } from "../../utils/message.js";
import { withDelay, convertToLimitOffset } from "../../utils/common.js";
import { getCategoryPath } from "../../utils/category.js";
import { formatTime } from "../../utils/date.js";

const SK_LEDGER_ALL_CATEGORY = "LEDGER_ALL_CATEGORY";

export default {
    name: "TransactionsIndex",
    components: {
        MyTable,
        AddTransactions,
        UpdateTransactions,
        TranDetail,
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
            transactions: Array.from({ length: 5 }, () => ({})),
            category: [],
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

        onSyncsize(s, p) {
            this.pageSize = s;
            this.page = p;
            this.onRefresh();
        },

        onCurrentChange(p) {
            this.page = p;
            this.loadGetBill(this.pageSize, p);
        },

        onSizeChange(s) {
            this.pageSize = s;
            this.page = 1;
            this.loadGetBill(s, 1);
        },

        async loadGetBill(pageSize, page) {
            this.loading = true;
            const params = convertToLimitOffset(page, pageSize);
            try {
                const res = await withDelay(() => GetTransactions(params));
                const resp = res.payload?.items || [];

                const categoryStr = window.localStorage.getItem(SK_LEDGER_ALL_CATEGORY);
                let categoryTree = [];
                try {
                    categoryTree = categoryStr ? JSON.parse(categoryStr) : [];
                } catch (e) {
                    this.$message.error(msgcon("分类数据解析失败：" + e));
                    categoryTree = [];
                }

                resp.forEach((item) => {
                    item.category_title = getCategoryPath(categoryTree, item.category_id) || "未知分类";
                });

                this.transactions = resp;
                this.pageTotal = res.payload?.page_info?.total || 0;
            } catch (err) {
                this.$message.error(msgcon("账单列表请求异常：" + err));
            } finally {
                this.loading = false;
            }
        },

        /**
         * Load income & expense categories, merge and save to localStorage
         */
        async loadGetCategory() {
            try {
                const results = await Promise.allSettled([Getcategory({ direction: 1 }), Getcategory({ direction: 2 })]);

                let mergedCategoryList = [];
                results.forEach((result) => {
                    if (result.status === "fulfilled") {
                        const res = result.value;
                        const rawItems = res.payload?.items;
                        const safeItems = Array.isArray(rawItems) ? rawItems : [];
                        mergedCategoryList.push(...safeItems);
                    } else {
                        // Only print to console to avoid frequent popup notifications
                        console.error("Category sub-request failed:", result.reason);
                        // this.$message.error(msgcon("分类请求异常 " + result.reason));
                    }
                });

                window.localStorage.setItem(SK_LEDGER_ALL_CATEGORY, JSON.stringify(mergedCategoryList));
            } catch (err) {
                this.$message.error(msgcon("分类加载异常：" + err));
            }
        },

        /**
         * Refresh data. Load categories first to eliminate race condition
         */
        async onRefresh() {
            this.loading = true;
            // Wait category request finished before loading bills, prevent unknown category display
            await this.loadGetCategory();
            this.loadGetBill(this.pageSize, this.page);
        },

        onAddBill() {
            this.$refs.AddTransactions.onOpenDialog();
        },

        onDetail(val) {
            this.$refs.TranDetail.onOpenDialog(val);
        },
        onUpdate(val) {
            this.$refs.UpdateTransactions.onOpenDialog(val);
        },

        onDeleteTransactions(val) {
            ElMessageBox.confirm(`确定删除这条金额为【${val.amount}】的账单吗？`, "删除确认", {
                confirmButtonText: "确认删除",
                cancelButtonText: "取消",
                type: "warning",
                draggable: true,
            })
                .then(async () => {
                    await DelTransactions(val.id);
                    this.$message.success(msgcon("删除成功"));
                    this.onRefresh();
                })
                .catch((err) => {
                    if (err) {
                        const msg = err.data?.metadata?.message || "未知错误";
                        this.$message.error(msgcon("删除失败：" + msg));
                    }
                });
        },
    },
    created() {
        this.$globalBus.emit("updateActivePath", "/transactions");
        this.refreshHandler = () => this.onRefresh();
        this.$globalBus.on("onRefresh", this.refreshHandler);
    },
    beforeUnmount() {
        // Remove global bus listener to avoid memory leak and duplicate trigger
        this.$globalBus.off("onRefresh", this.refreshHandler);
    },
};
</script>
