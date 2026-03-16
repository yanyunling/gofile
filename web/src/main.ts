import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "@/store";
import ElementPlus from "element-plus";
import locale from "element-plus/es/locale/lang/zh-cn";
import GoCaptcha from "go-captcha-vue";
import "./styles/index.scss";
import "go-captcha-vue/dist/style.css";

const app = createApp(App);
app.use(store);
app.use(router);
app.use(ElementPlus, { locale });
app.use(GoCaptcha);
app.mount("#app");
