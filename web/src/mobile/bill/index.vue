<template>
    <div>
        <van-list v-model:loading="loading" :finished="finished" finished-text="没有更多了" @load="onLoad" :offset="100">
            <!-- 日期分组 -->
            <div v-for="(group, index) in groupedTransactions" :key="index" class="card-box">
                <div class="date-title">{{ group.dateText }}</div>
                <div class="title-line"></div>

                <div style="padding: 0px 20px">
                    <!-- 点击跳转到详情 -->
                    <div class="item" v-for="item in group.items" :key="item.id" @click="toDetail(item.id)" style="cursor: pointer">
                        <div style="display: flex; justify-content: space-between; padding: 4px 0px">
                            <div style="font-size: 15px; color: #000">
                                {{ item.category_title || "未知分类" }}
                            </div>
                            <div style="font-size: 15px">
                                <span :style="{ color: item.amount > 0 ? 'red' : 'black' }">
                                    {{ item.amount }}
                                </span>
                            </div>
                        </div>
                        <div style="display: flex; justify-content: flex-start; padding: 4px 0px">
                            <div style="font-size: 12px; color: #999; margin-top: 2px; margin-right: 10px">
                                {{ formatTimeOnly(item.occ_time) }}
                            </div>
                            <div style="font-size: 12px; color: #999; margin-top: 2px">{{ item.remark }}</div>
                        </div>
                    </div>
                </div>
            </div>
        </van-list>

        <van-dialog v-model:show="showDeleteDialog" title="确认删除" show-cancel-button @confirm="onDeleteTransactions">
            <div style="margin: 10px; display: flex; justify-content: center">确定要删除该记录吗？</div>
        </van-dialog>

        <add-transactions ref="AddTransactions" />

        <!-- 回到顶部 -->
        <van-back-top bottom="60" :offset="1000" />
    </div>
</template>

<script>
import { GetTransactions, Getcategory, DelTransactions } from "../../api/basic.js";
import { withDelay, convertToLimitOffset } from "../../utils/common.js";
import { getCategoryPath } from "../../utils/category.js";

import { List as VanList, Dialog as VanDialog, BackTop as VanBackTop } from "vant";
import { showSuccessToast } from "vant";

