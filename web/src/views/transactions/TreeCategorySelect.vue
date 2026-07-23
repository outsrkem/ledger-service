<template>
    <div v-loading="loading">
        <!-- 一级分类 整行横向排列（固定不会被撑开） -->
        <div class="category-row parent-row">
            <div
                v-for="item in categoryList"
                :key="item.id"
                class="cate-item parent-cate-item"
                :class="{ 'active-parent': activeParentId === item.id }"
                @click="selectParent(item)">
                {{ item.name }}
                <span v-if="item.children && item.children.length" class="dot-mark">•</span>
            </div>
        </div>
        <div class="sub-row">
            <div v-if="activeSubList.length" class="category-row child-row line">
                <div v-for="sub in activeSubList" :key="sub.id" class="cate-item" :class="{ 'active-sub': activeSubId === sub.id }" @click="selectSub(sub)">
                    {{ sub.name }}
                </div>
            </div>

            <div v-else-if="activeParentId !== null" class="empty-tip">该分类暂无子类</div>
            <div v-else class="empty-tip">请选择分类</div>
        </div>
    </div>
</template>

<script>
import { Getcategory } from "@/api/basic.js";

export default {
    name: "CategoryPcAutoSelect",
    props: {
        modelValue: {
            type: [Number, String, Object],
            default: "",
        },
        direction: {
            type: Number,
            required: true,
        },
    },
    data() {
        return {
            categoryList: [],
            activeParentId: null,
            activeSubList: [],
            activeSubId: null,
            loading: false,
        };
    },
    watch: {
        modelValue: {
            handler(val) {
                this.$nextTick(() => this.matchSelected(val));
            },
            immediate: true,
        },
        direction: {
            handler() {
                this.loadCategory();
                this.resetState();
            },
            immediate: true,
        },
    },
    methods: {
        async loadCategory() {
            this.loading = true;
            try {
                const res = await Getcategory({ direction: this.direction });
                this.categoryList = res.payload?.items || [];
            } catch (err) {
                console.error("加载分类失败", err);
            } finally {
                this.loading = false;
            }
        },

        matchSelected(targetCid) {
            if (!targetCid) {
                this.resetState();
                return;
            }
            let targetParent = null;
            let targetChild = null;

            for (const parent of this.categoryList) {
                if (parent.id === targetCid) {
                    targetParent = parent;
                    break;
                }
                if (parent.children) {
                    targetChild = parent.children.find((c) => c.id === targetCid);
                    if (targetChild) {
                        targetParent = parent;
                        break;
                    }
                }
            }

            if (!targetParent) return;
            this.activeParentId = targetParent.id;
            this.activeSubList = targetParent.children || [];
            this.activeSubId = targetChild ? targetChild.id : null;
        },

        selectParent(parent) {
            this.activeParentId = parent.id;
            this.activeSubList = parent.children || [];
            this.activeSubId = null;
            this.$emit("update:modelValue", parent.id);
            this.$parent.$emit("el.form.change", parent.id);
        },

        selectSub(sub) {
            this.activeSubId = sub.id;
            this.$emit("update:modelValue", sub.id);
            this.$parent.$emit("el.form.change", sub.id);
        },

        resetState() {
            this.activeParentId = null;
            this.activeSubList = [];
            this.activeSubId = null;
        },

        reset() {
            this.resetState();
            this.$emit("update:modelValue", "");
        },
    },
};
</script>

<style scoped>
.category-row {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-bottom: 6px;

    padding-top: 6px;
}
.line {
    border-top: 1px dashed #dcdfe6;
}
.child-row {
    margin-bottom: 0;
}

.cate-item {
    flex: 0 0 auto; /* 禁止自动拉伸！固定尺寸不自动变宽 */
    min-width: 56px;
    padding: 4px 14px;
    border-radius: 5px;
    font-size: 13px;
    font-weight: 500;
    text-align: center;
    white-space: nowrap;
    cursor: pointer;
    user-select: none;
    transition: all 0.2s ease;
    border: 1px solid transparent;
    background-color: #f1f3f5;
    color: #2c3e50;
    line-height: 1.6;
}

.parent-cate-item {
    position: relative;
}

.dot-mark {
    position: absolute;
    right: 6px;
    bottom: 2px;
    font-size: 10px;
    letter-spacing: 1px;
    pointer-events: none;
}

.parent-row .cate-item {
    background-color: #f1f3f5;
    border-color: #e8ecf0;
    color: #1e293b;
}
.parent-row .cate-item:not(.active-parent) .dot-mark {
    color: #64748b;
}
.parent-row .cate-item:hover {
    background-color: #e6e9ef;
    border-color: #cbd5e1;
}
.parent-row .cate-item.active-parent {
    background-color: #409eff;
    color: #ffffff;
    border-color: #409eff;
    box-shadow: 0 2px 8px rgba(26, 110, 255, 0.25);
}
.parent-row .cate-item.active-parent .dot-mark {
    color: #fff;
}

.child-row .cate-item {
    background-color: #f5f9ff;
    color: #409eff;
    border-color: #dce6f5;
    font-size: 12px;
    padding: 3px 12px;
    min-width: 48px;
}
.child-row .cate-item:hover {
    background-color: #eaf0fe;
    border-color: #b0c8f0;
}
.child-row .cate-item.active-sub {
    background-color: #409eff;
    color: #ffffff;
    border-color: #409eff;
    box-shadow: 0 2px 8px rgba(26, 110, 255, 0.25);
}
.sub-row {
    min-height: 35px;
}
.empty-tip {
    color: #94a3b8;
    font-size: 12px;
    border-top: 1px dashed #e2e8f0;
    margin-top: 2px;
}
</style>
