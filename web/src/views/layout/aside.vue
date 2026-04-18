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
        };
    },
    computed: {},
    watch: {},
    methods: {
        OnSwitchRoutes(activePath) {
            this.$router.push({ path: activePath });
        },
    },
    created() {
        this.$globalBus.on("updateActivePath", (data) => {
            this.activePath = data || "/";
        });
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
