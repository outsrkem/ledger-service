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
            <el-table :data="transactions" style="width: 100%" v-loading="loading">
                <el-table-column prop="category_title" label="类别" />
                <el-table-column label="金额(¥)">
                    <template #default="scope">
                        <span
                            :style="{
                                color: scope.row.amount > 0 ? 'red' : scope.row.amount < 0 ? 'green' : 'black',
                            }"
                        >
                            {{ scope.row.amount }}
                        </span>
                    </template>
                </el-table-column>
                <el-table-column prop="remark" label="备注" />
                <el-table-column prop="create_time" label="交易时间" width="200">
                    <template #default="scope">{{ formatDate(scope.row.occ_time) }}</template>
                </el-table-column>
                <el-table-column width="300" label="操作">
                    <template #default="scope">
                        <el-button link type="primary" @click="onDetail(scope.row)">查看详情</el-button>
                        <el-button link type="primary" @click="onUpdate(scope.row)">修改</el-button>
                        <el-popconfirm
                            class="box-item"
                            :title="`删除：${scope.row.amount}`"
                            placement="left-end"
                            @confirm="onDeleteTransactions(scope.row)"
                        >
                            <template #reference>
                                <el-button link type="primary">删除</el-button>
                            </template>
                        </el-popconfirm>
                    </template>
                </el-table-column>
            </el-table>
            <div class="pagination">
                <pagination :pageTotal="pageTotal" :pageSize="pageSize" @CurrentChange="onCurrentChange" @SizeChange="onSizeChange" />
            </div>
        </el-card>
        <add-transactions ref="AddTransactions" />
    </div>
</template>

<script>
import { msgcon } from "../../utils/message.js";
import { GetTransactions, Getcategory, DelTransactions } from "../../api/basic.js";
import { withDelay, convertToLimitOffset } from "../../utils/common.js";
import { getCategoryPath } from "../../utils/category.js";
import { Refresh } from "@element-plus/icons-vue";
import { formatTime } from "../../utils/date.js";
export default {
    name: "TransactionsIndex",
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
