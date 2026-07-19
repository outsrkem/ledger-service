<template>
    <div class="details-container">
        <van-nav-bar title="账单详情" left-arrow @click-left="$router.back()" right-text="刷新" @click-right="refreshData" />

        <div class="card-box">
            <div class="title">{{ categoryName }}</div>

            <div class="price">{{ detailData.amount }}</div>
            <div class="label">时间：{{ formatTime(detailData.occ_time) }}</div>
            <div class="label">备注：{{ detailData.remark || "" }}</div>

            <!-- 明细：名称 + 严格对齐 -->
            <div class="detail" v-if="detailData.detail && detailData.detail.length > 0">
                明细
                <div class="detail-item" v-for="(item, index) in detailData.detail" :key="index">
                    <!-- 商品名称 -->
                    <span class="d-name">{{ item.name }}</span>

                    <!-- 对齐区域：等宽布局，符号全部对齐 -->
                    <span class="d-align-box">
                        <span class="part">{{ item.quantity }}{{ item.unit }}</span>
                        <span class="sign">×</span>
                        <span class="part">{{ item.price }}元</span>
                        <span class="sign">=</span>
                        <span class="part total">{{ item.total }}元</span>
                    </span>
                </div>
            </div>

            <div class="divider"></div>

            <div class="btn-group">
                <button class="btn delete-btn" @click="handleDelete">删除</button>
                <button class="btn edit-btn">编辑</button>
            </div>
        </div>
    </div>
</template>

<script>
import dayjs from "dayjs";
import { showConfirmDialog, showToast } from "vant";
import { GetBillDetails, DelTransactions, Getcategory } from "../../api/basic.js";
import { getCategoryPath } from "../../utils/category.js";

export default {
    name: "MobileBillDetails",
    data() {
        return {
            billId: "",
            detailData: {
                cid: "",
                amount: "",
                occ_time: "",
                remark: "",
                detail: [],
            },
            categoryName: "加载中...",
        };
    },
    mounted() {
        this.billId = this.$route.params.id;
        this.getBillDetail();
    },
    methods: {
        formatTime(time) {
            if (!time) return "";
            return dayjs(time).format("YYYY-MM-DD HH:mm:ss");
        },

        async getCategoryList() {
            let categoryList = [];
            try {
                categoryList = JSON.parse(localStorage.getItem("category")) || [];
            } catch (e) {
                console.log(e);
            }

            if (categoryList && categoryList.length > 0) {
                return categoryList;
            }

            try {
                const [res1, res2] = await Promise.all([Getcategory({ direction: 2 }), Getcategory({ direction: 1 })]);
                categoryList = [...(res1.payload?.items || []), ...(res2.payload?.items || [])];
                localStorage.setItem("category", JSON.stringify(categoryList));
            } catch (e) {
                console.error("加载分类失败", e);
            }
            return categoryList;
        },

        async getBillDetail() {
            try {
                const res = await GetBillDetails(this.billId);
                this.detailData = res.payload || {};
                const categoryList = await this.getCategoryList();
                const name = getCategoryPath(categoryList, this.detailData.cid);
                this.categoryName = name || "未分类";
            } catch (err) {
                console.error("获取详情失败：", err);
                this.categoryName = "获取失败";
            }
        },

        async refreshData() {
            try {
                localStorage.removeItem("category");
                showToast("刷新中...");
                await this.getBillDetail();
                showToast("刷新成功");
            } catch {
                showToast("刷新失败");
            }
        },

        async handleDelete() {
            try {
                await showConfirmDialog({
                    title: "确认删除",
                    message: "确定要删除这条账单记录吗？删除后无法恢复",
                    confirmButtonText: "删除",
                    cancelButtonText: "取消",
                    confirmButtonColor: "#ee0a24",
                });
                await DelTransactions(this.billId);
                this.$router.back();
            } catch {
                showToast("刷新失败");
            }
        },
    },
};
</script>

<style lang="less" scoped>
.details-container {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: #e9ebee;
}
.card-box {
    background: #fff;
    border-radius: 8px;
    margin: 16px;
    padding: 20px;
}
.title {
    text-align: center;
    font-size: 16px;
    margin-top: 50px;
    margin-bottom: 12px;
}
.price {
    text-align: center;
    font-size: 24px;
    font-weight: bold;
    color: #333;
    margin-top: 20px;
    margin-bottom: 20px;
}
.label {
    font-size: 14px;
    color: #666;
    text-align: left;
    margin-bottom: 20px;
}

/* 明细布局 - 超强对齐 */
.detail {
    font-size: 14px;
    color: #666;
    margin-bottom: 20px;

    .detail-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 8px;
        color: #333;
    }

    .d-name {
        font-weight: 500;
        flex-shrink: 0;
        margin-right: 12px;
    }

    /* 对齐核心容器 */
    .d-align-box {
        display: flex;
        align-items: center;
        gap: 6px;
        flex-shrink: 0;
    }

    /* 每一段固定宽度 */
    .part {
        display: inline-block;
        width: 40px;
        text-align: right;
    }

    /* 符号居中 */
    .sign {
        width: 12px;
        text-align: center;
    }

    .total {
        width: 50px;
        font-weight: 500;
    }
}

.divider {
    height: 1px;
    background-color: #f0f0f0;
    margin: 10px 0 20px;
}
.btn-group {
    display: flex;
    gap: 12px;
    margin-top: 10px;
}
.btn {
    flex: 1;
    height: 44px;
    border-radius: 6px;
    font-size: 14px;
    border: none;
    cursor: pointer;
}
.delete-btn {
    background-color: #f5f5f5;
    color: #666;
}
.edit-btn {
    background-color: #409eff;
    color: #fff;
}
</style>
