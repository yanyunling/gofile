<template>
  <div class="page-layout">
    <div class="top-view">
      <div class="left-view">
        <div class="title-view">文件服务器</div>
        <el-menu class="menu-view" :default-active="menuDefaultActive" mode="horizontal" @select="menuSelect">
          <el-menu-item index="home">主页</el-menu-item>
          <el-menu-item index="user">用户管理</el-menu-item>
        </el-menu>
      </div>
      <div class="right-view">
        <el-dropdown>
          <div class="text-view">
            {{ nickName }}
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-if="!isAdmin" style="user-select: none" @click="updatePasswordDialogVisible = true">修改密码</el-dropdown-item>
              <el-dropdown-item style="user-select: none" @click="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
    <router-view class="content-view"></router-view>
    <update-password-dialog v-model:visible="updatePasswordDialogVisible" />
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from "vue";
import { ElMessageBox } from "element-plus";
import { useTokenStore } from "@/store/token";
import { storeToRefs } from "pinia";
import TokenApi from "@/api/token";
import updatePasswordDialog from "./update-password-dialog.vue";
import router from "@/router";

const tokenStore = useTokenStore();
const { isAdmin, nickName } = storeToRefs(tokenStore);
const updatePasswordDialogVisible = ref(false);
const menuDefaultActive = ref("");

watch(
  () => router.currentRoute.value.name,
  (val) => {
    menuDefaultActive.value = val?.toString();
  },
);

onMounted(() => {
  menuDefaultActive.value = router.currentRoute.value.name?.toString();
});

/**
 * 菜单选择
 * @param key
 */
const menuSelect = (key: string) => {
  router.push({ name: key });
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
  height: 49px;
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
    flex: 1;
    gap: 20px;
    .title-view {
      font-size: 18px;
      font-weight: bold;
      display: flex;
      align-items: center;
    }
    .menu-view {
      height: 100%;
      flex: 1;
      border-bottom: none;
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
    }
    .text-view:hover {
      color: #0094c1;
    }
  }
}
.content-view {
  width: 100%;
  height: calc(100% - 50px);
  overflow: auto;
  background-color: #fcfcfc;
}
</style>
