<template>
    <el-card style="width: 100%">
        <template #header>
            <div class="my_refresh" style="display: flex; justify-content: space-between; align-items: center">
                <el-row>
                    <span>分类管理</span>
                </el-row>
                <el-row>
                    <el-button type="primary" :icon="Refresh" @click="onRefresh" :loading="loading">刷新</el-button>
                </el-row>
            </div>
        </template>
        <div style="display: flex; justify-content: flex-start; gap: 24px">
            <!-- 收入分类 direction:1 -->
            <div>
                <div style="display: flex; justify-content: center; padding: 12px">
                    <el-text style="font-weight: 500">收入分类</el-text>
                </div>
                <div style="border-radius: 8px; border: 1px dashed #aaa; padding: 8px; width: 400px">
                    <el-tree v-loading="loading" :data="inCategoryList" :props="defaultProps" @node-click="handleNodeClick"></el-tree>
                </div>
            </div>
            <!-- 支出分类 direction:2 -->
            <div>
                <div style="display: flex; justify-content: center; padding: 12px">
                    <el-text style="font-weight: 500">支出分类</el-text>
                </div>
                <div style="border-radius: 8px; border: 1px dashed #aaa; padding: 8px; width: 400px">
                    <el-tree v-loading="loading" :data="outCategoryList" :props="defaultProps" :expand-on-click-node="false" @node-click="handleNodeClick">
                        <template #default="{ node, data }">
                            <div class="custom-tree-node">
                                <span>{{ node.label }}</span>
                                <div>
                                    <el-button type="primary" link @click.stop="append(data)"> 新增 </el-button>
                                    <el-button style="margin-left: 4px" type="danger" link @click.stop="remove(node, data)"> 删除 </el-button>
                                </div>
                            </div>
                        </template>
                    </el-tree>
                </div>
            </div>
        </div>
    </el-card>
</template>

<script>
import { Getcategory } from "../../api/basic.js";
import { withDelay } from "../../utils/common.js";
import { Refresh } from "@element-plus/icons-vue";

export default {
    name: "CategoryIndex",
    setup() {
        return {
            Refresh,
        };
    },
    data() {
        return {
            loading: false,
            // 支出分类数据
            outCategoryList: [],
            // 收入分类数据
            inCategoryList: [],
            defaultProps: {
                label: "name",
            },
        };
    },
    methods: {
        // 通用加载分类
        async loadCategory(direction) {
            const params = { direction };
            const res = await withDelay(() => Getcategory(params));
            return res.payload?.items || [];
        },

        // 刷新全部分类
        async onRefresh() {
            this.loading = true;
            try {
                const [inList, outList] = await Promise.all([this.loadCategory(1), this.loadCategory(2)]);
                this.inCategoryList = inList;
                this.outCategoryList = outList;
            } catch (error) {
                console.error("加载分类失败：", error);
            } finally {
                this.loading = false;
            }
        },

        handleNodeClick() {
            // 预留树节点点击事件
        },
        append() {
            console.log("----");
        },
        remove() {
            console.log("--remove--");
        },
    },
    created() {
        this.$globalBus.emit("updateActivePath", "/category");
        this.onRefresh();
    },
};
</script>
<style scoped lang="less">
.custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    padding-right: 8px;
}
</style>
