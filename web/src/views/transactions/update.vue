<template>
    <el-dialog title="修改记账" v-model="dialogVisible" :close-on-click-modal="false" draggable @close="handleDialogClose">
        <div v-loading="loading">
            <el-form :model="formdata" label-width="auto" ref="form" :rules="formRules">
                <el-row>
                    <el-col :span="8">
                        <el-form-item label="资金流向">
                            <el-segmented
                                v-model="formdata.direction"
                                :options="[
                                    { value: 2, label: '支出' },
                                    { value: 1, label: '收入' },
                                ]"
                                @change="onSwitchCategory" />
                        </el-form-item>
                    </el-col>
                    <el-col :span="16">
                        <el-form-item label="时间" prop="occ_time">
                            <el-date-picker v-model="formdata.occ_time" type="datetime" value-format="YYYY-MM-DDTHH:mm:ssZZ" />
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-form-item label="金额" prop="amount">
                    <el-input v-model="formdata.amount" />
                </el-form-item>
                <el-form-item label="分类" prop="category">
                    <div class="category-container">
                        <div class="category-wrapper">
                            <div v-for="(row, rowIndex) in categorizedCategories" :key="rowIndex" class="category-row">
                                <el-segmented
                                    v-model="formdata.category"
                                    :options="row"
                                    @change="(val) => onChangeCategory(val, rowIndex)"
                                    class="category-segment">
                                    <template #default="scope">
                                        <div class="segment-item">
                                            {{ scope.item.name }}
                                        </div>
                                    </template>
                                </el-segmented>
                                <el-segmented
                                    v-if="formdata.category && activeRowIndex === rowIndex"
                                    v-model="formdata.subcategory"
                                    :options="subcategory"
                                    class="subcategory-segment">
                                    <template #default="scope">
                                        <div class="segment-item">
                                            {{ scope.item.name }}
                                        </div>
                                    </template>
                                </el-segmented>
                            </div>
                        </div>
                    </div>
                </el-form-item>

                <el-form-item label=" ">
                    <div class="detail-top-container">
                        <el-button type="primary" style="width: 340px" :icon="Plus" @click="addDetailRow">添加明细 </el-button>
                    </div>

                    <div v-for="(row, index) in detailRows" :key="row.id" class="detail-row">
                        <el-input v-model="row.name" placeholder="物品名称" style="width: 25%; margin-right: 1%" @input="calculateTotal(index)"></el-input>
                        <el-input v-model="row.price" placeholder="单价" style="width: 10%; margin-right: 1%" @input="calculateTotal(index)"></el-input>
                        <el-input v-model="row.quantity" placeholder="数量" style="width: 15%; margin-right: 1%" @input="calculateTotal(index)"></el-input>
                        <el-select v-model="row.unit" placeholder="单位" filterable allow-create style="width: 10%; margin-right: 1%">
                            <el-option label="斤" value="斤" />
                            <el-option label="克" value="克" />
                            <el-option label="个" value="个" />
                            <el-option label="份" value="份" />
                        </el-select>
                        <el-input v-model="row.total" placeholder="总额" style="width: 20%; margin-right: 1%"></el-input>
                        <el-button type="text" :icon="Remove" @click="removeDetailRow(index)"></el-button>
                    </div>

                    <div class="detail-total" v-if="detailRows.length > 0">
                        <span>明细总计: {{ totalAmount.toFixed(4) }}</span>
                    </div>
                </el-form-item>

                <el-form-item label="写点备注">
                    <el-input v-model="formdata.remark" />
                </el-form-item>
            </el-form>
            <div style="display: flex; justify-content: flex-end; gap: 12px">
                <el-button style="width: 150px" @click="onCloseDialog">取消</el-button>
                <el-button style="width: 150px" type="primary" @click="onSubmit">确认</el-button>
            </div>
        </div>
    </el-dialog>
</template>

<script>
import dayjs from "dayjs";
import { msgcon } from "../../utils/message.js";
import { Getcategory, GetTranDetails, EditTransactions } from "../../api/basic.js";
import { Plus, Remove } from "@element-plus/icons-vue";

const DEFAULT_TRAN_FORM = {
    direction: 2,
    occ_time: "",
    amount: "",
    category: "",
    subcategory: "",
    remark: "",
};

