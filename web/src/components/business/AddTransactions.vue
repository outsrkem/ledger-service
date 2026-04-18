<template>
    <el-dialog title="添加记账" v-model="dialogVisible" :close-on-click-modal="false" draggable>
        <div>
            <el-form :model="formdata" label-width="auto" :rules="formRules" ref="form">
                <el-row>
                    <el-col :span="8">
                        <el-form-item label="资金流向">
                            <el-segmented
                                v-model="formdata.direction"
                                :options="[
                                    { value: 2, label: '支出' },
                                    { value: 1, label: '收入' },
                                ]"
                                @change="onSwitchCategory"
                            />
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
                        <!-- 父分类容器，使用flex-wrap实现自动换行 -->
                        <div class="category-wrapper">
                            <div v-for="(row, rowIndex) in categorizedCategories" :key="rowIndex" class="category-row">
                                <!-- 行内的每个分类项 -->
                                <el-segmented
                                    v-model="formdata.category"
                                    :options="row"
                                    @change="(val) => onChangeCategory(val, rowIndex)"
                                    class="category-segment"
                                >
                                    <template #default="scope">
                                        <div class="segment-item">
                                            {{ scope.item.name }}
                                        </div>
                                    </template>
                                </el-segmented>
                                <!-- 子分类，只在当前选中行显示 -->
                                <el-segmented
                                    v-if="formdata.category && activeRowIndex === rowIndex"
                                    v-model="formdata.subcategory"
                                    :options="subcategory"
                                    class="subcategory-segment"
                                >
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

                <!-- 明细部分 -->
                <el-form-item label=" ">
                    <!-- 按钮和表头容器，使用Flex布局 -->
                    <div class="detail-top-container">
                        <el-button type="primary" style="width: 340px" :icon="Plus" @click="addDetailRow">添加明细 </el-button>
                    </div>

                    <!-- 明细行 -->
                    <div v-for="(row, index) in detailRows" :key="row.id" class="detail-row">
                        <el-input
                            v-model="row.name"
                            placeholder="物品名称"
                            style="width: 20%; margin-right: 1%"
                            @input="calculateTotal(index)"
                        ></el-input>

                        <el-input
                            v-model="row.quantity"
                            placeholder="数量"
                            style="width: 15%; margin-right: 1%"
                            @input="calculateTotal(index)"
                        ></el-input>
                        <el-select v-model="row.unit" placeholder="单位" filterable allow-create style="width: 10%; margin-right: 1%">
                            <el-option label="斤" value="斤" />
                            <el-option label="克" value="克" />
                            <el-option label="个" value="个" />
                            <el-option label="份" value="份" />
                        </el-select>
                        <el-input
                            v-model="row.price"
                            placeholder="单价"
                            style="width: 15%; margin-right: 1%"
                            @input="calculateTotal(index)"
                        ></el-input>

                        <el-input v-model="row.total" placeholder="总额" style="width: 20%; margin-right: 1%"></el-input>
                        <el-button type="text" :icon="Remove" @click="removeDetailRow(index)"></el-button>
                    </div>

                    <!-- 总计 -->
                    <div class="detail-total" v-if="detailRows.length > 0">
                        <span>明细总计: {{ totalAmount.toFixed(4) }}</span>
                    </div>
                </el-form-item>

                <el-form-item label="写点备注">
                    <el-input v-model="formdata.remark" />
                </el-form-item>
            </el-form>
            <div style="display: flex; justify-content: flex-end">
                <el-button style="width: 150px" @click="onCloseDialog">取消</el-button>
                <el-button style="width: 150px" type="primary" @click="onSubmit(1)">添加后继续</el-button>
                <el-button style="width: 150px" type="primary" @click="onSubmit(0)">添加后关闭</el-button>
            </div>
        </div>
    </el-dialog>
</template>

