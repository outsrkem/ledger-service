<template>
    <div>
        <el-form :model="form.pesion" label-width="120px" :rules="formRules" ref="pesion-form">
            <el-row>
                <el-col :span="8">
                    <el-form-item label="姓名" prop="name">
                        <el-input v-model="form.pesion.name"></el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="8">
                    <el-form-item label="身份证号码" prop="idCard">
                        <el-input v-model="form.pesion.idCard" @input="setGender"></el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="8">
                    <el-form-item label="性别">
                        <el-radio-group v-model="form.pesion.gender" disabled>
                            <el-radio value="1">男</el-radio>
                            <el-radio value="2">女</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </el-col>
            </el-row>
            <el-row>
                <el-col :span="8">
                    <el-form-item label="联系电话">
                        <el-input v-model="form.pesion.phone"></el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="8">
                    <el-form-item label="紧急联系人">
                        <el-input v-model="form.pesion.emergencyContact"></el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="8">
                    <el-form-item label="紧急联系电话">
                        <el-input v-model="form.pesion.emergencyPhone"></el-input>
                    </el-form-item>
                </el-col>
            </el-row>
            <el-form-item label="户籍地址">
                <el-input v-model="form.pesion.registerProvince" style="width: 110px" placeholder="省"></el-input>
                <el-input v-model="form.pesion.registerCity" style="width: 110px; margin: 0 5px" placeholder="市"></el-input>
                <el-input v-model="form.pesion.registerDistrict" style="width: 110px; margin: 0 5px" placeholder="区/县"></el-input>
                <el-input v-model="form.pesion.registerTown" style="width: 110px; margin: 0 5px" placeholder="镇/乡"></el-input>
                <el-input v-model="form.pesion.registerVillage" style="width: 110px; margin: 0 5px" placeholder="村"></el-input>
                <el-input v-model="form.pesion.registerGroup" style="width: 110px; margin: 0 5px" placeholder="社/小组"></el-input>
            </el-form-item>
            <el-row>
                <el-col :span="8">
                    <el-form-item label="与户主关系" prop="relationship">
                        <el-select v-model="form.pesion.relationship" placeholder="请选择与户主关系" allow-create filterable>
                            <el-option label="本人" value="本人" />
                            <el-option label="配偶" value="配偶" />
                            <el-option label="子女" value="子女" />
                            <el-option label="父母" value="父母" />
                            <el-option label="祖父母/外祖父母" value="祖父母/外祖父母" />
                            <el-option label="孙子女/外孙子女" value="孙子女/外孙子女" />
                            <el-option label="兄弟姐妹" value="兄弟姐妹" />
                        </el-select>
                    </el-form-item>
                </el-col>
                <el-col :span="8">
                    <el-form-item label="健康状况" prop="healthStatus">
                        <el-select v-model="form.pesion.healthStatus" placeholder="请选择健康状况" allow-create filterable>
                            <el-option label="健康" value="健康" />
                            <el-option label="亚健康" value="亚健康" />
                            <el-option label="患慢性病" value="患慢性病" />
                            <el-option label="残疾" value="残疾" />
                            <el-option label="重病" value="重病" />
                            <el-option label="其他" value="其他" />
                        </el-select>
                    </el-form-item>
                </el-col>
                <el-col :span="8">
                    <el-form-item label="户口在册" prop="isRegistered">
                        <el-radio-group v-model="form.pesion.isRegistered">
                            <el-radio :value="true">是</el-radio>
                            <el-radio :value="false">否</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </el-col>
            </el-row>
            <el-row>
                <el-col :span="8">
                    <el-form-item label="学历" prop="educationId">
                        <el-select v-model="form.pesion.educationId" placeholder="请选择学历">
                            <el-option label="未上过学" value="未上过学" />
                            <el-option label="学前教育" value="学前教育" />
                            <el-option label="小学" value="小学" />
                            <el-option label="初中" value="初中" />
                            <el-option label="高中" value="高中" />
                            <el-option label="中等职业教育" value="中等职业教育" />
                            <el-option label="专科（高职）" value="专科（高职）" />
                            <el-option label="本科" value="本科" />
                            <el-option label="硕士研究生" value="硕士研究生" />
                            <el-option label="博士研究生" value="博士研究生" />
                        </el-select>
                    </el-form-item>
                </el-col>
                <el-col :span="8">
                    <el-form-item label="政治面貌" prop="politicalStatusId">
                        <el-select v-model="form.pesion.politicalStatusId" placeholder="请选择政治面貌">
                            <el-option label="群众" value="mass" />
                            <el-option label="共青团员" value="league_member" />
                            <el-option label="中共预备党员" value="probationary_member" />
                            <el-option label="中共党员" value="party_member" />
                            <el-option label="民主党派" value="democratic_party" />
                            <el-option label="无党派人士" value="nonpartisan" />
                            <el-option label="其他" value="other" />
                        </el-select>
                    </el-form-item>
                </el-col>
                <el-col :span="8">
                    <el-form-item label="民族" prop="nationalityId">
                        <el-select v-model="form.pesion.nationalityId" placeholder="请选择民族">
                            <el-option label="汉族" value="汉族" />
                            <el-option label="蒙古族" value="蒙古族" />
                            <el-option label="回族" value="回族" />
                            <el-option label="藏族" value="藏族" />
                            <el-option label="维吾尔族" value="维吾尔族" />
                            <el-option label="苗族" value="苗族" />
                            <el-option label="彝族" value="彝族" />
                            <el-option label="壮族" value="壮族" />
                            <el-option label="布依族" value="布依族" />
                            <el-option label="朝鲜族" value="朝鲜族" />
                            <el-option label="满族" value="满族" />
                            <el-option label="侗族" value="侗族" />
                            <el-option label="瑶族" value="瑶族" />
                            <el-option label="白族" value="白族" />
                            <el-option label="土家族" value="土家族" />
                            <el-option label="哈尼族" value="哈尼族" />
                            <el-option label="哈萨克族" value="哈萨克族" />
                            <el-option label="傣族" value="傣族" />
                            <el-option label="黎族" value="黎族" />
                            <el-option label="傈僳族" value="傈僳族" />
                            <el-option label="佤族" value="佤族" />
                            <el-option label="畲族" value="畲族" />
                            <el-option label="高山族" value="高山族" />
                            <el-option label="拉祜族" value="拉祜族" />
                            <el-option label="水族" value="水族" />
                            <el-option label="东乡族" value="东乡族" />
                            <el-option label="纳西族" value="纳西族" />
                            <el-option label="景颇族" value="景颇族" />
                            <el-option label="柯尔克孜族" value="柯尔克孜族" />
                            <el-option label="土族" value="土族" />
                            <el-option label="达斡尔族" value="达斡尔族" />
                            <el-option label="仫佬族" value="仫佬族" />
                            <el-option label="羌族" value="羌族" />
                            <el-option label="布朗族" value="布朗族" />
                            <el-option label="撒拉族" value="撒拉族" />
                            <el-option label="毛南族" value="毛南族" />
                            <el-option label="仡佬族" value="仡佬族" />
                            <el-option label="锡伯族" value="锡伯族" />
                            <el-option label="阿昌族" value="阿昌族" />
                            <el-option label="普米族" value="普米族" />
                            <el-option label="塔吉克族" value="塔吉克族" />
                            <el-option label="怒族" value="怒族" />
                            <el-option label="乌孜别克族" value="乌孜别克族" />
                            <el-option label="俄罗斯族" value="俄罗斯族" />
                            <el-option label="鄂温克族" value="鄂温克族" />
                            <el-option label="德昂族" value="德昂族" />
                            <el-option label="保安族" value="保安族" />
                            <el-option label="裕固族" value="裕固族" />
                            <el-option label="京族" value="京族" />
                            <el-option label="塔塔尔族" value="塔塔尔族" />
                            <el-option label="独龙族" value="独龙族" />
                            <el-option label="鄂伦春族" value="鄂伦春族" />
                            <el-option label="赫哲族" value="赫哲族" />
                            <el-option label="门巴族" value="门巴族" />
                            <el-option label="珞巴族" value="珞巴族" />
                            <el-option label="基诺族" value="基诺族" />
                        </el-select>
                    </el-form-item>
                </el-col>
            </el-row>
            <el-row>
                <el-col :span="10">
                    <el-form-item label="职业">
                        <el-input v-model="form.pesion.occupation"></el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="14">
                    <el-form-item label="工作单位">
                        <el-input v-model="form.pesion.workUnit"></el-input>
                    </el-form-item>
                </el-col>
            </el-row>

            <el-row>
                <el-col :span="6">
                    <el-form-item label="出入状态" prop="inOutStatus">
                        <el-select v-model="form.pesion.inOutStatus">
                            <el-option label="外出" value="外出" />
                            <el-option label="在家" value="在家" />
                        </el-select>
                    </el-form-item>
                </el-col>
                <el-col :span="7">
                    <el-form-item label="外出时间">
                        <el-date-picker v-model="form.pesion.outStartDate" type="date" />
                    </el-form-item>
                </el-col>
                <el-col :span="11">
                    <el-form-item label="外出地址">
                        <el-input v-model="form.pesion.outAddress"></el-input>
                    </el-form-item>
                </el-col>
            </el-row>

            <el-form-item label="备注">
                <el-input v-model="form.pesion.remarks" type="textarea" maxlength="128" show-word-limit />
            </el-form-item>
        </el-form>
    </div>
