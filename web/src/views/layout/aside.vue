<template>
    <div>
        <div>
            <el-aside class="aside" style="width: auto">
                <p style="text-align: center">记账薄</p>
                <el-menu :default-active="activePath" unique-opened>
                    <el-menu-item index="/overview" @click="OnSwitchRoutes('/overview')">
                        <el-icon><DataLine /></el-icon>
                        <template #title><span>数据总览</span></template>
                    </el-menu-item>
                    <el-menu-item index="/transactions" @click="OnSwitchRoutes('/transactions')">
                        <el-icon><House /></el-icon>
                        <template #title><span>记账管理</span></template>
                    </el-menu-item>
                    <el-menu-item index="/category" @click="OnSwitchRoutes('/category')">
                        <el-icon><User /></el-icon>
                        <template #title><span>分类管理</span></template>
                    </el-menu-item>
                </el-menu>
            </el-aside>
        </div>
        <div style="padding: 20px">
            <el-button-group class="ml-4" size="small">
                <el-button :type="size.x" @click="onSetSize('small')">小</el-button>
                <el-button :type="size.z" @click="onSetSize('default')">中</el-button>
                <el-button :type="size.d" @click="onSetSize('large')">大</el-button>
            </el-button-group>
        </div>
    </div>
</template>

<script>
export default {
    name: "AppAside",
    components: {},
    props: {},
    data() {
        return {
            activePath: "",
            size: { x: "", z: "primary", d: "" },
        };
    },
    computed: {},
    watch: {},
    methods: {
        OnSwitchRoutes(activePath) {
            this.$router.push({ path: activePath });
        },
        // 设置元素尺寸
        onSetSize(size) {
            this.$globalBus.emit("element-size", size);
            const sizeConfigMap = {
                small: { x: "primary", z: "", d: "" },
                default: { x: "", z: "primary", d: "" },
                large: { x: "", z: "", d: "primary" },
            };
            this.size = sizeConfigMap[size] || sizeConfigMap.default;
        },
    },
    created() {
        this.$globalBus.on("updateActivePath", (data) => {
            this.activePath = data || "/";
        });
        this.onSetSize(window.localStorage.getItem("element-size")); // 加载默认元素尺寸
    },
};
</script>

<style scoped>
.el-aside {
    /* 处理菜单右边的阴影 */
    background-color: #ffffff;
    .el-menu {
        border-right: none;
    }
}
</style>
