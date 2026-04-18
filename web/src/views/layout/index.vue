<template>
    <div class="main-container">
        <!-- 顶部导航 -->
        <header class="header-content">
            <div class="header-left">
                <div class="header-left-wrap">
                    <el-link class="console-name" href="/console">
                        <el-icon class="console-icon"><Menu /></el-icon>
                        <span class="console-text">控制台</span>
                    </el-link>
                </div>
            </div>

            <div class="header-right">
                <span>{{ dateMessage }}</span>
                <span>欢迎您，{{ displayedName }}</span>
                <el-button link @click="onUserCenter">个人信息</el-button>
                <el-button link @click="logout">退出</el-button>
            </div>
        </header>

        <!-- 主体布局 -->
        <div class="admin-main">
            <aside class="admin-sidebar">
                <app-aside />
            </aside>
            <main class="admin-content">
                <router-view />
            </main>
        </div>
    </div>
</template>

<script>
import AppAside from "./aside.vue";
import { toLoginPage, toUserCenter } from "../../utils/common.js";
import { logout, basicInfo } from "../../api/basic.js";

export default {
    name: "LayoutIndex",
    components: { AppAside },

    data() {
        return {
            userInfo: { username: "", account: "" }, // 用户信息
            dateMessage: "", // 顶部日期文案
        };
    },

    computed: {
        // 拼接显示用户名+账号
        displayedName() {
            const { username = "", account = "" } = this.userInfo;
            return `${username}(${account})`;
        },
    },

    methods: {
        // 获取用户基础信息
        async getBasicInfo() {
            try {
                const res = await basicInfo();
                this.userInfo = res.payload?.userinfo || {};
            } catch (err) {
                console.error("获取用户信息失败:", err);
            }
        },

        // 退出登录
        async logout() {
            try {
                await logout();
                sessionStorage.removeItem("active-path");
                this.$cookies.remove("session");
                toLoginPage();
            } catch (err) {
                console.error("退出登录失败:", err);
            }
        },

        // 前往个人中心
        onUserCenter() {
            toUserCenter();
        },

        // 格式化顶部日期
        getCurrentTime() {
            const now = new Date();
            const year = now.getFullYear();
            const month = String(now.getMonth() + 1).padStart(2, "0");
            const day = String(now.getDate()).padStart(2, "0");
            const week = ["星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"];
            this.dateMessage = `今天是${year}年${month}月${day}日 ${week[now.getDay()]}`;
        },
    },

    created() {
        this.getBasicInfo();
        this.getCurrentTime();
    },
};
</script>

<style scoped lang="less">
/* 整体布局 */
.main-container {
    min-width: 1366px;
    height: 100vh;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    overflow: visible;
}

/* 顶部导航 */
.header-content {
    height: 50px;
    padding: 0 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #fff;
    border-bottom: 1px solid #e5e7eb;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.header-left {
    display: flex;
    align-items: center;
    gap: 8px;
}

.header-left-wrap {
    padding-left: 170px;
}

.console-icon {
    font-size: 20px;
    color: #409eff;
}

.console-name {
    font-size: 17px;
    font-weight: 500;
    color: #333;
}

.console-text {
    margin-left: 10px;
}

.header-right {
    display: flex;
    align-items: center;
    gap: 12px;
    color: #666;
    font-size: 14px;
}

/* 主体内容区 */
.admin-main {
    display: flex;
    flex: 1;
    overflow: hidden;
}

/* 左侧菜单 */
.admin-sidebar {
    width: 200px;
    background: #fff;
    border-right: 1px solid #e5e7eb;
    overflow-y: overlay;

    &::-webkit-scrollbar {
        width: 2px;
        height: 4px;
    }

    &::-webkit-scrollbar-track {
        background: transparent;
        border-radius: 2px;
    }

    &::-webkit-scrollbar-thumb {
        background: #e0e0e0;
        border-radius: 2px;
    }

    &::-webkit-scrollbar-thumb:hover {
        background: #d0d0d0;
    }
}

/* 右侧内容 */
.admin-content {
    flex: 1;
    background: #f9fafb;
    padding: 15px;
    overflow-y: auto;
    overflow-x: hidden;
}

/* 移动端适配 */
@media (max-width: 768px) {
    .admin-sidebar {
        width: auto;
    }
    .admin-content {
        background: #fff;
        padding: 0;
    }
}
</style>
