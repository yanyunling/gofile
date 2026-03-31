<template>
  <div class="log-card-list">
    <div v-for="(item, index) in logData" :key="index" class="log-card">
      <div class="card-header">
        <div class="header-left">
          <el-tag :type="getLevelType(item.level)" size="small" disable-transitions>
            {{ item.level }}
          </el-tag>
          <span class="title-text">{{ item.title }}</span>
        </div>
        <span class="time-text">{{ formatTime(item.createTime, "YYYY-MM-DD HH:mm:ss") }}</span>
      </div>
      <div class="card-content">
        <div class="content-row">
          <span class="label">内容：</span>
          <span class="text">{{ item.content }}</span>
        </div>
        <div v-if="item.username" class="content-row">
          <span class="label">操作用户：</span>
          <span class="text">{{ item.username }}</span>
        </div>
      </div>
    </div>
    <div v-if="logData.length === 0" class="empty-tip">
      <el-empty description="暂无日志" :image-size="80" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { formatTime } from "@/utils";

interface Log {
  level: string;
  title: string;
  content: string;
  username?: string;
  createTime: number;
}

interface Props {
  logData: Log[];
}

const props = defineProps<Props>();

const getLevelType = (level: string) => {
  if (level === "Info") return "primary";
  if (level === "Warn") return "warning";
  if (level === "Error") return "danger";
  return "info";
};
</script>

<style lang="scss" scoped>
.log-card-list {
  width: calc(100% - 20px);
  padding: 10px;
  background: #fff;

  .log-card {
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

      .header-left {
        display: flex;
        align-items: center;
        gap: 10px;
        flex: 1;
        min-width: 0;

        .title-text {
          font-size: 14px;
          font-weight: 500;
          color: #333;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }

      .time-text {
        font-size: 12px;
        color: #999;
        flex-shrink: 0;
        margin-left: 10px;
      }
    }

    .card-content {
      padding: 12px 15px;

      .content-row {
        margin-bottom: 8px;
        display: flex;
        align-items: flex-start;
        gap: 5px;

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

        .text {
          font-size: 13px;
          color: #333;
          word-break: break-all;
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
