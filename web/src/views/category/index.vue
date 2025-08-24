<template>{{ category }}</template>

<script>
import { Getcategory } from "../../api/basic.js";
import { withDelay, convertToLimitOffset } from "../../utils/common.js";
import { Refresh } from "@element-plus/icons-vue";
import { formatTime } from "../../utils/date.js";
export default {
    name: "categoryIndex",
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
        this.onRefresh();
    },
};
</script>

<style scoped lang="less"></style>
