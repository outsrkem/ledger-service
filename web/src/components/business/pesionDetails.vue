<template>
    <el-dialog v-model="dialogVisible" width="950px" :close-on-click-modal="false" draggable>
        <el-descriptions :title="dialogTitle" border>
            <el-descriptions-item label="姓名">{{ pesionInfo.name }}</el-descriptions-item>
            <el-descriptions-item label="电话">{{ pesionInfo.phone }}</el-descriptions-item>
            <el-descriptions-item label="身份证号">{{ pesionInfo.idCard }}</el-descriptions-item>
            <el-descriptions-item label="性别">{{ pesionInfo.gender }}</el-descriptions-item>
            <el-descriptions-item label="出生日期">{{ pesionInfo.birthDate }}</el-descriptions-item>
            <el-descriptions-item label="年龄">{{ pesionInfo.age }}</el-descriptions-item>
            <el-descriptions-item label="与户主关系">{{ pesionInfo.relationship }}</el-descriptions-item>
            <el-descriptions-item label="类型">{{ pesionInfo.memberType }} </el-descriptions-item>
            <el-descriptions-item label="户籍地址">{{ pesionInfo.registerAddress }}</el-descriptions-item>
        </el-descriptions>
    </el-dialog>
</template>

<script>
export default {
    name: "DetailsIndex",
    props: {
        vdata: {
            type: Object,
            default: () => ({}),
        },
    },
    data() {
        return {
            dialogVisible: false,
            dialogTitle: "查看居民详情",
        };
    },
    computed: {
        pesionInfo() {
            return this.vdata;
        },
    },
    methods: {
        onOpenDialog() {
            this.dialogVisible = true;
        },
    },
};
</script>
