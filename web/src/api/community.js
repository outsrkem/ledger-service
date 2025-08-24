import ajax from "./ajax.js";

export const GetPerson = (params) => ajax("/v1/community/familyMember/getPageUsers", "GET", params, null);
export const Gethousehold = (params) => ajax("/v1/community/familyMember/getFamilyMemberByMemberType", "GET", params, null);
export const AddPerson = (data) => ajax("/v1/community/familyMember/addFamilyMember", "POST", null, data);

// 根据ID删除家庭成员
export const DelPerson = (params) => ajax("/v1/community/familyMember/deleteFamilyMemberById", "DELETE", params, null);

// 根据用户身份证号查询家庭成员
export const FindPersonByIdCard = (params) => ajax("/v1/community/familyMember/getFamilyMemberByIdCard", "GET", params, null);

// 根据户ID查询所有家庭成员
export const GetfamilyMember = (params) => ajax("/v1/community/familyMember/familyMemberByID", "GET", params, null);

// 根据ID修改家庭成员
export const UpdateFamilyMember = (data) => ajax("/v1/community/familyMember/updateFamilyMemberById", "PUT", null, data);

// 根据ID查询家庭成员
export const GetFamilyMemberById = (params) => ajax("/v1/community/familyMember/getFamilyMemberById", "GET", params, null);