</template>

<script>
import { msgcon } from "../../utils/message.js";
import { UpdateFamilyMember, GetFamilyMemberById } from "../../api/community.js";
import { getGenderCodeByIdCard } from "../../utils/person.js";
export default {
    name: "UpdatePersonIndex",
    props: {
        vdata: {
            type: Object,
            default: () => ({}),
        },
    },
    computed: {},
    watch: {
        vdata: {
            handler(newVal) {
                this.onLoadvdata();
            },
            deep: true,
        },
    },
    data() {
        return {
            huzhu: {
                data: [],
                select: null,
            },
            radio: null,
            form: {
                household: {
                    fullAddress: "",
                },
                pesion: {
                    name: "",
                    idCard: "",
                    gender: "",
                    registerAddress: "",
                    memberType: "",
                    remarks: "",
                    phone: "",
                },
            },
            // TODO 输入检查,需要继续完善
            formRules: {
                name: [{ required: true, type: "string", message: "请输入姓名", trigger: ["blur", "change"] }],
                idCard: [
                    { required: true, message: "请输入身份证号码", trigger: "blur" },
                    {
                        pattern: /(^[0-9]{18}$)|(^[0-9]{17}([0-9]|X)$)/,
                        message: "请输入有效的18位身份证号码",
                        trigger: "change",
                    },
                ],
                nationalityId: [{ required: true, message: "请选择名族" }],
                politicalStatusId: [{ required: true, message: "请选择政治面貌" }],
                educationId: [{ required: true, message: "请选择学历" }],
                isRegistered: [{ required: true, message: "户籍在册必选" }],
                healthStatus: [{ required: true, message: "健康状况必选" }],
                relationship: [{ required: true, message: "与户主关系必选" }],
                inOutStatus: [{ required: true, message: "出入状态必选" }],
            },
        };
    },
    methods: {
        setGender() {
            this.form.pesion.gender = getGenderCodeByIdCard(this.form.pesion.idCard);
        },
        // 拼接地址
        fullAddress() {
            const fullAddress = [
                this.form.pesion.registerProvince,
                this.form.pesion.registerCity,
                this.form.pesion.registerDistrict,
                this.form.pesion.registerTown,
                this.form.pesion.registerVillage,
                this.form.pesion.registerGroup,
            ]
                .filter(Boolean)
                .join(",");
            return fullAddress;
        },
        // 构造数据
        createData() {
            let data = {
                name: this.form.pesion.name,
                idCard: this.form.pesion.idCard,
                registerAddress: this.fullAddress(),
                memberType: 2, //2：非户主，1：户主
                relationship: this.form.pesion.relationship, // 与户主关系
                nationalityId: this.form.pesion.nationalityId, // # 民族ID
                politicalStatusId: this.form.pesion.politicalStatusId, // # 政治面貌ID 政治面貌
                educationId: this.form.pesion.educationId, // 学历ID
                occupation: this.form.pesion.occupation, // 职业
                workUnit: this.form.pesion.workUnit, // 工作单位
                healthStatus: this.form.pesion.healthStatus, // 健康状况
                isRegistered: this.form.pesion.isRegistered, // 是否户口在册
                inOutStatus: this.form.pesion.inOutStatus,
                outAddress: this.form.pesion.outAddress,
                outStartDate: this.form.pesion.outStartDate,
                phone: this.form.pesion.phone,
                emergencyContact: this.form.pesion.emergencyContact, // 紧急联系人
                emergencyPhone: this.form.pesion.emergencyPhone, //  紧急联系电话
                remarks: this.form.pesion.remarks,
            };
            data.kid = this.vdata.kid;
            console.debug(data);
            return data;
        },
        async validateForm() {
            return new Promise((resolve) => {
                this.$refs["pesion-form"].validate((valid) => {
                    resolve(valid);
                });
            });
        },
        async onSubmit() {
            const pesionValid = await this.validateForm();
            if (!pesionValid) return;

            // 构造数据
            const data = this.createData();

            // 发送请求
            UpdateFamilyMember(data).then(() => {
                this.$message.success(msgcon("修改成功"));
                this.$emit("callCloseDialog");
            });
        },
        // 载入更新数据
        async onLoadvdata() {
            this.form.pesion = this.vdata;
        },
        // 提交更新
        async onSubmitUpdate() {},
    },
    created() {},
};
</script>

<style scoped lang="less"></style>