<script>
import dayjs from "dayjs";
import { msgcon } from "../../utils/message.js";
import { Getcategory, AddTransactions } from "../../api/basic.js";
import { Plus, Remove } from "@element-plus/icons-vue";
export default {
    name: "AddTransactionsIndex",
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
            formdata: {
                direction: 2,
                occ_time: "",
                amount: "",
                category: "",
                subcategory: "",
                remark: "",
            },
            // 明细行数据
            detailRows: [],
            category: [],
            subcategory: [],
            activeRowIndex: -1, // 当前激活的父分类行索引
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
        // 将分类数据按每行rowmax个进行分组
        categorizedCategories() {
            const rowmax = 12; // 每12个分为一组
            const rows = [];
            // 遍历所有分类
            for (let i = 0; i < this.category.length; i += rowmax) {
                rows.push(this.category.slice(i, i + rowmax));
            }
            return rows;
        },
        // 计算所有明细的总额
        totalAmount() {
            return this.detailRows.reduce((sum, row) => {
                return sum + (Number(row.total) || 0);
            }, 0);
        },
    },
    methods: {
        onOpenDialog() {
            this.formdata.amount = "";
            this.formdata.category = "";
            this.formdata.subcategory = "";
            this.formdata.remark = "";
            this.detailRows = []; // 清空明细
            this.onSetNowTime();
            this.dialogVisible = true;
            this.$nextTick(() => {
                if (this.$refs["form"]) {
                    this.$refs["form"].resetFields();
                }
            });
        },
        onCloseDialog() {
            this.dialogVisible = false;
        },
        async validateForm() {
            // 验证主表单
            const mainFormValid = await new Promise((resolve) => {
                this.$refs["form"].validate((valid) => {
                    resolve(valid);
                });
            });
            // 验证明细行
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
                    if (row.unit === null || row.unit === "") {
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
        async onSubmit(val) {
            const valid = await this.validateForm();
            if (!valid) return;

            let cid = 0;
            if (this.formdata.subcategory === "") {
                cid = this.formdata.category;
            } else {
                cid = this.formdata.subcategory;
            }

            // 将总金额转换为数字
            let amount = Number(this.formdata.amount);
            if (this.formdata.direction === 2) {
                amount = amount * -1;
            }

            // 处理明细数据，确保数量和金额为数字类型
            const formattedDetails = this.detailRows.map((row) => ({
                name: row.name,
                quantity: Number(row.quantity), // 转换为数字
                price: Number(row.price), // 转换为数字
                total: Number(row.total), // 转换为数字
                unit: row.unit,
            }));

            // 准备提交的数据，包含格式化后的明细
            const data = {
                cid: cid,
                occ_time: this.formdata.occ_time,
                amount: amount,
                remark: this.formdata.remark,
                detail: formattedDetails, // 提交转换后的明细数据
            };

            this.$nextTick(() => {
                if (this.$refs["form"]) {
                    this.$refs["form"].resetFields();
                }
            });

            AddTransactions(data)
                .then(() => {
                    this.$message.success(msgcon("添加成功"));
                    this.$globalBus.emit("onRefresh");
                    if (val === 1) {
                        // 保留对话框，清空部分字段
                        this.formdata.amount = "";
                        this.formdata.category = "";
                        this.formdata.subcategory = "";
                        this.formdata.remark = "";
                        this.detailRows = []; // 清空明细
                        this.onSetNowTime();
                    } else {
                        this.dialogVisible = false;
                        this.detailRows = []; // 清空明细
                    }
                })
                .catch((err) => {
                    let msg = err.data.metadata.message;
                    this.$message.error(msgcon("添加失败 " + msg));
                });
        },
        onSetNowTime() {
            this.formdata.occ_time = dayjs().format("YYYY-MM-DDTHH:mm:ssZZ");
        },
        onSwitchCategory(val) {
            this.loadGetCategory(val);
            this.formdata.category = "";
            this.formdata.subcategory = "";
            this.subcategory = this.getChildrenWithValue(this.category, val);
        },
        loadGetCategory: async function (direction = 2) {
            Getcategory({ direction: direction }).then((res) => {
                this.category = res.payload.items;
                this.category.map((item) => {
                    item.value = item.id;
                });
            });
        },
        getChildrenWithValue(data, targetId) {
            // 递归查找函数
            function findItem(items) {
                for (const item of items) {
                    // 找到目标ID项
                    if (item.id === targetId) {
                        // 处理children（确保返回数组）
                        const children = item.children || [];
                        // 为每个子项添加value字段（值为自身id）
                        return children.map((child) => ({
                            ...child,
                            value: child.id, // 添加value字段，值等于id
                        }));
                    }

                    // 若当前项有children，递归查找
                    if (item.children && item.children.length > 0) {
                        const result = findItem(item.children);
                        if (result !== null) {
                            return result;
                        }
                    }
                }
                // 未找到目标ID，返回空数组
                return null;
            }

            const result = findItem(data);
            return result !== null ? result : [];
        },

        // 父分类变化时触发，记录当前行索引
        onChangeCategory(value, rowIndex) {
            this.activeRowIndex = rowIndex;
            this.formdata.subcategory = "";
            this.subcategory = this.getChildrenWithValue(this.category, value);
        },

        // 添加明细行
        addDetailRow() {
            this.detailRows.push({
                name: "",
                quantity: null,
                price: null,
                total: null,
            });
        },

        // 移除明细行
        removeDetailRow(index) {
            this.detailRows.splice(index, 1);
        },

        // 计算单行总额
        calculateTotal(index) {
            const row = this.detailRows[index];
            if (row && row.quantity !== null && row.price !== null) {
                const quantity = Number(row.quantity);
                const price = Number(row.price);
                row.total = (quantity * price).toFixed(4);
            }
        },
    },
    created() {
        this.loadGetCategory();
    },
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
    margin-bottom: 2px; /* 行之间的间距 */
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

/* 明细样式 */
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
}

/* 淡入动画效果 */
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
