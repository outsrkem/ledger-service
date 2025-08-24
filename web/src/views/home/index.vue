<template></template>

<script>
export default {
    name: "HomeIndex",
    data() {
        return {};
    },
    methods: {},
    created() {
        this.$router.push({ name: "overview" });
    },
};
</script>

<style scoped lang="less"></style>
