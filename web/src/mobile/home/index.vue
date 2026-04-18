<template>
    <div class="home-container">
        <router-view class="page-content" />

        <van-tabbar route v-model="active">
            <!-- 账单：添加双击刷新 -->
            <van-tabbar-item to="/bill" icon="balance-list" @dblclick.native="onBillDoubleClick"> 账单 </van-tabbar-item>

            <van-tabbar-item to="/report" icon="info">报表</van-tabbar-item>
            <van-tabbar-item icon="records" @click.prevent="onAddTransactions">添加</van-tabbar-item>
            <van-tabbar-item to="/setting" icon="setting">设置</van-tabbar-item>
        </van-tabbar>
    </div>
    <AddTrans ref="AddTrans" />
</template>

<script>
import AddTrans from "../comp/AddTrans.vue";
export default {
    name: "HomeIndexMob",
    components: { AddTrans },
    data() {
        return {
            active: "",
        };
    },
    watch: {
        $route: {
            immediate: true,
            handler() {
                this.active = this.$route.path;
            },
        },
    },
    methods: {
        // 打开添加弹窗
        onAddTransactions() {
            this.$refs.AddTrans.onOpenDialog();
        },

        // ====================== 双击账单刷新 ======================
        onBillDoubleClick() {
            // 只在账单页才刷新
            if (this.active === "/bill") {
                this.$globalBus.emit("onRefresh");
            }
        },
    },
};
</script>

<style scoped>
.home-container {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: #e9ebee;
}
.page-content {
    flex: 1;
    overflow-y: auto;
}
</style>
