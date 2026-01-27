import { createRouter, createWebHashHistory } from "vue-router";
import Layout from "@/views/layout/index.vue";
import Home from "@/views/home/index.vue";

const routes = [
  {
    path: "/",
    name: "layout",
    component: Layout,
    redirect: "/home",
    children: [{ path: "/home", name: "home", component: Home }],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
