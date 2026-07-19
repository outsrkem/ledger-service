<template>
    <div class="chart-card" ref="chart" style="height: 300px; width: 100%"></div>
</template>

<script>
import { GetChartDay } from "../../../api/basic.js";
import { withDelay } from "../../../utils/common.js";
import * as echarts from "echarts";

export default {
    name: "DayChart",
    props: {
        queryParams: Object,
    },
    data() {
        return {
            myChart: null,
            resizeTimer: null,
            loadingFlag: false,
        };
    },
    mounted() {
        this.myChart = echarts.init(this.$refs.chart);
        this.getBillChartData();

        window.addEventListener("resize", this.handleResize);
    },
    beforeUnmount() {
        clearTimeout(this.resizeTimer);
        window.removeEventListener("resize", this.handleResize);
        if (this.myChart) {
            this.myChart.dispose();
            this.myChart = null;
        }
    },
    watch: {
        queryParams: {
            handler(newVal, oldVal) {
                // 只在更新时执行，排除首次初始化
                if (oldVal !== undefined) {
                    this.getBillChartData();
                }
            },
            deep: true,
        },
    },
    methods: {
        handleResize() {
            clearTimeout(this.resizeTimer);
            this.resizeTimer = setTimeout(() => {
                this.myChart && this.myChart.resize();
            }, 100);
        },

        async getBillChartData() {
            if (!this.queryParams || this.loadingFlag) return;
            this.loadingFlag = true;

            try {
                this.myChart?.showLoading({
                    text: "加载中...",
                    textColor: "#666",
                    maskColor: "rgba(255,255,255,0.7)",
                });

                const res = await withDelay(() => GetChartDay(this.queryParams));
                const list = res.payload || [];

                const xData = [];
                const incomeData = [];
                const expenseData = [];

                list.forEach((item) => {
                    xData.push(item.day.slice(5));
                    incomeData.push(Number(item.income));
                    expenseData.push(Number(item.expense));
                });

                const option = {
                    grid: {
                        left: 80,
                        right: 80,
                        top: 40,
                        bottom: 40,
                    },
                    title: {
                        text: "月度每日收支统计",
                        left: "center",
                        textStyle: { fontSize: 16 },
                    },
                    legend: {
                        data: ["收入", "支出"],
                        right: 0,
                        top: "",
                        itemWidth: 14,
                        itemHeight: 10,
                    },
                    tooltip: {
                        trigger: "axis",
                        formatter(params) {
                            let tip = params[0].axisValue;
                            params.forEach((p) => {
                                tip += `<br/>${p.seriesName}：${p.value.toFixed(2)}元`;
                            });
                            return tip;
                        },
                    },
                    xAxis: {
                        type: "category",
                        data: xData,
                        axisLabel: { interval: 2, rotate: 30 },
                        boundaryGap: false,
                    },
                    yAxis: {
                        type: "value",
                        name: "金额(元)",
                    },
                    series: [
                        {
                            name: "收入",
                            type: "line",
                            smooth: true,
                            data: incomeData,
                            itemStyle: { color: "#36cbcb" },
                            lineStyle: { width: 2 },
                            symbol: "circle",
                            symbolSize: 6,
                        },
                        {
                            name: "支出",
                            type: "line",
                            smooth: true,
                            data: expenseData,
                            itemStyle: { color: "#f56c6c" },
                            lineStyle: { width: 2 },
                            symbol: "circle",
                            symbolSize: 6,
                        },
                    ],
                };

                this.myChart.setOption(option, { notMerge: false });
                this.myChart?.hideLoading();
            } catch (err) {
                this.myChart?.hideLoading();
                console.error("获取每日收支图表失败：", err);
            } finally {
                this.loadingFlag = false;
            }
        },

        // 父组件手动刷新调用
        refreshChart() {
            this.getBillChartData();
        },
    },
};
</script>
