export function getGenderCodeByIdCard(idCard) {
    // 校验身份证格式（18位，最后可能是X）
    idCard = idCard.trim().toUpperCase();
    if (!/^[0-9]{17}[0-9X]$/.test(idCard)) return undefined; // 修正正则表达式

    // 提取第17位数字（索引16）
    const genderDigit = parseInt(idCard.charAt(16), 10);

    // 奇数→男（1），偶数→女（2）
    return genderDigit % 2 === 1 ? "1" : "2";
}
