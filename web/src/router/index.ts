import { createRouter, createWebHashHistory } from "vue-router";
import Layout from "@/views/layout/index.vue";

const routes = [
  {
    path: "/",
    name: "layout",
    component: Layout,
    redirect: "/home",
    children: [
      {
        path: "/home",
        name: "home",
        component: () => import("@/views/home/index.vue"),
      },
      {
        path: "/log",
        name: "log",
        component: () => import("@/views/log/index.vue"),
      },
      {
        path: "/share",
        name: "share",
        component: () => import("@/views/share/index.vue"),
      },
      {
        path: "/user",
        name: "user",
        component: () => import("@/views/user/index.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
