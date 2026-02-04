<template>
  <div class="page-home" v-loading="loading">
    <div class="left-view">
      <div class="item-view" :class="parentDir === 'public' ? 'selected' : ''" @click="changeType('public')">
        <el-icon><FolderChecked /></el-icon>
        <span>公开文件</span>
      </div>
      <div v-if="accessToken" class="item-view" :class="parentDir === 'protected' ? 'selected' : ''" @click="changeType('protected')">
        <el-icon><Lock /></el-icon>
        <span>保护文件</span>
      </div>
      <div v-if="accessToken" class="item-view" :class="parentDir === 'private' ? 'selected' : ''" @click="changeType('private')">
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
          <el-table-column prop="type" label="文件类型" align="center" width="150">
            <template #default="scope">
              {{ scope.row.isDir ? "--" : Upload.getFileExtension(scope.row.name) }}
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
              <el-button v-if="!scope.row.isDir && accessToken" type="success" link size="small" @click="shareClick(scope.row)">分享</el-button>
              <el-button v-if="updateAuth()" type="danger" link size="small" @click="deleteClick(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
    <el-drawer
      v-model="drawerData.visible"
      :title="drawerData.isUpload ? '文件上传' : '文件下载'"
      direction="ltr"
      :close-on-press-escape="false"
      :close-on-click-modal="drawerData.isFinished"
    >
      <el-descriptions direction="vertical" :column="1" border>
        <el-descriptions-item label="文件名">{{ drawerData.fileName }}</el-descriptions-item>
        <el-descriptions-item label="文件大小">{{ drawerData.total }}</el-descriptions-item>
        <el-descriptions-item :label="drawerData.isUpload ? '已上传大小' : '已下载大小'">{{ drawerData.loaded }}</el-descriptions-item>
        <el-descriptions-item :label="drawerData.isUpload ? '上传结果' : '下载结果'">
          <el-tag v-if="drawerData.finished === '成功'" type="success" disable-transitions>成功</el-tag>
          <el-tag v-else-if="drawerData.finished === '失败'" type="danger" disable-transitions>失败</el-tag>
          <el-tag v-else disable-transitions>进行中</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="drawerData.isUpload ? '上传进度' : '下载进度'">
          <el-progress :percentage="drawerData.percent" :stroke-width="20" striped striped-flow :duration="20" />
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref, onMounted } from "vue";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import FileApi from "@/api/file";
import { saveAs } from "file-saver";
import { Upload } from "@/utils/upload";
import { useTokenStore } from "@/store/token";
import { storeToRefs } from "pinia";
import { FolderChecked, Lock, User } from "@element-plus/icons-vue";
import { formatTime, uuid } from "@/utils";
import { AxiosProgressEvent } from "axios";

const hostUrl = ref(location.origin);
const tokenStore = useTokenStore();
const { accessToken, publicAuth, protectedAuth, privateAuth } = storeToRefs(tokenStore);
const parentDir = ref("public");
const loading = ref(false);
const fileList: Ref<FileInfo[]> = ref([]);
const pathList: Ref<string[]> = ref([]);
const drawerData = ref({
  visible: false,
  id: "",
  fileName: "",
  isUpload: false,
  isFinished: false,
  finished: "",
  loaded: "",
  total: "",
  percent: 0,
});

onMounted(() => {
  queryFileList();
});

/**
 * 判断是否有更新权限
 */
const updateAuth = () => {
  if (parentDir.value === "public" && publicAuth.value) {
    return true;
  }
  if (parentDir.value === "protected" && protectedAuth.value) {
    return true;
  }
  if (parentDir.value === "private" && privateAuth.value) {
    return true;
  }
  return false;
};

/**
 * 文件分类变更
 * @param menu
 */
const changeType = (menu: string) => {
  parentDir.value = menu;
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
  loading.value = true;
  FileApi.list(parentDir.value, pathList.value.join("/"))
    .then((res) => {
      fileList.value = res.data;
    })
    .finally(() => {
      loading.value = false;
    });
};

/**
 * 分享文件
 * @param row
 */
