<template>
    <div class="el-aside">
        <p style="text-align: center">社区管理系统</p>
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
            menusList: [],
        };
    },
    computed: {},
    watch: {},
    created() {
        this.$globalBus.on("updateActivePath", (data) => {
            this.activePath = data || "/";
        });
    },
    methods: {
        OnSwitchRoutes(activePath) {
            this.$router.push({ path: activePath });
        },
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
