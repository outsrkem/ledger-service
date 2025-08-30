<template>
    <div class="admin-layout">
        <header class="header-content">
            <div class="header-left">
                <el-space :size="10" spacer="">
                    <span style="padding-left: 160px"></span>
                    <el-icon class="console-icon"><Menu /></el-icon>
                    <el-link class="console-name" href="/console">控制台</el-link>
                </el-space>
            </div>
            <div class="header-right">
                <span>{{ dateMessage }}</span>
                <span>欢迎您，{{ displayedName }}</span>
                <div>
                    <el-button size="default" link @click="onUserCenter">个人信息</el-button>
                </div>
                <div>
                    <el-button size="default" link @click="Logout">退出</el-button>
                </div>
            </div>
        </header>
        <!-- 2. 主体内容区（左右分栏） -->
        <div class="admin-main">
            <aside class="admin-sidebar">
                <app-aside />
            </aside>
            <main class="admin-content">
                <!-- 子路由出口 -->
                <router-view />
            </main>
        </div>
    </div>
</template>

<script>
import AppAside from "./aside.vue";
import { toLoginPage, toUserCenter, toConsole } from "../../utils/common.js";
import { logout, basicInfo } from "../../api/basic.js";
export default {
    name: "LayoutIndex",
    components: {
        AppAside,
    },
    props: {},
    data() {
        return {
            userInfo: {},
            breadcrumb: [], // 面包屑导航
            dateMessage: "",
        };
    },
    computed: {
        displayedName() {
            return this.userInfo.username + "(" + this.userInfo.account + ")";
        },
    },
    methods: {
        LoadLogOut: async function () {
            await logout().then(() => {
                window.sessionStorage.removeItem("active-path");
                this.$cookies.remove("session");
                toLoginPage();
            });
        },
        GetbasicInfo: async function () {
            const res = await basicInfo();
            this.userInfo = res.payload.userinfo;
        },
        Logout() {
            this.LoadLogOut();
        },
        onUserCenter() {
            toUserCenter();
        },
        onConsole() {
            toConsole();
        },
        CurrentTime() {
            // 返回一个对象，包含日期、时间和星期几
            const now = new Date();
            const year = now.getFullYear();
            const month = String(now.getMonth() + 1).padStart(2, "0"); // 月份是从0开始的，所以要+1
            const day = String(now.getDate()).padStart(2, "0");
            // const hours = String(now.getHours()).padStart(2, "0");
            // const minutes = String(now.getMinutes()).padStart(2, "0");
            // const seconds = String(now.getSeconds()).padStart(2, "0");
            const weekdays = ["星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"];
            const _weekday = weekdays[now.getDay()];
            this.dateMessage = `今天是${year}年${month}月${day}日 ${_weekday}`;
        },
    },
    created() {
        this.GetbasicInfo();
        this.CurrentTime();
    },
};
</script>

<style scoped lang="less">
.admin-layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
    overflow: hidden;
    margin: 0; /* 清除body默认margin导致的留白 */
    padding: 0;
    min-width: 1200px;
}

.header-content {
    height: 50px;
    padding: 0 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    background-color: #ffffff;
    border-bottom: 1px solid #e5e7eb;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.header-left {
    display: flex;
    align-items: center;
    gap: 8px; /* 图标与文字间距 */
}

.console-icon {
    font-size: 20px;
    color: #409eff;
}

.console-name {
    font-size: 17px;
    font-weight: 500;
    color: #333333;
}

.header-right {
    display: flex;
    align-items: center;
    gap: 12px; /* 文字与按钮间距 */
    color: #666666;
    font-size: 14px;
}

.admin-main {
    display: flex;
    flex: 1;
    overflow: hidden;
    flex-grow: 1;
}

/* 左侧菜单样式（默认白色） */
.admin-sidebar {
    width: 200px; /* 固定菜单宽度 */
    background-color: #ffffff; /* 菜单默认白色 */
    border-right: 1px solid #e5e7eb; /* 右侧分隔线 */
}

/* 右侧内容区样式（灰色底色） */
.admin-content {
    flex: 1; /* 占满剩余宽度 */
    background-color: #f9fafb; /* 灰色底色 */
    padding: 10px; /* 内边距，避免内容贴边 */
    overflow-y: auto; /* 内容超出时仅右侧出现垂直滚动条 */
    height: 100%; /* 强制占满主体区高度 */
}

/* 4. 修复移动端菜单宽度：避免菜单过宽导致横向滚动 */
@media (max-width: 768px) {
    .admin-sidebar {
        width: auto; /* 移动端缩小菜单宽度，适配小屏幕 */
    }
    .system-name {
        font-size: 14px; /* 缩小页眉文字，避免换行 */
    }
    .admin-content {
        background-color: #fff; /* 设置白底 */
        padding: 0px; /* 取消边距 */
    }
}
</style>
