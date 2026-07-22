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
            <!-- 支出分类 direction:2 -->
            <div>
                <div style="display: flex; justify-content: center; padding: 12px">
                    <el-text style="font-weight: 500">支出分类</el-text>
                </div>
                <!-- 添加固定高度与滚动 -->
                <div style="border-radius: 8px; border: 1px dashed #aaa; padding: 8px; width: 400px; max-height: 450px; overflow-y: auto">
                    <el-tree
                        ref="outTreeRef"
                        v-loading="loading"
                        :data="outCategoryList"
                        :default-expanded-keys="expandedKeys"
                        :props="defaultProps"
                        :expand-on-click-node="false"
                        @node-click="handleNodeClick">
                        <template #default="{ node, data }">
                            <div class="custom-tree-node">
                                <span>{{ node.label }}</span>
                                <div>
                                    <!-- 根节点level=1 显示新增主分类 -->
                                    <el-button v-if="node.level === 1" type="primary" link @click.stop="appendMain(data)"> 新增主分类 </el-button>
                                    <!-- 一级子节点level=2 显示新增子分类 -->
                                    <el-button v-if="node.level === 2" type="primary" link @click.stop="appendSub(data)"> 新增子分类 </el-button>
                                    <!-- 根节点不显示删除；level2、level3显示删除 -->
                                    <el-button
                                        v-if="node.level === 2 || node.level === 3"
                                        style="margin-left: 4px"
                                        type="danger"
                                        link
                                        @click.stop="remove(node, data)">
                                        删除
                                    </el-button>
                                </div>
                            </div>
                        </template>
                    </el-tree>
                </div>
            </div>
            <!-- 收入分类 direction:1 -->
            <div>
                <div style="display: flex; justify-content: center; padding: 12px">
                    <el-text style="font-weight: 500">收入分类</el-text>
                </div>
                <!-- 添加固定高度与滚动 -->
                <div style="border-radius: 8px; border: 1px dashed #aaa; padding: 8px; width: 400px; max-height: 450px; overflow-y: auto">
                    <el-tree
                        ref="inTreeRef"
                        v-loading="loading"
                        :data="inCategoryList"
                        :default-expanded-keys="expandedKeys"
                        :props="defaultProps"
                        @node-click="handleNodeClick">
                        <template #default="{ node, data }">
                            <div class="custom-tree-node">
                                <span>{{ node.label }}</span>
                                <div>
                                    <!-- 根节点level=1 显示新增主分类 -->
                                    <el-button v-if="node.level === 1" type="primary" link @click.stop="appendMain(data)"> 新增主分类 </el-button>
                                    <!-- 一级子节点level=2 显示新增子分类 -->
                                    <el-button v-if="node.level === 2" type="primary" link @click.stop="appendSub(data)"> 新增子分类 </el-button>
                                    <!-- 根节点不显示删除；level2、level3显示删除 -->
                                    <el-button
                                        v-if="node.level === 2 || node.level === 3"
                                        style="margin-left: 4px"
                                        type="danger"
                                        link
                                        @click.stop="remove(node, data)">
                                        删除
                                    </el-button>
                                </div>
                            </div>
                        </template>
                    </el-tree>
                </div>
            </div>
        </div>
    </el-card>

    <el-dialog v-model="createCategoryDialog" title="新建分类" width="500">
        <span>输入分类名称：</span>
        <span v-if="parentName">{{ parentName }} - </span>
        <el-input v-model="name" style="width: 240px" placeholder="Please input" />
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="createCategoryDialog = false">取消</el-button>
                <el-button type="primary" @click="onCreateCategory"> 确定 </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script>
