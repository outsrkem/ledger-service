<template>
    <el-dialog title="修改记账" v-model="dialogVisible" :close-on-click-modal="false" draggable @close="handleDialogClose">
        <div v-loading="loading">
            <el-form :model="formdata" label-width="auto" ref="form" :rules="formRules">
                <el-row :gutter="20">
                    <el-col :span="8">
                        <el-form-item label="金额" prop="amount">
                            <el-input v-model="formdata.amount" />
                        </el-form-item>
                    </el-col>
                    <el-col :span="5">
                        <el-form-item label="">
                            <el-segmented
                                v-model="formdata.direction"
                                :options="[
                                    { value: 2, label: '支出' },
                                    { value: 1, label: '收入' },
                                ]" />
                        </el-form-item>
                    </el-col>
                    <el-col :span="11">
                        <el-form-item label="时间" prop="occ_time">
                            <el-date-picker v-model="formdata.occ_time" type="datetime" value-format="YYYY-MM-DDTHH:mm:ssZZ" />
                        </el-form-item>
                    </el-col>
                </el-row>

                <el-form-item label="分类" prop="category">
                    <TreeCategorySelect ref="treeCat" v-model="formdata.category" :direction="formdata.direction" />
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
import { Plus, Remove } from "@element-plus/icons-vue";
import dayjs from "dayjs";
import TreeCategorySelect from "./TreeCategorySelect.vue";
import { GetTranDetails, EditTransactions } from "../../api/basic.js";
import { msgcon } from "../../utils/message.js";
import { withDelay } from "../../utils/common.js";

const DEFAULT_TRAN_FORM = {
    direction: 2,
    occ_time: "",
    amount: "",
    category: null,
    remark: "",
};

export default {
    name: "UpdateTransactions",
    components: { TreeCategorySelect },
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
            tranId: null,
            targetCidCache: null, // 缓存需要回显的分类id
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
        totalAmount() {
            return this.detailRows.reduce((sum, row) => {
                return sum + (Number(row.total) || 0);
            }, 0);
        },
    },
    methods: {
        async onOpenDialog(val) {
            this.formdata = { ...DEFAULT_TRAN_FORM };
            this.detailRows = [];
            this.tranId = val.id;
            this.targetCidCache = null;
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
        handleDialogClose() {
            this.loading = false;
            this.$nextTick(() => {
                if (this.$refs["form"]) {
                    this.$refs["form"].clearValidate();
                }
            });
        },
        async GetDetails() {
            const res = await withDelay(() => GetTranDetails(this.tranId));
            const payload = res.payload || {};

            const rawAmount = Number(payload.amount);
            this.formdata.direction = rawAmount >= 0 ? 1 : 2;
            this.formdata.amount = Math.abs(rawAmount).toString();

            this.formdata.occ_time = dayjs(payload.occ_time).format("YYYY-MM-DDTHH:mm:ssZZ");
            this.formdata.remark = payload.remark || "";
            this.detailRows = payload.detail ? [...payload.detail] : [];

            // 缓存cid，**不直接赋值**，等组件加载完毕再赋值
            this.targetCidCache = payload.cid;
            this.formdata.category = null;

            // 等待组件挂载+分类加载完成，执行回显
            this.$nextTick(async () => {
                if (!this.$refs.treeCat) return;
                // 等待组件内部分类请求完成
                await this.$refs.treeCat.loadCategory();
                // 赋值触发匹配
                this.formdata.category = this.targetCidCache;
            });
        },
        onSwitchDirection() {
            this.formdata.category = null;
            this.targetCidCache = null;
        },
        async validateForm() {
            const mainFormValid = await new Promise((resolve) => {
                this.$refs["form"].validate((valid) => {
                    resolve(valid);
                });
            });
            if (mainFormValid && !this.formdata.category) {
                this.$message.error("请选择类别");
                return false;
            }
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

            if (!this.formdata.category || !this.formdata.category) {
                this.$message.error(msgcon("请选择有效的分类"));
                return;
            }
            const cid = this.formdata.category;

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

            const data = {
                cid: cid,
                occ_time: this.formdata.occ_time,
                amount: amount,
                remark: this.formdata.remark,
                detail: formattedDetails,
            };

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
};
</script>

<style scoped lang="less">
.flex {
    display: flex;
    align-items: center;
}
.text-center {
    text-align: center;
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
</style>
