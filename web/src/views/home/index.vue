<template>
  <div class="page-home">
    <el-button @click="publicClick">公开文件</el-button>
    <el-button @click="protectedClick">保护文件</el-button>
    <el-button @click="privateClick">私有文件</el-button>
    <div>{{ publicList }}</div>
    <div>{{ protectedList }}</div>
    <div>{{ privateList }}</div>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref, onMounted } from "vue";
import FileApi from "@/api/file";

const publicList: Ref<FileInfo[]> = ref([]);
const protectedList: Ref<FileInfo[]> = ref([]);
const privateList: Ref<FileInfo[]> = ref([]);

const publicClick = () => {
  FileApi.listPublic("").then((res) => {
    publicList.value = res.data;
    console.log(res);
  });
};

const protectedClick = () => {
  FileApi.listProtected("").then((res) => {
    protectedList.value = res.data;
    console.log(res);
  });
};

const privateClick = () => {
  FileApi.listPrivate("").then((res) => {
    privateList.value = res.data;
    console.log(res);
  });
};
</script>

<style lang="scss" scoped>
.page-home {
}
</style>
