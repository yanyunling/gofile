<template>
  <el-dialog
    :model-value="props.visible"
    :title="formData.id ? '修改用户信息' : '新增用户信息'"
    append-to-body
    destroy-on-close
    width="500"
    :before-close="beforeClose"
    center
  >
    <el-form :model="formData" label-width="auto" :rules="rules" ref="formRef">
      <el-form-item label="用户名" prop="username">
        <el-input v-model.trim="formData.username" placeholder="请填写用户名" :readonly="formData.id !== ''"></el-input>
      </el-form-item>
      <el-form-item label="密码" prop="password" v-if="!formData.id">
        <el-input v-model.trim="formData.password" placeholder="请填写密码"></el-input>
      </el-form-item>
      <el-form-item label="状态" prop="enable">
        <el-radio-group v-model="formData.enable">
          <el-radio :value="1" border>启用</el-radio>
          <el-radio :value="0" border>停用</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="公开文件权限" prop="publicAuth">
        <el-radio-group v-model="formData.publicAuth">
          <el-radio :value="1" border>有</el-radio>
          <el-radio :value="0" border>无</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="保护文件权限" prop="protectedAuth">
        <el-radio-group v-model="formData.protectedAuth">
          <el-radio :value="1" border>有</el-radio>
          <el-radio :value="0" border>无</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="私有文件权限" prop="privateAuth">
        <el-radio-group v-model="formData.privateAuth">
          <el-radio :value="1" border>有</el-radio>
          <el-radio :value="0" border>无</el-radio>
        </el-radio-group>
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

const emit = defineEmits<{ save: []; "update:visible": [visible: boolean] }>();

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
  username: [
    { required: true, message: "请填写用户名", trigger: "blur" },
    { pattern: /^[a-zA-Z0-9_\-\u4e00-\u9fa5]{1,20}$/, message: "只允许20个字符以内的数字、英文、汉字、中划线、下划线", trigger: "blur" },
  ],
  password: [
    { required: true, message: "请填写密码", trigger: "blur" },
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
      UserApi.save(formData.value)
        .then(() => {
          ElMessage.success("保存成功");
          emit("save");
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
  formData.value.username = "";
  formData.value.password = "";
  formData.value.enable = 1;
  formData.value.publicAuth = 0;
  formData.value.protectedAuth = 0;
  formData.value.privateAuth = 0;
};
</script>
