<template>
  <div class="page-layout">
    <div class="top-view">
      <div class="left-view">
        <svg-icon name="logo" customStyle="width: 15px; height: 15px; margin-top: 2px"></svg-icon>
        <div class="title-view">文件服务器</div>
      </div>
      <div class="right-view">
        <el-dropdown v-if="username">
          <div class="text-view">
            {{ username }}
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-if="!isAdmin" style="user-select: none" @click="updatePasswordDialogVisible = true">修改密码</el-dropdown-item>
              <el-dropdown-item v-if="isAdmin" style="user-select: none" @click="userDialogVisible = true">用户管理</el-dropdown-item>
              <el-dropdown-item style="user-select: none" @click="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <div v-else class="text-view" @click="login">登录</div>
      </div>
    </div>
    <router-view class="content-view"></router-view>
    <update-password-dialog v-model:visible="updatePasswordDialogVisible" />
    <user-dialog v-model:visible="userDialogVisible" />
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ElMessageBox } from "element-plus";
import { useTokenStore } from "@/store/token";
import { useCommonStore } from "@/store/common";
import { storeToRefs } from "pinia";
import TokenApi from "@/api/token";
import updatePasswordDialog from "@/views/user/update-password-dialog.vue";
import userDialog from "@/views/user/index.vue";
import svgIcon from "@/components/svg-icon";

const tokenStore = useTokenStore();
const { isAdmin, username } = storeToRefs(tokenStore);
const updatePasswordDialogVisible = ref(false);
const userDialogVisible = ref(false);

/**
 * 点击登录
 */
const login = () => {
  useCommonStore().loginVisible = true;
};

/**
 * 点击退出登录
 */
const logout = () => {
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
}
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
  .left-view {
    display: flex;
    align-items: center;
    height: 100%;
    gap: 5px;
    color: #3d5eb9;
    .title-view {
      font-size: 15px;
      font-weight: bold;
      display: flex;
      align-items: center;
    }
  }
  .right-view {
    height: 100%;
    display: flex;
    align-items: center;
    .text-view {
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.3s;
      color: #303133;
      outline: none;
      font-size: 14px;
    }
    .text-view:hover {
      color: #3d5eb9;
    }
  }
}
.content-view {
  width: 100%;
  height: calc(100% - 40px);
  overflow: auto;
  background-color: #fcfcfc;
}
</style>
