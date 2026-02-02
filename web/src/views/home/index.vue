<template>
  <div class="page-home" v-loading="loading">
    <div class="left-view">
      <div class="item-view" :class="currentMenu === 'public' ? 'selected' : ''" @click="changeType('public')">
        <el-icon><FolderChecked /></el-icon>
        <span>公开文件</span>
      </div>
      <div v-if="accessToken" class="item-view" :class="currentMenu === 'protected' ? 'selected' : ''" @click="changeType('protected')">
        <el-icon><Lock /></el-icon>
        <span>保护文件</span>
      </div>
      <div v-if="accessToken" class="item-view" :class="currentMenu === 'private' ? 'selected' : ''" @click="changeType('private')">
        <el-icon><User /></el-icon>
        <span>个人文件</span>
      </div>
    </div>
    <div class="right-view">
      <div class="right-inner">
        <div class="path-view">
          <span class="path-text" :class="{ active: pathList.length === 0 }" @click="pathClick(-1)">根目录</span>
          <span class="split-text">/</span>
          <template v-for="(path, index) in pathList" :key="index">
            <span class="path-text" :class="{ active: pathList.length === index + 1 }" @click="pathClick(index)">{{ path }}</span>
            <span class="split-text">/</span>
          </template>
        </div>
        <el-table class="table-view" ref="tableRef" :data="fileList" height="calc(100% - 49px)" size="small">
          <el-table-column prop="" label="序号" align="center" width="60">
            <template #default="scope">
              {{ scope.$index + 1 }}
            </template>
          </el-table-column>
          <el-table-column prop="name" label="文件名" align="center">
            <template #default="scope">
              <span v-if="scope.row.isDir" class="folder-name" @click="folderClick(scope.row)">{{ scope.row.name }}</span>
              <span v-else>{{ scope.row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="isDir" label="类型" align="center" width="80">
            <template #default="scope">
              <el-tag v-if="scope.row.isDir" type="primary" disable-transitions>目录</el-tag>
              <el-tag v-else type="success" disable-transitions>文件</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="size" label="文件大小" align="center" width="150">
            <template #default="scope">
              {{ scope.row.isDir ? "--" : Upload.formatFileSize(scope.row.size) }}
            </template>
          </el-table-column>
          <el-table-column prop="updateTime" label="更新时间" align="center" width="150">
            <template #default="scope">
              {{ formatTime(scope.row.updateTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="操作" label="" align="center" width="180">
            <template #header>
              <template v-if="updateAuth()">
                <el-button type="primary" link size="small" @click="uploadClick">上传</el-button>
                <el-button type="primary" link size="small" @click="createFolderClick">创建目录</el-button>
              </template>
              <template v-else>操作</template>
            </template>
            <template #default="scope">
              <el-button v-if="!scope.row.isDir" type="primary" link size="small" @click="downloadClick(scope.row)">下载</el-button>
              <el-button v-if="updateAuth()" type="danger" link size="small" @click="deleteClick(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import FileApi from "@/api/file";
import { saveAs } from "file-saver";
import { Upload } from "@/utils/upload";
import { useTokenStore } from "@/store/token";
import { storeToRefs } from "pinia";
import { FolderChecked, Lock, User } from "@element-plus/icons-vue";
import { formatTime } from "@/utils";

const tokenStore = useTokenStore();
const { isAdmin, accessToken, hasUpdate } = storeToRefs(tokenStore);
const currentMenu = ref("public");
const loading = ref(false);
const fileList: Ref<FileInfo[]> = ref([]);
const pathList: Ref<string[]> = ref([]);

onMounted(() => {
  queryFileList();
});

/**
 * 判断是否有更新权限
 */
const updateAuth = () => {
  if (isAdmin.value) {
    return true;
  }
  if (currentMenu.value === "private" && hasUpdate.value) {
    return true;
  }
  return false;
};

/**
 * 文件分类变更
 * @param menu
 */
const changeType = (menu: string) => {
  currentMenu.value = menu;
  pathList.value = [];
  queryFileList();
};

/**
 * 点击路径
 * @param index -1代表根路径
 */
const pathClick = (index: number) => {
  if (index === -1) {
    pathList.value = [];
  } else {
    pathList.value = pathList.value.slice(0, index + 1);
  }
  queryFileList();
};

/**
 * 点击目录
 * @param row
 */
const folderClick = (row: FileInfo) => {
  pathList.value.push(row.name);
  queryFileList();
};

/**
 * 查询文件列表
 */
const queryFileList = () => {
  fileList.value = [];
  const fileFunc =
    currentMenu.value === "public" ? FileApi.listPublic : currentMenu.value === "protected" ? FileApi.listProtected : FileApi.listPrivate;
  loading.value = true;
  fileFunc(pathList.value.join("/"))
    .then((res) => {
      fileList.value = res.data;
    })
    .finally(() => {
      loading.value = false;
    });
};

/**
 * 上传文件
 */
const uploadClick = () => {
  Upload.openFiles().then((fileList) => {
    if (fileList[0].size >= 1000 * 1000 * 1000) {
      ElMessage.error("文件大小不可超过1GB");
      return;
    }
    const fileFunc =
      currentMenu.value === "public" ? FileApi.uploadPublic : currentMenu.value === "protected" ? FileApi.uploadProtected : FileApi.uploadPrivate;
    ElMessage.info("正在上传文件，请稍候...");
    loading.value = true;
    fileFunc(pathList.value.join("/"), fileList[0])
      .then((res) => {
        ElMessage.success(res.message);
        queryFileList();
      })
      .catch(() => {
        loading.value = false;
      });
  });
};

/**
 * 创建目录
 */
const createFolderClick = () => {
  ElMessageBox.prompt("请输入目录名称", "创建目录", {
    confirmButtonText: "创建",
    cancelButtonText: "取消",
  }).then(({ value }) => {
    const fileFunc =
      currentMenu.value === "public" ? FileApi.folderPublic : currentMenu.value === "protected" ? FileApi.folderProtected : FileApi.folderPrivate;
    loading.value = true;
    fileFunc(pathList.value.join("/"), value)
      .then((res) => {
        ElMessage.success(res.message);
        queryFileList();
      })
      .catch(() => {
        loading.value = false;
      });
  });
};

/**
 * 下载
 */
const downloadClick = (row: FileInfo) => {
  const fileFunc =
    currentMenu.value === "public" ? FileApi.downloadPublic : currentMenu.value === "protected" ? FileApi.downloadProtected : FileApi.downloadPrivate;
  ElMessage.info("正在下载文件，请稍候...");
  fileFunc(pathList.value.join("/"), row.name).then((res) => {
    saveAs(res, row.name);
    ElMessage.success("下载成功");
  });
};

/**
 * 删除
 */
const deleteClick = (row: FileInfo) => {
  ElMessageBox.confirm(`是否删除${row.isDir ? "目录" : "文件"}：${row.name}？${row.isDir ? "这将会删除目录内的全部文件。" : ""}`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    const fileFunc =
      currentMenu.value === "public" ? FileApi.deletePublic : currentMenu.value === "protected" ? FileApi.deleteProtected : FileApi.deletePrivate;
    loading.value = true;
    fileFunc(pathList.value.join("/"), row.name)
      .then((res) => {
        ElMessage.success(res.message);
        queryFileList();
      })
      .catch(() => {
        loading.value = false;
      });
  });
};
</script>

<style lang="scss" scoped>
.page-home {
  display: flex;
  overflow: hidden;
  .left-view {
    background: #fff;
    height: 100%;
    width: 180px;
    overflow: auto;
    padding-top: 10px;
    .item-view {
      width: calc(100% - 3px);
      height: 50px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-wrap: nowrap;
      gap: 5px;
      user-select: none;
      cursor: pointer;
      border-left: 3px #fff solid;
      transition: 0.3s;
      color: #1946c0;
      .el-icon {
        margin-left: -10px;
      }
    }
    .item-view.selected {
      background: #e8ecf8;
      border-left-color: #3d5eb9;
    }
    .item-view:hover {
      background: #dde1f0;
      border-left-color: #3d5eb9;
    }
  }
  .right-view {
    width: calc(100% - 180px);
    display: flex;
    padding: 10px;
    .right-inner {
      width: 100%;
      background: #fff;
      border-radius: 10px;
      overflow: hidden;
      .path-view {
        height: 50px;
        display: flex;
        align-items: center;
        font-size: 14px;
        color: #3d5eb9;
        user-select: none;
        gap: 5px;
        padding: 0 15px;
        border-bottom: 1px #eee solid;
        .path-text {
          cursor: pointer;
        }
        .path-text.active {
          color: #aaa;
        }
        .split-text {
          color: #ccc;
        }
      }
      .table-view {
        .folder-name {
          cursor: pointer;
          color: #3d5eb9;
        }
        .folder-name:hover,
        .folder-name:active {
          color: #032c9e;
          text-decoration: underline;
        }
        .el-scrollbar__wrap {
          display: flex;
        }
      }
    }
  }
}
</style>
