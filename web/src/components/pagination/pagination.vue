<template>
    <div class="el-pagination">
        <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="currentPage"
            :page-size="pageSize"
            :page-sizes="[10, 15, 25, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pageTotal" />
    </div>
</template>

<script>
export default {
    name: "PaginationIndex",
    props: {
        pageTotal: { type: Number, default: 0 },
        pageSize: { type: Number, default: 15 },
    },
    data() {
        return {
            currentPage: 1,
        };
    },
    methods: {
        handleSizeChange(page_size) {
            this.currentPage = 1;
            this.$emit("SizeChange", page_size);
            // this.setUrlQuery(page_size)
        },
        handleCurrentChange(current_page) {
            this.currentPage = current_page;
            this.$emit("CurrentChange", current_page, this.pageSize);
            // this.setUrlQuery()
        },
    },
};
</script>

<style scoped>
.el-pagination {
    padding-top: 0px;
    padding-bottom: 0;
}
</style>