import { Getcategory, DelCategory, AddSubCategory, AddMainCategory } from "../../api/basic.js";
import { withDelay } from "../../utils/common.js";
import { Refresh } from "@element-plus/icons-vue";
import { ElMessageBox } from "element-plus";
import { msgcon } from "../../utils/message.js";

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
                children: "children",
            },
            createCategoryDialog: false,
            parentName: "",
            name: "",
            // 当前操作的数据（用于区分新增主分类还是子分类）
            currentData: null,
            // 当前操作类型: 'main' | 'sub'
            currentActionType: "main",
            // 根节点ID标记
            rootInId: "root-in",
            rootOutId: "root-out",
            // 默认只展开根节点
            expandedKeys: [],
            // 树引用
            inTreeRef: null,
            outTreeRef: null,
        };
    },
    watch: {
        // 监听数据变化，确保展开根节点
        inCategoryList: {
            handler(newVal) {
                if (newVal && newVal.length > 0) {
                    this.$nextTick(() => {
                        this.expandRootNode("inTreeRef");
                    });
                }
            },
            deep: true,
            immediate: true,
        },
        outCategoryList: {
            handler(newVal) {
                if (newVal && newVal.length > 0) {
                    this.$nextTick(() => {
                        this.expandRootNode("outTreeRef");
                    });
                }
            },
            deep: true,
            immediate: true,
        },
    },
    methods: {
        // 展开根节点方法
        expandRootNode(refName) {
            try {
                const treeRef = this.$refs[refName];
                if (!treeRef) return;

                // 获取树的所有节点
                const treeInstance = treeRef;
                if (treeInstance.store && treeInstance.store.root) {
                    const rootNode = treeInstance.store.root;
                    // 展开根节点
                    if (rootNode.childNodes && rootNode.childNodes.length > 0) {
                        rootNode.childNodes.forEach((node) => {
                            node.expanded = true;
                        });
                    }
                }
            } catch (error) {
                console.warn("展开根节点失败:", error);
            }
        },

        // 通用加载分类 - 自动包裹根节点
        async loadCategory(direction) {
            const params = { direction };
            const res = await withDelay(() => Getcategory(params));
            const list = res.payload?.items || [];

            // 收入分类 - 包裹自定义根节点
            if (direction === 1) {
                return [
                    {
                        id: this.rootInId,
                        name: "收入分类",
                        children: list,
                    },
                ];
            }
            // 支出分类 - 包裹自定义根节点
            if (direction === 2) {
                return [
                    {
                        id: this.rootOutId,
                        name: "支出分类",
                        children: list,
                    },
                ];
            }
            return list;
        },

        // 刷新全部分类
        async onRefresh() {
            this.loading = true;
            try {
                const [inList, outList] = await Promise.all([this.loadCategory(1), this.loadCategory(2)]);
                this.inCategoryList = inList;
                this.outCategoryList = outList;
                // 设置展开根节点的keys
                this.expandedKeys = [this.rootInId, this.rootOutId];

                // 等待DOM更新后强制展开根节点
                await this.$nextTick();
                this.expandRootNode("inTreeRef");
                this.expandRootNode("outTreeRef");
            } catch (error) {
                console.error("加载分类失败：", error);
                this.$message.error("加载分类列表失败");
            } finally {
                this.loading = false;
            }
        },

        handleNodeClick() {
            // 预留树节点点击事件
        },

        // 打开新增主分类弹窗（根节点调用）
        appendMain(data) {
            this.currentData = data;
            this.currentActionType = "main";
            this.parentName = data.name;
            this.name = "";
            this.createCategoryDialog = true;
        },

        // 打开新增子分类弹窗（一级子节点调用）
        appendSub(data) {
            this.currentData = data;
            this.currentActionType = "sub";
            this.parentName = data.name;
            this.name = "";
            this.createCategoryDialog = true;
        },

        // 新增分类
        async onCreateCategory() {
            if (!this.name || !this.name.trim()) {
                return this.$message.warning(msgcon("请输入分类名称"));
            }

            const trimmedName = this.name.trim();

            try {
                if (this.currentActionType === "main") {
                    // 新增主分类
                    // 根据当前点击的根节点判断 direction
                    // 收入分类根节点ID为 rootInId，支出分类根节点ID为 rootOutId
                    let direction = 1; // 默认收入
                    if (this.currentData.id === this.rootOutId) {
                        direction = 2; // 支出
                    } else if (this.currentData.id === this.rootInId) {
                        direction = 1; // 收入
                    }

                    await withDelay(() =>
                        AddMainCategory({
                            name: trimmedName,
                            direction: direction,
                        }),
                    );
                    this.$message.success(msgcon("新增主分类成功"));
                } else if (this.currentActionType === "sub") {
                    // 新增子分类
                    await withDelay(() => AddSubCategory({ id: this.currentData.id }, { name: trimmedName }));
                    this.$message.success(msgcon("新增子分类成功"));
                }

                // 关闭弹窗
                this.createCategoryDialog = false;
                this.name = "";
                this.currentData = null;

                // 刷新列表
                await this.onRefresh();
            } catch (error) {
                console.error("新增分类失败：", error);
                this.$message.error(msgcon("新增分类失败，请稍后重试"));
            }
        },

        // 删除功能
        remove(node, data) {
            // 禁止删除根节点
            if (data.id === this.rootInId || data.id === this.rootOutId) {
                return this.$message.warning(msgcon("根节点禁止删除"));
            }
            // 存在子分类禁止删除
            if (node.childNodes && node.childNodes.length > 0) {
                return this.$message.warning(msgcon("该分类下存在子分类，无法删除"));
            }

            ElMessageBox.confirm(`确定要删除【${data.name}】分类吗？`, "删除提示", {
                type: "warning",
            })
                .then(async () => {
                    try {
                        await withDelay(() => DelCategory(data.id));
                        this.$message.success(msgcon("删除成功"));
                        await this.onRefresh();
                    } catch (err) {
                        console.error("删除分类失败：", err);
                        this.$message.error(msgcon("删除失败，请稍后重试"));
                    }
                })
                .catch(() => {
                    // 取消删除不处理
                });
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
