<template>
    <el-card style="width: 100%">
        <template #header>
            <div class="my_refresh">
                <el-row>
                    <span>分类管理</span>
                    <span style="padding-left: 5px; padding-right: 5px"></span>
                </el-row>
                <el-row>
                    <el-button type="primary" :icon="Refresh" @click="onRefresh" :loading="loading" style="margin-left: 10px">刷新</el-button>
                </el-row>
            </div>
        </template>
        <el-tree v-loading="loading" show-checkbox :data="category" :props="defaultProps" @node-click="handleNodeClick"></el-tree>
    </el-card>
</template>

<script>
import { Getcategory } from "../../api/basic.js";
import { withDelay } from "../../utils/common.js";
import { Refresh } from "@element-plus/icons-vue";
import { formatTime } from "../../utils/date.js";
export default {
    name: "CategoryIndex",
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
            category: [{}, {}, {}],
            defaultProps: {
                label: "name",
            },
        };
    },
    methods: {
        formatDate(time) {
            return formatTime(time);
        },
        onCurrentChange(p) {
            this.page = p;
            this.loadGetCategory(this.pageSize, p);
        },
        onSizeChange(s) {
            this.pageSize = s;
            this.page = 1;
            this.loadGetCategory(s, 1);
        },
        loadGetCategory: async function (pageSize, page) {
            this.loading = true;
            const params = { direction: 2 };
            withDelay(() => Getcategory(params))
                .then((res) => {
                    this.category = res.payload.items || [];
                    return res;
                })
                .finally(() => {
                    this.loading = false;
                });
        },
        onRefresh() {
            this.loading = true;
            this.loadGetCategory(this.pageSize, this.page);
        },
    },
    created() {
        this.$globalBus.emit("updateActivePath", "/category");
        this.onRefresh();
    },
};
</script>

<style scoped lang="less"></style>
