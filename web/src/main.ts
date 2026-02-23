import { createApp } from "vue";
import "./styles/index.scss";
import ElementPlus from "element-plus";
import locale from "element-plus/es/locale/lang/zh-cn";
import App from "./App.vue";
import router from "./router";
import store from "@/store";
import "virtual:svg-icons-register";
import "go-captcha-vue/dist/style.css";
import GoCaptcha from "go-captcha-vue";

const app = createApp(App);
app.use(store);
app.use(router);
app.use(ElementPlus, { locale });
app.use(GoCaptcha);
app.mount("#app");
