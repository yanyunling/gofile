<template>
  <el-dialog :model-value="props.visible" title="重置密码" append-to-body destroy-on-close width="400" :before-close="beforeClose" center>
    <el-form :model="formData" label-width="auto" :rules="rules" ref="formRef">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="formData.username" readonly></el-input>
      </el-form-item>
      <el-form-item label="新密码" prop="password">
        <el-input v-model.trim="formData.password" placeholder="请填写新密码"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="beforeClose">取消</el-button>
      <el-button type="primary" @click="saveClick" :disabled="formLoading">保存</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, Ref, watch, PropType } from "vue";
import { ElMessage } from "element-plus";
import type { FormInstance } from "element-plus";
import UserApi from "@/api/user";

const emit = defineEmits<{ "update:visible": [visible: boolean] }>();

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  formData: {
    type: Object as PropType<User>,
  },
});

const formRef = ref<FormInstance>();
const formLoading = ref(false);
const formData: Ref<User> = ref({
  id: "",
  username: "",
  password: "",
  enable: 1,
  publicAuth: 0,
  protectedAuth: 0,
  privateAuth: 0,
});
const rules = {
  password: [
    { required: true, message: "请填写新密码", trigger: "blur" },
    { min: 6, message: "密码至少6个字符", trigger: "blur" },
  ],
};

watch(
  () => props.visible,
  (val) => {
    if (val) {
      setFormData(props.formData);
    }
  },
);

/**
 * 保存
 */
const saveClick = () => {
  formRef.value.validate((valid) => {
    if (valid) {
      formLoading.value = true;
      UserApi.resetPassword(formData.value)
        .then(() => {
          ElMessage.success("保存成功");
          beforeClose();
        })
        .finally(() => {
          formLoading.value = false;
        });
    }
  });
};

/**
 * 设置表单数据
 * @param user
 */
const setFormData = (user: User) => {
  formData.value.id = user.id;
  formData.value.username = user.username;
  formData.value.enable = user.enable;
  formData.value.publicAuth = user.publicAuth;
  formData.value.protectedAuth = user.protectedAuth;
  formData.value.privateAuth = user.privateAuth;
};

/**
 * 弹窗关闭
 */
const beforeClose = () => {
  emit("update:visible", false);
  formData.value.id = "";
  formData.value.username = "";
  formData.value.password = "";
  formData.value.enable = 1;
  formData.value.publicAuth = 0;
  formData.value.protectedAuth = 0;
  formData.value.privateAuth = 0;
};
</script>
