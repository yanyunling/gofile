import { createRouter, createWebHashHistory } from "vue-router";

const routes = [
  {
    path: "/",
    name: "layout",
    component: () => import("@/views/layout/index.vue"),
    redirect: "/home",
    children: [
      { path: "/home", name: "home", component: () => import("@/views/home/index.vue") },
      { path: "/user", name: "user", component: () => import("@/views/user/index.vue") },
    ],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  next();
});

export default router;
