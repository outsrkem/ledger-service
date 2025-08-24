import { createApp } from "vue";
import router from "./router";
import EventBusPlugin from "./utils/event-bus.js";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import App from "./App.vue";
import "./styles/index.less";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import zhCn from "element-plus/es/locale/lang/zh-cn";
import config from "./config/config";
import vueCookies from "vue-cookies";
const app = createApp(App);

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component); // 注册所有 ElementPlus 图标
}

app.config.globalProperties.$config = config;
app.use(ElementPlus, { locale: zhCn });
app.use(router);
app.use(EventBusPlugin);
app.use(vueCookies);
app.mount("#app");
