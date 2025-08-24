<script>
export default {
    data() {
        return {
            size: "default", // 可选值: large/default/small
        };
    },
    created() {
        this.$globalBus.on("element-size", (data) => {
            this.size = data ?? "default";
            window.localStorage.setItem("element-size", this.size);
        });
    },
};
</script>

<template>
    <el-config-provider :size="size">
        <div>
            <router-view />
        </div>
    </el-config-provider>
</template>

<style scoped></style>
