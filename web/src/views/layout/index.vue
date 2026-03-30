<template>
  <div class="page-layout">
    <div class="top-view">
      <div class="left-view" @click="routerChange('/home')">
        <logo style="width: 15px; height: 15px; margin-top: 3px; fill: currentColor" />
        <div class="title-view">文件服务器</div>
      </div>
      <div class="right-view">
        <div class="page-title-view" title="返回首页" @click="routerChange('/home')">{{ pageTitle }}</div>
        <el-dropdown v-if="accessToken">
          <div class="text-view">
            {{ username }}
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-if="!isAdmin" style="user-select: none" @click="updatePasswordDialogVisible = true">修改密码</el-dropdown-item>
              <el-dropdown-item v-if="isAdmin" style="user-select: none" @click="routerChange('/user')">用户管理</el-dropdown-item>
              <el-dropdown-item style="user-select: none" @click="routerChange('/share')">分享记录</el-dropdown-item>
              <el-dropdown-item style="user-select: none" @click="routerChange('/log')">操作日志</el-dropdown-item>
              <el-dropdown-item style="user-select: none" @click="logoutClick">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <div v-else class="text-view" @click="loginClick">登录</div>
      </div>
    </div>
    <router-view class="content-view"></router-view>
    <login-dialog />
    <update-password-dialog v-model:visible="updatePasswordDialogVisible" />
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import { ElMessageBox } from "element-plus";
import { useTokenStore } from "@/store/token";
import { useCommonStore } from "@/store/common";
import { storeToRefs } from "pinia";
import TokenApi from "@/api/token";
import loginDialog from "./login-dialog.vue";
import updatePasswordDialog from "./update-password-dialog.vue";
import logo from "@/icons/logo.svg";
import router from "@/router";

const tokenStore = useTokenStore();
const { isAdmin, accessToken, username } = storeToRefs(tokenStore);
const updatePasswordDialogVisible = ref(false);

/**
 * 计算页面标题
 */
const pageTitle = computed(() => {
  const route = router.currentRoute.value;
  switch (route.name) {
    case "log":
      return "操作日志";
    case "share":
      return "分享记录";
    case "user":
      return "用户管理";
    default:
      return "";
  }
});

/**
 * 点击登录
 */
const loginClick = () => {
  useCommonStore().loginVisible = true;
};

/**
 * 点击退出登录
 */
const logoutClick = () => {
  ElMessageBox.confirm("是否退出登录？", "提示", {
    confirmButtonText: "退出登录",
    cancelButtonText: "取消",
    type: "info",
  })
    .then(() => {
      TokenApi.signOut();
      tokenStore.removeToken();
      location.reload();
    })
    .catch(() => {});
};

const routerChange = (path: string) => {
  router.push(path);
};
</script>

<style lang="scss" scoped>
.page-layout {
  position: fixed;
  width: 100%;
  height: 100%;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  .top-view {
    height: 39px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    user-select: none;
    border-bottom: 1px #eee solid;
    padding: 0 2%;
    white-space: nowrap;
    gap: 20px;
    background: #3d5eb9;
    .left-view {
      display: flex;
      align-items: center;
      height: 100%;
      gap: 5px;
      color: #fff;
      cursor: pointer;
      .title-view {
        font-size: 15px;
        display: flex;
        align-items: center;
      }
    }
    .right-view {
      height: 100%;
      display: flex;
      align-items: center;
      gap: 15px;
      .page-title-view {
        color: #fff;
        font-size: 14px;
        padding: 0 10px;
        cursor: pointer;
      }
      .text-view {
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all 0.3s;
        color: #fff;
        outline: none;
        font-size: 14px;
      }
      .text-view:hover {
        color: #eee;
      }
    }
  }
  .content-view {
    width: 100%;
    height: calc(100% - 40px);
    background-color: #fcfcfc;
  }
}
</style>
