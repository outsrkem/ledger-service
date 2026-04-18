<template>
    <div>
        <div class="card-box">
            <div class="text">
                <div style="width: 50%">
                    <div style="font-size: 12px; color: #999; margin-bottom: 8px">支出(元)</div>
                    <div style="font-size: 16px; font-weight: bold">{{ expense }}</div>
                </div>
                <div style="width: 50%">
                    <div style="font-size: 12px; color: #999; margin-bottom: 8px">收入(元)</div>
                    <div style="font-size: 16px; font-weight: bold">{{ income }}</div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { msgcon } from "../../utils/message.js";
import { GetStatistic } from "../../api/basic.js";

export default {
    name: "BillStatistic",
    data() {
        return {
            income: "0.00",
            expense: "0.00",
            pt: "month", // 默认查本月
        };
    },
    mounted() {
        this.getBillStatistic();
    },
    methods: {
        async getBillStatistic() {
            try {
                // 1. 获取当前年月 2026-03
                const now = new Date();
                const year = now.getFullYear();
                const month = (now.getMonth() + 1).toString().padStart(2, "0");
                const monthStr = `${year}-${month}`;

                // 2. 调用接口
                const res = await GetStatistic({
                    sa: "s01",
                    pt: this.pt,
                    from: monthStr,
                    to: monthStr,
                });

                // 3. 按你的接口格式取值 payload
                if (res?.payload) {
                    this.income = res.payload.income || "0.00";
                    this.expense = res.payload.expense || "0.00";
                }
            } catch (err) {
                console.error("获取统计失败：", err);
                msgcon.error("获取收支统计失败");
            }
        },
    },
};
</script>

<style lang="less" scoped>
.card-box {
    background: #fff;
    border-radius: 8px;
    margin: 16px;
}
.text {
    font-size: 14px;
    color: #333;
    padding: 20px 20px;
    display: flex;
}
</style>
