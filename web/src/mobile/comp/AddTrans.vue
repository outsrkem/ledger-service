<template>
    <div>
        <van-popup v-model:show="showBottom" closeable round position="bottom" :style="{ height: '85%' }">
            <div style="margin: 40px 15px 0 15px; padding-bottom: 20px">
                <!-- 顶部：收支 + 时间 -->
                <div style="display: flex; justify-content: space-between">
                    <div>
                        <van-button
                            round
                            size="mini"
                            style="width: 50px; margin-right: 5px"
                            :type="formdata.direction === 2 ? 'primary' : 'default'"
                            @click="switchDirection(2)"
                        >
                            支出
                        </van-button>
                        <van-button
                            round
                            size="mini"
                            style="width: 50px"
                            :type="formdata.direction === 1 ? 'primary' : 'default'"
                            @click="switchDirection(1)"
                        >
                            收入
                        </van-button>
                    </div>
                    <div>
                        <van-button round type="default" size="mini" style="min-width: 50px" @click="showDatePicker = true">
                            {{ dateText }}
                        </van-button>
                        <van-button round type="default" size="mini" style="min-width: 50px">
                            {{ currentTime }}
                        </van-button>
                    </div>
                </div>

                <!-- 金额输入 -->
                <div style="margin: 10px 0; display: flex; align-items: center; gap: 8px">
                    <input
                        style="font-size: 18px; padding: 8px; border: 1px solid #eee; border-radius: 6px"
                        placeholder="请输入金额"
                        v-model="formdata.amount"
                    />
                    <van-button type="default" style="height: 39px; font-size: 16px">添加明细</van-button>
                </div>

                <!-- 分类：水平撑满、均分 -->
                <div v-loading="loading" style="margin-top: 10px">
                    <template v-for="(row, rowIdx) in categoryRows" :key="rowIdx">
                        <div class="category-row">
                            <div
                                v-for="item in row"
                                :key="item.id"
                                class="cate-item"
                                :style="{
                                    backgroundColor: activeParentId === item.id ? '#1989fa' : '#f2f3f5',
                                    color: activeParentId === item.id ? '#fff' : '#333',
                                }"
                                @click="selectParent(item, rowIdx)"
                            >
                                {{ item.name }}
                            </div>
                        </div>

                        <div v-if="activeRowIndex === rowIdx && activeSubList.length" class="category-row">
                            <div
                                v-for="sub in activeSubList"
                                :key="sub.id"
                                class="cate-item"
                                :style="{
                                    backgroundColor: activeSubId === sub.id ? '#1989fa' : '#e8f4ff',
                                    color: activeSubId === sub.id ? '#fff' : '#1989fa',
                                }"
                                @click="selectSub(sub)"
                            >
                                {{ sub.name }}
                            </div>
                        </div>
                    </template>
                </div>

                <!-- 备注 -->
                <van-field v-model="formdata.remark" placeholder="写点备注..." style="margin-top: 15px" />

                <!-- 按钮 -->
                <div style="display: flex; gap: 10px; margin-top: 20px">
                    <van-button type="default" block @click="showBottom = false">取消</van-button>
                    <van-button type="primary" block @click="onSubmit(0)">保存</van-button>
                </div>
            </div>
        </van-popup>

        <!-- 日期选择 -->
        <van-calendar v-model:show="showDatePicker" :min-date="new Date(2020, 0, 1)" :row-height="40" @confirm="onConfirmDate" />
    </div>
</template>

<script>
import { Getcategory, AddTransactions } from "../../api/basic.js";
import dayjs from "dayjs";

export default {
    name: "AddTransIndex",
    data() {
        return {
            showBottom: false,
            showDatePicker: false,
            currentTime: dayjs().format("HH:mm:ss"),
            selectDate: dayjs().format("YYYY-MM-DD"),

            formdata: {
                direction: 2,
                amount: "",
                cid: "",
                occ_time: dayjs().format("YYYY-MM-DDTHH:mm:ssZZ"),
                remark: "",
            },

            categoryList: [],
            categoryRows: [],
            activeRowIndex: -1,
            activeParentId: null,
            activeSubList: [],
            activeSubId: null,
            loading: false,
        };
    },
    computed: {
        dateText() {
            const now = dayjs().format("YYYY-MM-DD");
            return this.selectDate === now ? "今天" : this.selectDate;
        },
    },
    methods: {
        onOpenDialog() {
            this.formdata = {
                direction: 2,
                amount: "",
                cid: "",
                occ_time: dayjs().format("YYYY-MM-DDTHH:mm:ssZZ"),
                remark: "",
            };
            this.selectDate = dayjs().format("YYYY-MM-DD");
            this.resetAll();
            this.showBottom = true;
            this.loadCategory();
        },

        onConfirmDate(date) {
            const d = dayjs(date);
            this.selectDate = d.format("YYYY-MM-DD");
            this.formdata.occ_time = d.format("YYYY-MM-DD") + "T" + this.currentTime + "+0800";
            this.showDatePicker = false;
        },

        async switchDirection(val) {
            this.formdata.direction = val;
            this.resetAll();
            await this.loadCategory();
        },

        async loadCategory() {
            this.loading = true;
            const res = await Getcategory({ direction: this.formdata.direction });
            this.categoryList = res.payload?.items || [];
            const rows = [];
            for (let i = 0; i < this.categoryList.length; i += 5) {
                rows.push(this.categoryList.slice(i, i + 5));
            }
            this.categoryRows = rows;
            this.loading = false;
        },

        selectParent(parent, rowIdx) {
            this.activeParentId = parent.id;
            this.activeRowIndex = rowIdx;
            this.activeSubList = parent.children || [];
            this.activeSubId = null;
            this.formdata.cid = parent.id;
        },

        selectSub(sub) {
            this.activeSubId = sub.id;
            this.formdata.cid = sub.id;
        },

        async onSubmit(keepOpen) {
            if (!this.formdata.occ_time) {
                return this.$message.error("请选择时间");
            }
            if (!this.formdata.amount) {
                return this.$message.error("请输入金额");
            }
            if (!this.formdata.cid) {
                return this.$message.error("请选择分类");
            }

            let amount = Number(this.formdata.amount);
            if (this.formdata.direction === 2) {
                amount = -Math.abs(amount);
            }

            const data = {
                cid: this.formdata.cid,
                occ_time: this.formdata.occ_time,
                amount: amount,
                remark: this.formdata.remark,
                detail: [],
            };

            try {
                await AddTransactions(data);
                this.$message.success("添加成功");
                this.$globalBus.emit("onRefresh");

                if (keepOpen) {
                    this.formdata.amount = "";
                    this.formdata.cid = "";
                    this.formdata.remark = "";
                    this.formdata.occ_time = dayjs().format("YYYY-MM-DDTHH:mm:ssZZ");
                    this.selectDate = dayjs().format("YYYY-MM-DD");
                    this.resetAll();
                } else {
                    this.showBottom = false;
                }
            } catch (err) {
                const msg = err.data?.metadata?.message || "提交失败";
                this.$message.error(msg);
            }
        },

        resetAll() {
            this.activeRowIndex = -1;
            this.activeParentId = null;
            this.activeSubList = [];
            this.activeSubId = null;
        },
    },
};
</script>

<style scoped>
.category-row {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-bottom: 8px;
}

/* 关键：水平撑满、均分宽度 */
.cate-item {
    flex: 1;
    min-width: 60px;
    padding: 8px 4px;
    border-radius: 6px;
    font-size: 13px;
    text-align: center;
    white-space: nowrap;
    box-sizing: border-box;
}
</style>
