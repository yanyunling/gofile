import { createRouter, createWebHistory } from "vue-router";
import Layout from "@/views/layout/index.vue";

const routes = [
  {
    path: "/",
    name: "layout",
    component: Layout,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
