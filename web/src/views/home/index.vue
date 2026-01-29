<template>
  <div class="page-home">
    <el-button @click="publicClick">公开文件</el-button>
    <el-button @click="protectedClick">保护文件</el-button>
    <el-button @click="privateClick">私有文件</el-button>
    <el-button @click="publicDownloadClick">下载公开文件</el-button>
    <el-button @click="protectedDownloadClick">下载保护文件</el-button>
    <el-button @click="privateDownloadClick">下载私有文件</el-button>
    <el-button @click="publicUploadClick">上传公开文件</el-button>
    <el-button @click="publicFolderClick">创建公开目录</el-button>
    <el-button @click="protectedFolderClick">创建保护目录</el-button>
    <el-button @click="privateFolderClick">创建私有目录</el-button>
    <el-button @click="protectedUploadClick">上传保护文件</el-button>
    <el-button @click="privateUploadClick">上传私有文件</el-button>
    <el-button @click="publicDeleteClick">删除公开文件</el-button>
    <el-button @click="protectedDeleteClick">删除保护文件</el-button>
    <el-button @click="privateDeleteClick">删除私有文件</el-button>
    <div>{{ publicList }}</div>
    <div>{{ protectedList }}</div>
    <div>{{ privateList }}</div>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import FileApi from "@/api/file";
import { saveAs } from "file-saver";
import { Upload } from "@/utils/upload";

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

const publicFolderClick = () => {
  FileApi.folderPublic("", "123").then((res) => {
    console.log(res);
    ElMessage.success(res.message);
  });
};

const protectedFolderClick = () => {
  FileApi.folderProtected("", "试试").then((res) => {
    console.log(res);
    ElMessage.success(res.message);
  });
};

const privateFolderClick = () => {
  FileApi.folderPrivate("", "ffdsa2说/|?").then((res) => {
    console.log(res);
    ElMessage.success(res.message);
  });
};

const publicUploadClick = () => {
  Upload.openFiles().then((fileList) => {
    if (fileList[0].size >= 1000 * 1000 * 1000) {
      ElMessage.error("文件大小不可超过1GB");
      return;
    }
    FileApi.uploadPublic("", fileList[0]).then((res) => {
      console.log(res);
      ElMessage.success(res.message);
    });
  });
};

const protectedUploadClick = () => {
  Upload.openFiles().then((fileList) => {
    if (fileList[0].size >= 1000 * 1000 * 1000) {
      ElMessage.error("文件大小不可超过1GB");
      return;
    }
    FileApi.uploadProtected("", fileList[0]).then((res) => {
      console.log(res);
      ElMessage.success(res.message);
    });
  });
};

const privateUploadClick = () => {
  Upload.openFiles().then((fileList) => {
    if (fileList[0].size >= 1000 * 1000 * 1000) {
      ElMessage.error("文件大小不可超过1GB");
      return;
    }
    FileApi.uploadPrivate("", fileList[0]).then((res) => {
      console.log(res);
      ElMessage.success(res.message);
    });
  });
};

const publicDeleteClick = () => {
  FileApi.deletePublic("", "王五").then((res) => {
    console.log(res);
    ElMessage.success(res.message);
  });
};

const protectedDeleteClick = () => {
  FileApi.deleteProtected("", "生成图标.png").then((res) => {
    console.log(res);
    ElMessage.success(res.message);
  });
};

const privateDeleteClick = () => {
  FileApi.deletePrivate("", "生成图标.png").then((res) => {
    console.log(res);
    ElMessage.success(res.message);
  });
};
</script>

<style lang="scss" scoped>
.page-home {
}
</style>
