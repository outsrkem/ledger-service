<template>
    <div>
        <el-dialog v-model="dialogVisible" :title="dialogTitle" width="950px" :close-on-click-modal="false" draggable>
            <el-form :model="huzhu" label-width="120px">
                <el-row>
                    <el-col :span="20">
                        <el-form-item label="户主">
                            <el-input v-model="huzhu.select" disabled></el-input>
                        </el-form-item>
                    </el-col>
                </el-row>
            </el-form>
            <add-pesion ref="AddPesion" :vdata="vdata" @callCloseDialog="onCloseDialog" />
            <div style="display: flex; justify-content: flex-end">
                <el-button @click="onCloseDialog">取消</el-button>
                <el-button type="primary" @click="onSubmit">添加</el-button>
            </div>
        </el-dialog>
    </div>
</template>

<script>
export default {
    name: "UtPersonIndex",
    props: {
        vdata: {
            type: Object,
            default: () => ({}),
        },
    },
    data() {
        return {
            dialogTitle: "添加家庭成员",
            dialogVisible: false,
            huzhu: {
                select: "",
            },
        };
    },
    watch: {
        // 监听vdata变化，当父组件更新vdata时同步更新显示
        vdata: {
            handler(newVal) {
                this.setHuzhu();
            },
            deep: true,
        },
    },
    methods: {
        onOpenDialog() {
            this.dialogVisible = true;
        },
        onCloseDialog() {
            this.dialogVisible = false;
        },
        setHuzhu() {
            // 检查vdata是否存在且包含必要的字段
            if (this.vdata && this.vdata.idCard && this.vdata.name && this.vdata.registerAddress) {
                this.huzhu.select = `${this.vdata.idCard}/${this.vdata.name}/${this.vdata.registerAddress}`;
            }
        },
        // 添加按钮事件处理
        onSubmit() {
            // 调用子页函数提交
            this.$refs.AddPesion.onSubmit();
        },
    },
    created() {
        // 组件初始化时调用setHuzhu方法处理vdata
        this.setHuzhu();
    },
};
</script>

<style scoped lang="less"></style>
