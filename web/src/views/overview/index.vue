<template>
    <div style="padding-top: 10px">
        <el-card style="width: 100%">
            <el-statistic title="Daily active users" :value="268500" />
        </el-card>
    </div>
</template>

<script>
export default {
    name: "OverviewIndex",
    data() {
        return {};
    },
    methods: {},
    created() {
        this.$globalBus.emit("updateActivePath", "/overview");
    },
};
</script>

<style scoped lang="less"></style>
