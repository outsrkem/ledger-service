<template>
    <el-drawer v-model="drawer" direction="rtl" title="账单详情" :size="'35%'" @close="handleClose">
        <div v-loading="loading" class="detail-wrap">
            <!-- 主金额区块 -->
            <div class="amount-card">
                <div class="amount-label">账单总金额</div>
                <div :class="amountClass">{{ amountPrefix }}{{ disFrom.amount }}</div>
            </div>

            <!-- 基础信息 -->
            <div class="info-card">
                <div class="info-row">
                    <span class="label">记账时间</span>
                    <span class="value">{{ disFrom.occ_time || "-" }}</span>
                </div>
                <div class="info-row">
                    <span class="label">分类</span>
                    <span class="value">{{ disFrom.categoryName || "-" }}</span>
                </div>
                <div class="info-row">
                    <span class="label">备注</span>
                    <span class="value">{{ disFrom.remark || "无备注" }}</span>
                </div>
            </div>

            <!-- 明细列表 -->
            <div class="detail-card" v-if="disFrom.detail.length > 0">
                <div class="card-title">账目明细</div>
                <div class="detail-item" v-for="(val, inx) in disFrom.detail" :key="inx">
                    <div class="item-name">{{ val.name }}</div>
                    <div class="item-calc">
                        {{ val.price }}元 × {{ val.quantity }}{{ val.unit }} =
                        <span class="item-total">{{ val.total }}</span>
                    </div>
                </div>
            </div>
            <div class="empty-tip" v-else>暂无明细</div>
        </div>
    </el-drawer>
</template>

<script>
import dayjs from "dayjs";
import { GetTranDetails } from "../../api/basic.js";

const DEFAULT_TRAN_FORM = {
    occ_time: "",
    amount: "",
    category: "",
    categoryName: "",
    remark: "",
    detail: [
        {
            name: "",
            quantity: null,
            price: null,
            unit: "",
            total: null,
        },
    ],
};

export default {
    name: "TranDetail",
    data() {
        return {
            drawer: false,
            loading: true,
            tranId: "",
            tranRawAmount: 0, // 原始带正负金额，用来判断收支
            disFrom: { ...DEFAULT_TRAN_FORM },
        };
    },
    computed: {
        amountPrefix() {
            return this.tranRawAmount >= 0 ? "+" : "-";
        },
        amountClass() {
            return this.tranRawAmount >= 0 ? "income-text" : "expense-text";
        },
    },
    methods: {
        async onOpenDialog(val) {
            this.tranRawAmount = 0;
            this.loading = true;
            this.drawer = true;
            this.tranId = val.id;

            try {
                await this.GetDetails();
            } finally {
                this.loading = false;
            }
        },

        async GetDetails() {
            const res = await GetTranDetails(this.tranId);
            const payload = res.payload || {};

            this.tranRawAmount = Number(payload.amount);
            this.disFrom.occ_time = dayjs(payload.occ_time).format("YYYY-MM-DD HH:mm:ss");
            this.disFrom.amount = Math.abs(this.tranRawAmount).toString();
            this.disFrom.category = payload.cid;
            this.disFrom.categoryName = payload.categoryName || ""; // 账单分类名称
            this.disFrom.remark = payload.remark || "";
            this.disFrom.detail = payload.detail ? [...payload.detail] : [];
        },

        // 关闭抽屉重置数据，防止缓存旧数据
        handleClose() {
            this.tranId = "";
            // this.tranRawAmount = 0;
        },
    },
};
</script>

<style scoped lang="less">
.detail-wrap {
    padding: 4px 8px;
}

.amount-card {
    background: #f7f8fa;
    border-radius: 8px;
    padding: 16px 20px;
    margin-bottom: 16px;
    text-align: center;
    .amount-label {
        font-size: 14px;
        color: #666;
        margin-bottom: 6px;
    }
    .income-text {
        font-size: 30px;
        font-weight: bold;
        color: #f53f3f;
    }
    .expense-text {
        font-size: 30px;
        font-weight: bold;
        color: #00b42a;
    }
}

.info-card {
    border: 1px solid #e5e6eb;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
    .info-row {
        display: flex;
        padding: 8px 0;
        &:not(:last-child) {
            border-bottom: 1px solid #f2f3f5;
        }
        .label {
            width: 90px;
            color: #666;
            flex-shrink: 0;
        }
        .value {
            flex: 1;
            color: #222;
            word-break: break-all;
        }
    }
}

.detail-card {
    border: 1px solid #e5e6eb;
    border-radius: 8px;
    padding: 16px;
    .card-title {
        font-weight: bold;
        font-size: 15px;
        margin-bottom: 12px;
    }
    .detail-item {
        display: flex;
        justify-content: space-between;
        padding: 8px 0;
        &:not(:last-child) {
            border-bottom: 1px dashed #eee;
        }
        .item-name {
            color: #333;
        }
        .item-calc {
            color: #666;
            .item-total {
                color: #1890ff;
                font-weight: 500;
            }
        }
    }
}

.empty-tip {
    color: #999;
    text-align: center;
    padding: 30px 0;
    border: 1px dashed #e5e6eb;
    border-radius: 8px;
}
</style>