export default {
    name: "AccountListMob",
    components: {
        VanList,
        VanDialog,
        VanBackTop,
    },
    data() {
        return {
            loading: false,
            finished: false,
            refreshing: false,

            page: 1,
            pageSize: 10,
            transactions: [],
            groupedTransactions: [],

            showDeleteDialog: false,
            currentDeleteItem: null,
        };
    },
    methods: {
        // 跳转到账单详情页
        toDetail(id) {
            this.$router.push({
                path: `/bill/${id}`,
            });
        },

        // 格式化 时:分
        formatTimeOnly(time) {
            const date = new Date(time);
            const hour = String(date.getHours()).padStart(2, "0");
            const minute = String(date.getMinutes()).padStart(2, "0");
            return `${hour}:${minute}`;
        },

        // 智能日期：今天 / 昨天 / 年月日
        formatSmartDate(time) {
            const now = new Date();
            const target = new Date(time);

            const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
            const targetDay = new Date(target.getFullYear(), target.getMonth(), target.getDate());
            const diffTime = targetDay - today;
            const diffDays = diffTime / (1000 * 60 * 60 * 24);

            // 星期数组
            const weekDays = ["星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"];
            const week = weekDays[target.getDay()];

            // 格式化年月日
            const year = target.getFullYear();
            const month = String(target.getMonth() + 1).padStart(2, "0");
            const day = String(target.getDate()).padStart(2, "0");
            let dateStr = `${year}年${month}月${day}日 ${week}`;

            // 拼接今天/昨天标识
            if (diffDays === 0) {
                dateStr += "（今天）";
            } else if (diffDays === -1) {
                dateStr += "（昨天）";
            }

            return dateStr;
        },

        // 按日期分组 - 修复日期分组key生成方式
        groupTransactionsByDate(list) {
            const groups = {};
            const safeList = Array.isArray(list) ? list : [];

            safeList.forEach((item) => {
                // 使用更可靠的日期分组方式，避免时区和格式问题
                const date = new Date(item.occ_time);
                const year = date.getFullYear();
                const month = String(date.getMonth() + 1).padStart(2, "0");
                const day = String(date.getDate()).padStart(2, "0");
                const key = `${year}-${month}-${day}`;

                if (!groups[key]) {
                    groups[key] = {
                        dateText: this.formatSmartDate(item.occ_time),
                        items: [],
                    };
                }
                groups[key].items.push(item);
            });

            return Object.values(groups).sort((a, b) => {
                return new Date(b.items[0].occ_time) - new Date(a.items[0].occ_time);
            });
        },

        async onLoad() {
            try {
                await this.loadData();
            } catch {
                this.loading = false;
            }
        },

        async loadData() {
            const params = convertToLimitOffset(this.page, this.pageSize);
            try {
                const res = await withDelay(() => GetTransactions(params));
                const resp = res.payload?.items || [];

                // localStorage为空时兜底空数组，不再是null
                const cacheStr = localStorage.getItem("category");
                let categoryList = cacheStr ? JSON.parse(cacheStr) : [];

                resp.forEach((item) => {
                    item.category_title = getCategoryPath(categoryList, item.category_id) || "未知分类";
                });

                // 使用 Set 去重，基于 id，避免重复数据
                const existingIds = new Set(this.transactions.map((t) => t.id));
                const newItems = resp.filter((item) => !existingIds.has(item.id));

                if (this.page === 1) {
                    this.transactions = newItems;
                } else {
                    this.transactions = [...this.transactions, ...newItems];
                }

                // 重新分组
                this.groupedTransactions = this.groupTransactionsByDate(this.transactions);

                // 判断是否还有更多数据
                if (resp.length < this.pageSize) {
                    this.finished = true;
                } else {
                    this.page++;
                }
            } catch (err) {
                console.error("加载失败", err);
            } finally {
                this.loading = false;
                this.refreshing = false;
            }
        },

        async onRefresh() {
            this.refreshing = true;
            this.page = 1;
            this.finished = false;
            this.transactions = [];
            this.groupedTransactions = [];

            try {
                await this.loadAllCategory();
                await this.loadData();
            } catch {
                this.refreshing = false;
                this.loading = false;
            }
        },

        async loadAllCategory() {
            const [res1, res2] = await Promise.all([Getcategory({ direction: 2 }), Getcategory({ direction: 1 })]);
            const all = [...(res1.payload?.items || []), ...(res2.payload?.items || [])];
            localStorage.setItem("category", JSON.stringify(all));
        },

        onAddTransactions() {
            this.$refs.AddTransactions.onOpenDialog();
        },

        onDeleteConfirm(val) {
            this.currentDeleteItem = val;
            this.showDeleteDialog = true;
        },

        onDeleteTransactions() {
            if (!this.currentDeleteItem) return;
            DelTransactions(this.currentDeleteItem.id)
                .then(() => {
                    showSuccessToast("删除成功");
                    // 删除后重置并刷新列表
                    this.page = 1;
                    this.finished = false;
                    this.transactions = [];
                    this.groupedTransactions = [];
                    this.onLoad();
                })
                .catch((err) => {
                    console.log(err);
                })
                .finally(() => {
                    this.showDeleteDialog = false;
                    this.currentDeleteItem = null;
                });
        },
    },
    async created() {
        this.$globalBus.on("onRefresh", () => {
            this.onRefresh();
        });

        // 页面初始化判断分类缓存，无缓存先加载分类
        const cache = localStorage.getItem("category");
        if (!cache) {
            await this.loadAllCategory();
        }
        await this.onLoad();
    },
};
</script>

<style scoped lang="less">
/* 外层大区域 */
.card-box {
    background: #fff;
    border-radius: 8px;
    margin: 16px;
}

/* 日期标题 */
.date-title {
    font-size: 14px;
    color: #333;
    padding: 20px 20px;
}

/* 标题下方横线 */
.title-line {
    height: 1px;
    background: #eee;
    margin-bottom: 5px;
}

/* 列表项 */
.item {
    padding: 10px 0;
    position: relative;
    font-size: 14px;
    color: #666;
}

/* 两个 item 之间的横线（伪类实现） */
.item + .item::before {
    content: "";
    position: absolute;
    box-sizing: border-box;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: #ebedf0;
    transform: scaleY(0.5);
    pointer-events: none;
}
</style>
