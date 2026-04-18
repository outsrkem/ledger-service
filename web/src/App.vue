<script>
export default {
    data() {
        return {
            size: "default", // 可选值: large/default/small
        };
    },
    created() {},
};
</script>

<template>
    <el-config-provider size="default">
        <div>
            <router-view />
        </div>
    </el-config-provider>
</template>

<style scoped></style>
