<template>
    <div class="category-page">
        <!-- 顶部导航 -->
        <van-nav-bar fixed title="分类管理" left-arrow @click-left="$router.back()" />

        <!-- 分类内容 -->
        <div class="category-content" v-loading="loading">
            <!-- 一级分类 -->
            <div class="category-row">
                <div
                    v-for="item in parentList"
                    :key="item.id"
                    class="cate-item"
                    :style="{
                        backgroundColor: activeParentId === item.id ? '#1989fa' : '#f2f3f5',
                        color: activeParentId === item.id ? '#fff' : '#333',
                    }"
                    @click="selectParent(item)">
                    {{ item.name }}
                </div>
            </div>

            <!-- 二级分类 -->
            <div class="category-row" v-if="activeSubList.length > 0">
                <div
                    v-for="sub in activeSubList"
                    :key="sub.id"
                    class="cate-item"
                    :style="{
                        backgroundColor: activeSubId === sub.id ? '#1989fa' : '#e8f4ff',
                        color: activeSubId === sub.id ? '#fff' : '#1989fa',
                    }"
                    @click="selectSub(sub)">
                    {{ sub.name }}
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { Getcategory } from "../../api/basic.js";
import { withDelay } from "../../utils/common.js";

import { NavBar as VanNavBar } from "vant";

export default {
    name: "CategoryMob",
    components: {
        VanNavBar,
    },
    data() {
        return {
            loading: true,
            category: [],
            parentList: [],
            activeParentId: "",
            activeSubList: [],
            activeSubId: "",
        };
    },
    methods: {
        async loadGetCategory() {
            this.loading = true;
            const params = { direction: 2 };

            try {
                const res = await withDelay(() => Getcategory(params));
                console.log("接口返回分类：", res.payload?.items); // 调试看数据

                // 一级分类 = 接口返回的顶层数据
                this.parentList = res.payload?.items || [];

                // 默认选中第一个
                if (this.parentList.length) {
                    this.selectParent(this.parentList[0]);
                }
            } catch (err) {
                console.error("加载失败", err);
            } finally {
                this.loading = false;
            }
        },

        // 选中一级分类 → 直接取 children
        selectParent(item) {
            this.activeParentId = item.id;
            this.activeSubId = "";
            this.activeSubList = item.children || []; // ✅ 核心修复
        },

        selectSub(sub) {
            this.activeSubId = sub.id;
            console.log("选中二级：", sub);
        },

        onRefresh() {
            this.loadGetCategory();
        },
    },
    created() {
        this.$globalBus.emit("updateActivePath", "/category");
        this.onRefresh();
    },
};
</script>

<style scoped lang="less">
.category-page {
    background: #f7f8fa;
    min-height: 100vh;
    padding: 46px 0 20px;
}

.category-content {
    padding: 15px;
    margin: 15px;
    border-radius: 8px;
    background-color: #ffffff;
}

.category-row {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-bottom: 10px;
}

.cate-item {
    flex: 1;
    min-width: 60px;
    padding: 10px 4px;
    border-radius: 6px;
    font-size: 14px;
    text-align: center;
    box-sizing: border-box;
    cursor: pointer;
}
</style>