export default {
    name: "UpdateTransactions",
    props: {
        vdata: {
            type: Object,
            default: () => ({}),
        },
    },
    setup() {
        return {
            Plus,
            Remove,
        };
    },
    data() {
        return {
            dialogVisible: false,
            loading: false,
            formdata: { ...DEFAULT_TRAN_FORM },
            detailRows: [],
            category: [],
            subcategory: [],
            activeRowIndex: -1,
            tranId: null,
            // 和新增页面保持一致校验规则
            formRules: {
                occ_time: [{ required: true, type: "string", message: "请选择时间", trigger: ["blur", "change"] }],
                amount: [
                    { required: true, message: "请输入金额", trigger: ["blur"] },
                    { pattern: /(^[1-9]([0-9]+)?(\.[0-9]{1,4})?$)|(^(0){1}$)|(^[0-9]\.[0-9]{1,3}?$)/, message: "请输入正确额格式,可保留四位小数" },
                ],
                category: [{ required: true, message: "请选择类别", trigger: ["blur", "change"] }],
            },
        };
    },
    computed: {
        categorizedCategories() {
            const rowmax = 12;
            const rows = [];
            for (let i = 0; i < this.category.length; i += rowmax) {
                rows.push(this.category.slice(i, i + rowmax));
            }
            return rows;
        },
        totalAmount() {
            return this.detailRows.reduce((sum, row) => {
                return sum + (Number(row.total) || 0);
            }, 0);
        },
    },
    methods: {
        // 打开弹窗入口
        async onOpenDialog(val) {
            this.formdata = { ...DEFAULT_TRAN_FORM };
            this.detailRows = [];
            this.tranId = val.id;
            this.dialogVisible = true;
            this.loading = true;
            try {
                await this.GetDetails();
            } finally {
                this.loading = false;
            }
        },
        onCloseDialog() {
            this.dialogVisible = false;
        },
        // 弹窗关闭回调，清除校验红框 + 重置loading
        handleDialogClose() {
            this.loading = false;
            this.$nextTick(() => {
                if (this.$refs["form"]) {
                    this.$refs["form"].clearValidate();
                }
            });
        },
        // 获取账单详情 + 分类回显核心逻辑
        async GetDetails() {
            const res = await GetTranDetails(this.tranId);
            const payload = res.payload || {};

            // 金额正负处理
            const rawAmount = Number(payload.amount);
            this.formdata.direction = rawAmount >= 0 ? 1 : 2;
            this.formdata.amount = Math.abs(rawAmount).toString();

            this.formdata.occ_time = dayjs(payload.occ_time).format("YYYY-MM-DDTHH:mm:ssZZ");
            this.formdata.remark = payload.remark || "";
            this.detailRows = payload.detail ? [...payload.detail] : [];

            const targetCid = payload.cid;
            // 先加载对应流向分类
            await this.loadGetCategory(this.formdata.direction);

            // 递归查找cid，区分是父分类还是子分类
            const findCidInfo = (items) => {
                for (let i = 0; i < items.length; i++) {
                    const item = items[i];
                    if (item.id === targetCid) {
                        // 匹配到父分类
                        return {
                            type: "parent",
                            id: item.id,
                            rowIndex: this.getCategoryRowIndex(item.id),
                        };
                    }
                    if (item.children && item.children.length) {
                        const childItem = item.children.find((c) => c.id === targetCid);
                        if (childItem) {
                            const rowIdx = this.getCategoryRowIndex(item.id);
                            return {
                                type: "child",
                                parentId: item.id,
                                childId: childItem.id,
                                rowIndex: rowIdx,
                            };
                        }
                        const res = findCidInfo(item.children);
                        if (res) return res;
                    }
                }
                return null;
            };

            const cidInfo = findCidInfo(this.category);
            if (cidInfo) {
                this.activeRowIndex = cidInfo.rowIndex;
                if (cidInfo.type === "parent") {
                    this.formdata.category = cidInfo.id;
                } else {
                    this.formdata.category = cidInfo.parentId;
                    this.$nextTick(() => {
                        this.subcategory = this.getChildrenWithValue(this.category, cidInfo.parentId);
                        this.formdata.subcategory = cidInfo.childId;
                    });
                }
            }
        },
        // 根据父分类id，找到属于第几行（用于activeRowIndex）
        getCategoryRowIndex(parentId) {
            for (let r = 0; r < this.categorizedCategories.length; r++) {
                const rowItems = this.categorizedCategories[r];
                if (rowItems.some((item) => item.id === parentId)) {
                    return r;
                }
            }
            return -1;
        },

        // 表单+明细校验（完全对齐新增页面）
        async validateForm() {
            const mainFormValid = await new Promise((resolve) => {
                this.$refs["form"].validate((valid) => {
                    resolve(valid);
                });
            });
            if (this.detailRows.length > 0) {
                for (let i = 0; i < this.detailRows.length; i++) {
                    const row = this.detailRows[i];
                    if (!row.name) {
                        this.$message.error(`第${i + 1}行明细：请输入物品名称`);
                        return false;
                    }
                    if (row.quantity === null || isNaN(Number(row.quantity)) || Number(row.quantity) <= 0) {
                        this.$message.error(`第${i + 1}行明细：请输入有效的个数`);
                        return false;
                    }
                    if (!row.unit) {
                        this.$message.error(`第${i + 1}行明细：请输入有效的单位`);
                        return false;
                    }
                    if (row.price === null || isNaN(Number(row.price)) || Number(row.price) <= 0) {
                        this.$message.error(`第${i + 1}行明细：请输入有效的单价`);
                        return false;
                    }
                }
            }
            return mainFormValid;
        },

        async onSubmit() {
            const valid = await this.validateForm();
            if (!valid) return;

            const cid = this.formdata.subcategory ? this.formdata.subcategory : this.formdata.category;
            let amount = Number(this.formdata.amount);
            if (this.formdata.direction === 2) {
                amount = amount * -1;
            }

            const formattedDetails = this.detailRows.map((row) => ({
                name: row.name,
                quantity: Number(row.quantity),
                price: Number(row.price),
                total: Number(row.total),
                unit: row.unit,
            }));

            // body不再携带id
            const data = {
                cid: cid,
                occ_time: this.formdata.occ_time,
                amount: amount,
                remark: this.formdata.remark,
                detail: formattedDetails,
            };

            // 第一个参数传tranId，第二个传表单数据
            EditTransactions(this.tranId, data)
                .then(() => {
                    this.$message.success(msgcon("修改成功"));
                    this.$globalBus.emit("onRefresh");
                    this.dialogVisible = false;
                })
                .catch((err) => {
                    let msg = err.data?.metadata?.message || err.message;
                    this.$message.error(msgcon("修改失败 " + msg));
                });
        },

        onSwitchCategory(val) {
            this.loadGetCategory(val);
            this.formdata.category = "";
            this.formdata.subcategory = "";
            this.subcategory = this.getChildrenWithValue(this.category, val);
            this.activeRowIndex = -1;
        },

        loadGetCategory: async function (direction = 2) {
            const res = await Getcategory({ direction: direction });
            this.category = res.payload.items;
            this.category.map((item) => {
                item.value = item.id;
            });
        },

        getChildrenWithValue(data, targetId) {
            function findItem(items) {
                for (const item of items) {
                    if (item.id === targetId) {
                        const children = item.children || [];
                        return children.map((child) => ({
                            ...child,
                            value: child.id,
                        }));
                    }
                    if (item.children && item.children.length > 0) {
                        const result = findItem(item.children);
                        if (result !== null) return result;
                    }
                }
                return null;
            }
            const result = findItem(data);
            return result !== null ? result : [];
        },

        onChangeCategory(value, rowIndex) {
            this.activeRowIndex = rowIndex;
            this.formdata.subcategory = "";
            this.subcategory = this.getChildrenWithValue(this.category, value);
        },

        addDetailRow() {
            this.detailRows.push({
                id: Date.now() + Math.random(),
                name: "",
                quantity: null,
                price: null,
                total: null,
                unit: "",
            });
        },

        removeDetailRow(index) {
            this.detailRows.splice(index, 1);
        },

        calculateTotal(index) {
            const row = this.detailRows[index];
            if (row && row.quantity !== null && row.price !== null) {
                const quantity = Number(row.quantity);
                const price = Number(row.price);
                row.total = (quantity * price).toFixed(4);
            }
        },
    },
    // =========重要改动：删除created里自动加载分类========
    // created() {
    //     this.loadGetCategory();
    // },
};
</script>

<style scoped lang="less">
.el-segmented {
    word-break: break-word;
    white-space: pre-line;
}
.flex {
    display: flex;
    align-items: center;
}
.text-center {
    text-align: center;
}

.category-container {
    width: 100%;
}

.category-wrapper {
    width: 100%;
}

.category-row {
    margin-bottom: 2px;
}

.category-segment {
    margin-bottom: 4px;
}

.subcategory-segment {
    width: 100%;
    animation: fadeIn 0.2s ease-in-out;
}

.segment-item {
    height: 42px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.detail-top-container {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
}

.detail-row {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
}

.detail-total {
    margin-top: 10px;
    text-align: right;
    font-weight: bold;
    color: #1890ff;
    border-top: 1px solid #e8e8e8;
    padding-top: 8px;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(-5px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