const shareClick = (row: FileInfo) => {
  ElMessageBox.prompt("请输入分享时长(小时)，限制在1-720小时之间", "分享文件", {
    confirmButtonText: "创建",
    cancelButtonText: "取消",
  }).then(({ value }) => {
    if (!/^[1-9]\d*$/.test(value)) {
      ElMessage.warning("请输入1-720之间的正整数");
      return;
    }
    let shareHours = parseInt(value);
    if (shareHours > 720) {
      ElMessage.warning("请输入1-720之间的正整数");
      return;
    }
    loading.value = true;
    FileApi.share(parentDir.value, pathList.value.join("/"), row.name, shareHours)
      .then((res) => {
        ElNotification({
          type: "success",
          title: "文件分享成功",
          dangerouslyUseHTMLString: true,
          duration: 0,
          message: `<div>文件名：${row.name}</div>
          <div>有效期：${shareHours}小时</div>
          <div>链接：</div>
          <div style="color: #3d5eb9">${hostUrl.value}/api/open/file/share/${res.data}</div>`,
        });
      })
      .finally(() => {
        loading.value = false;
      });
  });
};

/**
 * 上传文件
 */
const uploadClick = () => {
  Upload.openFiles().then((fileList) => {
    if (fileList[0].size >= 1024 * 1024 * 1024) {
      ElMessage.error("文件大小不可超过1GB");
      return;
    }
    const id = uuid();
    openDrawer(id);
    FileApi.upload(parentDir.value, pathList.value.join("/"), fileList[0], onProgress, id)
      .then(() => {
        ElMessage.success("上传成功");
        if (id === drawerData.value.id) {
          drawerData.value.isFinished = true;
          drawerData.value.finished = "成功";
        }
        queryFileList();
      })
      .catch(() => {
        if (id === drawerData.value.id) {
          drawerData.value.isFinished = true;
          drawerData.value.finished = "失败";
        }
      });
  });
};

/**
 * 下载
 */
const downloadClick = (row: FileInfo) => {
  const fileName = row.name;
  const id = uuid();
  openDrawer(id);
  FileApi.download(parentDir.value, pathList.value.join("/"), fileName, onProgress, id)
    .then((res) => {
      saveAs(res, fileName);
      if (id === drawerData.value.id) {
        drawerData.value.isFinished = true;
        drawerData.value.finished = "成功";
      }
    })
    .catch(() => {
      if (id === drawerData.value.id) {
        drawerData.value.isFinished = true;
        drawerData.value.finished = "失败";
      }
    });
};

/**
 * 上传/下载进度
 * @param progressEvent
 * @param fileName
 * @param isUpload
 * @param id
 */
const onProgress = (progressEvent: AxiosProgressEvent, fileName: string, isUpload: boolean, id: string) => {
  // 根据id判断进度是否显示
  if (drawerData.value.id !== id) {
    return;
  }
  drawerData.value.isUpload = isUpload;
  drawerData.value.fileName = fileName;
  if (progressEvent.loaded) {
    drawerData.value.loaded = Upload.formatFileSize(progressEvent.loaded);
  }
  if (progressEvent.total) {
    drawerData.value.total = Upload.formatFileSize(progressEvent.total);
  }
  if (progressEvent.progress) {
    drawerData.value.percent = Math.round(progressEvent.progress * 100);
  }
};

/**
 * 打开进度抽屉
 * @param id
 */
const openDrawer = (id: string) => {
  drawerData.value.id = id;
  drawerData.value.isFinished = false;
  drawerData.value.finished = "进行中";
  drawerData.value.isUpload = false;
  drawerData.value.fileName = "";
  drawerData.value.loaded = "";
  drawerData.value.total = "";
  drawerData.value.percent = 0;
  drawerData.value.visible = true;
};

/**
 * 创建目录
 */
const createFolderClick = () => {
  ElMessageBox.prompt("请输入目录名称", "创建目录", {
    confirmButtonText: "创建",
    cancelButtonText: "取消",
  }).then(({ value }) => {
    loading.value = true;
    FileApi.folder(parentDir.value, pathList.value.join("/"), value)
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
 * 删除
 */
const deleteClick = (row: FileInfo) => {
  ElMessageBox.confirm(`是否删除${row.isDir ? "目录" : "文件"}：${row.name}？${row.isDir ? "这将会删除目录内的全部文件。" : ""}`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    loading.value = true;
    FileApi.delete(parentDir.value, pathList.value.join("/"), row.name)
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
    height: calc(100% - 10px);
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
