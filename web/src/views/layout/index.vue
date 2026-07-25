<template>
    <div class="main-container">
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
export default {
    name: "LayoutIndex",
    components: { AppAside },
};
</script>

<style scoped lang="less">
.main-container {
    min-width: 1366px;
    height: 100vh;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    overflow: visible;
    padding-top: 50px;
}
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
