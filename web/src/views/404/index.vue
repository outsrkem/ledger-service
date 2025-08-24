<template>
    <div class="not-found-container">
        <div class="error-code">404</div>
        <h1>页面未找到</h1>
        <p>抱歉，您访问的页面不存在或已被移除</p>
        <div class="button-group">
            <router-link to="/" class="btn home-btn">返回首页</router-link>
            <button @click="goBack" class="btn back-btn">回到上一页</button>
        </div>
    </div>
</template>

<script>
export default {
    name: "SimpleNotFound",
    methods: {
        goBack() {
            window.history.back();
        },
    },
};
</script>

<style scoped>
@import "./style.css";
</style>
