<template>
    <div style="padding-top: 10px">
        <el-card style="width: 100%">
            <div class="my_refresh">
                <div>
                    <el-space>
                        <el-text>
                            本月收入：
                            <span style="display: inline-block; width: 80px">{{ incomeTotal }}</span>
                        </el-text>
                        <el-text>
                            本月支出：
                            <span style="display: inline-block; width: 80px">{{ expendTotal }}</span>
                        </el-text>
                    </el-space>
                </div>
                <el-space>
                    <el-date-picker
                        v-if="mode === 'year'"
                        v-model="month"
                        type="year"
                        placeholder="选择年份"
                        format="YYYY"
                        value-format="YYYY"
                        style="width: 320px"
                        @change="handleMonthChange" />
                    <el-date-picker
                        v-if="mode === 'month'"
                        v-model="month"
                        type="month"
                        placeholder="选择月份"
                        format="YYYY-MM"
                        value-format="YYYY-MM"
                        style="width: 320px"
                        @change="handleMonthChange" />
                    <el-date-picker
                        v-if="mode === 'custom'"
                        v-model="custom"
                        type="daterange"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        format="YYYY-MM-DD"
                        value-format="YYYY-MM-DD"
                        style="width: 320px"
                        @change="handleMonthChange" />

                    <el-select v-model="mode" placeholder="报表类型" style="width: 240px" @change="handleModeChange">
                        <el-option label="年报" value="year" />
                        <el-option label="月报" value="month" />
                        <el-option label="自定义" value="custom" />
                    </el-select>

                    <el-button type="primary" :icon="Refresh" @click="onRefresh" :loading="loading">刷新</el-button>
                </el-space>
            </div>
        </el-card>
    </div>

    <div style="padding-top: 10px">
        <div class="chart-wrapper">
            <div class="chart-item">
                <Day ref="dayChartRef" :query-params="chartParams" />
            </div>
        </div>
    </div>
</template>

<script>
import { Refresh } from "@element-plus/icons-vue";
import { GetCycleSummary } from "../../api/basic.js";
import Day from "./charts/day.vue";

export default {
    name: "OverviewIndex",
    components: { Day },
    setup() {
        return { Refresh };
    },
    data() {
        return {
            loading: false,
            mode: "month",
            month: "",
            custom: [],
            chartParams: null,
            incomeTotal: 0,
            expendTotal: 0,
        };
    },
    methods: {
        getRangeParams() {
            const { mode, month, custom } = this;
            if (mode === "custom") {
                if (!Array.isArray(custom) || custom.length !== 2) return null;
                return {
                    from: custom[0],
                    to: custom[1],
                };
            }
            if (!month) return null;

            if (mode === "year") {
                return {
                    from: `${month}-01-01`,
                    to: `${month}-12-31`,
                };
            }

            if (mode === "month") {
                const [y, m] = month.split("-");
                const lastDay = new Date(y, m, 0).getDate();
                return {
                    from: `${month}-01`,
                    to: `${month}-${String(lastDay).padStart(2, "0")}`,
                };
            }
            return null;
        },

        async getSummary() {
            const params = this.getRangeParams();
            if (!params) return;
            const res = await GetCycleSummary(params);
            if (res.metadata.ecode === "Ledger.0000") {
                this.incomeTotal = res.payload.income;
                this.expendTotal = res.payload.expense;
            }
        },

        async handleMonthChange() {
            const params = this.getRangeParams();
            if (!params) return;
            this.chartParams = params;
            await this.getSummary();
            this.refreshChartComponent();
        },

        handleModeChange() {
            if (this.mode === "custom") {
                this.custom = [];
            } else {
                this.month = this.getCurrentModeDefault();
            }
            this.handleMonthChange();
        },

        getCurrentModeDefault() {
            const now = new Date();
            const year = now.getFullYear();
            const month = String(now.getMonth() + 1).padStart(2, "0");
            if (this.mode === "year") return String(year);
            if (this.mode === "month") return `${year}-${month}`;
            return "";
        },

        async onRefresh() {
            this.loading = true;
            try {
                await this.handleMonthChange();
            } finally {
                this.loading = false;
            }
        },

        refreshChartComponent() {
            this.$nextTick(() => {
                const chartRef = this.$refs.dayChartRef;
                if (chartRef && typeof chartRef.refreshChart === "function") {
                    chartRef.refreshChart();
                }
            });
        },
    },
    created() {
        this.$globalBus.emit("updateActivePath", "/overview");
        this.month = this.getCurrentModeDefault();
        this.handleMonthChange();
    },
};
</script>

<style scoped lang="less">
.chart-wrapper {
    display: flex;
    flex-wrap: wrap;
    gap: 15px;
}

.chart-item {
    flex: 1;
    min-width: calc(50%);
    max-width: calc(100%);
    height: 400px;
    border-radius: 8px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
    border: 1px solid #f5f7fa;
    background-color: #fff;
    overflow: hidden;
    padding: 20px;
    box-sizing: border-box;
}
</style>
