import { createRouter, createWebHashHistory } from "vue-router";
const NotFound = () => import("../views/404/index.vue");
const Overview = () => import("../views/overview/index.vue");
const Home = () => import("../views/home/index.vue");
const Layout = () => import("../views/layout/index.vue");
const Transactions = () => import("../views/transactions/index.vue");
const Category = () => import("../views/category/index.vue");

const routes = [
    {
        path: "/",
        component: Layout,
        children: [
            { meta: { title: "首页" }, path: "/", name: "home", component: Home },
            { meta: { title: "数据总览" }, path: "/overview", name: "overview", component: Overview },
            { meta: { title: "记账管理" }, path: "/transactions", name: "transactions", component: Transactions },
            { meta: { title: "分类管理" }, path: "/category", name: "category", component: Category },
        ],
    },
    { meta: { title: "404 页面未找到" }, path: "/:pathMatch(.*)*", component: NotFound },
];

const router = createRouter({
    history: createWebHashHistory("/ledger/"),
    routes,
});

router.beforeEach((to, from, next) => {
    if (to.meta && to.meta.title) {
        document.title = to.meta.title;
    }
    next();
});

export default router;
