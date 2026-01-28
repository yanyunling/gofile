<template>
  <div class="page-home">
    <el-button @click="publicClick">公开文件</el-button>
    <el-button @click="protectedClick">保护文件</el-button>
    <el-button @click="privateClick">私有文件</el-button>
    <el-button @click="publicDownloadClick">下载公开文件</el-button>
    <el-button @click="protectedDownloadClick">下载保护文件</el-button>
    <el-button @click="privateDownloadClick">下载私有文件</el-button>
    <div>{{ publicList }}</div>
    <div>{{ protectedList }}</div>
    <div>{{ privateList }}</div>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref, onMounted } from "vue";
import FileApi from "@/api/file";
import { saveAs } from "file-saver";

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

const publicDownloadClick = () => {
  FileApi.downloadPublic("王五", "aaa.docx").then((res) => {
    console.log(res);
    saveAs(res, "aaa.docx");
  });
};

const protectedDownloadClick = () => {
  FileApi.downloadProtected("", "沛县电子政务外网使用情况.xlsx").then((res) => {
    console.log(res);
    saveAs(res, "沛县电子政务外网使用情况.xlsx");
  });
};

const privateDownloadClick = () => {
  FileApi.downloadPrivate("测试", "ccc.docx").then((res) => {
    console.log(res);
    saveAs(res, "ccc.docx");
  });
};
</script>

<style lang="scss" scoped>
.page-home {
}
</style>
