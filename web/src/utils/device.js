// 检测是否为移动端
export const isMobile = () => {
    const userAgent = navigator.userAgent.toLowerCase();
    const mobileReg = /iphone|ipad|ipod|android|windows phone|mobile/;
    return mobileReg.test(userAgent) || window.innerWidth < 768;
};
