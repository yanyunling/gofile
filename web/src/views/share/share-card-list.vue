<template>
  <div class="share-card-list">
    <div v-for="(item, index) in shareData" :key="index" class="share-card">
      <div class="card-header">
        <div class="header-main">
          <div class="file-name">{{ item.name }}</div>
          <div class="file-meta">
            <span>{{ item.parentDir }}</span>
            <span v-if="item.path"> · {{ item.path }}</span>
          </div>
        </div>
        <el-button type="danger" size="small" @click="$emit('delete', item)">删除</el-button>
      </div>
      <div class="card-body">
        <div class="info-row">
          <span class="label">下载链接：</span>
          <span class="link-text" @click="$emit('copy-link', item.id)">{{ shareUrl + item.id }}</span>
        </div>
        <div class="info-row">
          <span class="label">分享时长：</span>
          <span class="text">{{ item.shareHours }}小时</span>
        </div>
        <div class="info-row">
          <span class="label">起始时间：</span>
          <span class="text">{{ formatTime(item.startTime, "YYYY-MM-DD HH:mm:ss") }}</span>
        </div>
        <div class="info-row">
          <span class="label">截止时间：</span>
          <span class="text">{{ formatTime(item.endTime, "YYYY-MM-DD HH:mm:ss") }}</span>
        </div>
        <div v-if="item.username" class="info-row">
          <span class="label">操作用户：</span>
          <span class="text">{{ item.username }}</span>
        </div>
      </div>
    </div>
    <div v-if="shareData.length === 0" class="empty-tip">
      <el-empty description="暂无分享记录" :image-size="80" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { formatTime } from "@/utils";

interface Share {
  id: string;
  parentDir: string;
  path?: string;
  name: string;
  username?: string;
  shareHours: number;
  startTime: number;
  endTime: number;
}

interface Props {
  shareData: Share[];
  shareUrl: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "delete", row: Share): void;
  (e: "copy-link", id: string): void;
}>();
</script>

<style lang="scss" scoped>
.share-card-list {
  width: calc(100% - 20px);
  padding: 10px;
  background: #fff;
  .share-card {
    background: #fff;
    border-radius: 8px;
    margin-bottom: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    border: 1px solid #eee;
    overflow: hidden;
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 15px;
      background: #fafafa;
      border-bottom: 1px solid #eee;
      .header-main {
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
      .el-button {
        flex-shrink: 0;
        margin-left: 10px;
      }
    }
    .card-body {
      padding: 12px 15px;
      .info-row {
        margin-bottom: 10px;
        display: flex;
        align-items: flex-start;
        gap: 8px;
        &:last-child {
          margin-bottom: 0;
        }
        .label {
          font-size: 12px;
          color: #666;
          font-weight: 500;
          white-space: nowrap;
          flex-shrink: 0;
        }
        .link-text {
          font-size: 13px;
          color: #3d5eb9;
          word-break: break-all;
          line-height: 1.6;
          cursor: pointer;
          &:active {
            opacity: 0.7;
          }
        }
        .text {
          font-size: 13px;
          color: #333;
          line-height: 1.6;
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
