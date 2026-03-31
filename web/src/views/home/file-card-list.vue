<template>
  <div class="file-card-list">
    <div v-if="hasUpdateAuth" class="action-bar">
      <el-button type="primary" @click="$emit('upload')">上传</el-button>
      <el-button type="primary" @click="$emit('create-folder')">创建目录</el-button>
    </div>
    <div v-for="(item, index) in fileList" :key="index" class="file-card">
      <div class="card-header" @click="handleCardClick(item)">
        <el-icon class="file-icon">
          <Folder v-if="item.isDir" />
          <Document v-else />
        </el-icon>
        <div class="file-info">
          <div class="file-name">{{ item.name }}</div>
          <div class="file-meta">
            <span>{{ item.isDir ? "目录" : getFileExtension(item.name) }}</span>
            <span v-if="!item.isDir"> · {{ formatFileSize(item.size) }}</span>
          </div>
        </div>
        <el-icon v-if="item.isDir" class="arrow-icon"><ArrowRight /></el-icon>
      </div>
      <div class="card-footer">
        <div class="time-text">{{ formatTime(item.updateTime) }}</div>
        <div class="action-buttons">
          <el-button v-if="!item.isDir" type="primary" size="small" @click="$emit('download', item)"> 下载 </el-button>
          <el-button v-if="!item.isDir && showShare" type="success" size="small" @click="$emit('share', item)"> 分享 </el-button>
          <el-button v-if="hasUpdateAuth" type="danger" size="small" @click="$emit('delete', item)"> 删除 </el-button>
        </div>
      </div>
    </div>
    <div v-if="fileList.length === 0" class="empty-tip">
      <el-empty description="暂无文件" :image-size="80" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Folder, Document, ArrowRight } from "@element-plus/icons-vue";
import { Upload } from "@/utils/upload";
import { formatTime } from "@/utils";

interface FileInfo {
  name: string;
  size: number;
  updateTime: number;
  isDir: boolean;
}

interface Props {
  fileList: FileInfo[];
  hasUpdateAuth: boolean;
  showShare?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  showShare: false,
});

const emit = defineEmits<{
  (e: "folder-click", row: FileInfo): void;
  (e: "download", row: FileInfo): void;
  (e: "share", row: FileInfo): void;
  (e: "delete", row: FileInfo): void;
  (e: "upload"): void;
  (e: "create-folder"): void;
}>();

const getFileExtension = (fileName: string) => {
  return Upload.getFileExtension(fileName);
};

const formatFileSize = (size: number) => {
  return Upload.formatFileSize(size);
};

const handleCardClick = (item: FileInfo) => {
  if (item.isDir) {
    emit("folder-click", item);
  }
};
</script>

<style lang="scss" scoped>
.file-card-list {
  padding: 10px;
  background: #fff;

  .action-bar {
    display: flex;
    justify-content: center;
    gap: 10px;
    margin-bottom: 15px;
    padding-bottom: 15px;
    border-bottom: 1px solid #eee;

    .el-button {
      flex: 1;
      max-width: 140px;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 5px;
    }
  }

  .file-card {
    background: #fff;
    border-radius: 8px;
    margin-bottom: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    border: 1px solid #eee;
    overflow: hidden;
    transition: all 0.3s;

    &:active {
      background: #f5f5f5;
    }

    .card-header {
      display: flex;
      align-items: center;
      padding: 12px 15px;
      gap: 12px;
      cursor: pointer;

      .file-icon {
        font-size: 28px;
        color: #3d5eb9;
        flex-shrink: 0;
      }

      .file-info {
        flex: 1;
        min-width: 0;

        .file-name {
          font-size: 15px;
          font-weight: 500;
          color: #333;
          margin-bottom: 4px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .file-meta {
          font-size: 12px;
          color: #999;
          display: flex;
          align-items: center;
          gap: 5px;
        }
      }

      .arrow-icon {
        font-size: 18px;
        color: #ccc;
        flex-shrink: 0;
      }
    }

    .card-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 10px 15px;
      background: #fafafa;
      border-top: 1px solid #eee;

      .time-text {
        font-size: 12px;
        color: #999;
      }

      .action-buttons {
        display: flex;
        gap: 8px;

        .el-button {
          padding: 5px 10px;
          font-size: 12px;
        }
      }
    }
  }

  .empty-tip {
    padding: 60px 0;
    text-align: center;
  }
}
</style>
