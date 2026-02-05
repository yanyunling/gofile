<template>
  <el-dialog v-model="loginVisible" append-to-body destroy-on-close width="500" :before-close="beforeClose" center>
    <div class="page-login">
      <div class="title-view">登录</div>
      <div class="content-view">
        <form>
          <el-input
            class="input-view"
            v-model.trim="signInData.username"
            size="large"
            clearable
            placeholder="请输入用户名"
            @keyup.enter.native="loginClick"
          ></el-input>
          <el-input
            class="input-view"
            v-model.trim="signInData.password"
            size="large"
            type="password"
            clearable
            placeholder="请输入密码"
            @keyup.enter.native="loginClick"
          ></el-input>
        </form>
        <el-button class="input-view" size="large" type="primary" @click="loginClick" :disabled="loading">登录</el-button>
      </div>
      <el-dialog
        class="captcha-dialog"
        :model-value="captchaVisible"
        append-to-body
        destroy-on-close
        :close-on-click-modal="false"
        :close-on-press-escape="false"
        :show-close="false"
        width="360"
        center
      >
        <gocaptcha-slide
          v-loading="captchaLoading"
          :config="{ iconSize: 0 }"
          :data="captchaData"
          :events="{
            confirm: captchaConfirm,
          }"
        />
      </el-dialog>
    </div>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, onMounted, Ref } from "vue";
import { useTokenStore } from "@/store/token";
import { useCommonStore } from "@/store/common";
import TokenApi from "@/api/token";
import { ElMessage } from "element-plus";
import { SlidePoint } from "go-captcha-vue/dist/components/slide/meta/data";
import { storeToRefs } from "pinia";

const tokenStore = useTokenStore();
const commonStore = useCommonStore();
const { loginVisible } = storeToRefs(commonStore);
const loading = ref(false);
const captchaLoading = ref(false);
const captchaVisible = ref(false);

// 登录数据
const signInData: Ref<SignIn> = ref({
  username: "",
  password: "",
  captchaId: "",
  captchaX: 0,
  captchaY: 0,
});
// 验证码数据
const captchaData = ref({
  thumbY: 0,
  thumbWidth: 0,
  thumbHeight: 0,
  image: "",
  thumb: "",
});

onMounted(() => {
  let usernameCache = tokenStore.username;
  if (usernameCache) {
    signInData.value.username = usernameCache;
  }
});

/**
 * 点击登录按钮
 */
const loginClick = () => {
  if (!signInData.value.username) {
    ElMessage.warning("请输入用户名");
    return;
  }
  if (!signInData.value.password) {
    ElMessage.warning("请输入密码");
    return;
  }
  // 打开验证码窗口
  captchaGet();
};

/**
 * 获取验证码
 */
const captchaGet = () => {
  loading.value = true;
  captchaLoading.value = true;
  TokenApi.captcha(signInData.value.captchaId)
    .then((res) => {
      signInData.value.captchaId = res.data.captchaId;
      captchaData.value.thumbY = res.data.y;
      captchaData.value.thumbWidth = res.data.width;
      captchaData.value.thumbHeight = res.data.width;
      captchaData.value.image = res.data.image;
      captchaData.value.thumb = res.data.tile;
      captchaVisible.value = true;
    })
    .catch(() => {
      loading.value = false;
      captchaVisible.value = false;
    })
    .finally(() => {
      captchaLoading.value = false;
    });
};

/**
 * 验证码滑动完成回调
 * @param point
 * @param reset
 */
const captchaConfirm = (point: SlidePoint, reset: () => void) => {
  signInData.value.captchaX = point.x;
  signInData.value.captchaY = point.y;
  captchaLoading.value = true;
  TokenApi.validateCaptcha(signInData.value)
    .then((res) => {
      // 验证成功
      if (res.data) {
        captchaVisible.value = false;
        signIn();
      } else {
        ElMessage.warning("验证失败");
        captchaGet();
      }
      reset();
    })
    .finally(() => {
      captchaLoading.value = false;
    });
};

/**
 * 调用登录接口
 */
const signIn = () => {
  loading.value = true;
  TokenApi.signIn(signInData.value)
    .then((res) => {
      ElMessage.success("登录成功");
      tokenStore.setToken(res.data);
      beforeClose();
      location.reload();
    })
    .finally(() => {
      loading.value = false;
    });
};

/**
 * 弹窗关闭
 */
const beforeClose = () => {
  loginVisible.value = false;
  signInData.value.username = "";
  signInData.value.password = "";
  signInData.value.captchaId = "";
  signInData.value.captchaX = 0;
  signInData.value.captchaY = 0;
};
</script>

<style lang="scss">
.page-login {
  display: flex;
  flex-direction: column;
  align-items: center;
  user-select: none;
  overflow: auto;
  .title-view {
    font-size: 24px;
    font-weight: bold;
    color: #3f3f3f;
  }
  .content-view {
    width: 300px;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    background: rgba(255, 255, 255, 0.2);
    padding: 30px 20px 50px 20px;
    .input-view {
      width: 100%;
      margin-top: 10px;
    }
  }
}
.captcha-dialog {
  .go-captcha .gc-drag-block {
    background-color: #3d5eb9;
  }
}
</style>
