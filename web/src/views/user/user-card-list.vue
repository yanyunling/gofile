<template>
  <div class="user-card-list">
    <div class="action-bar">
      <el-button type="primary" @click="$emit('add')">新增</el-button>
    </div>
    <div v-for="(item, index) in userData" :key="index" class="user-card">
      <div class="card-header">
        <div class="username-section">
          <el-icon class="user-icon"><User /></el-icon>
          <span class="username-text">{{ item.username }}</span>
        </div>
        <el-tag :type="item.enable === 1 ? 'success' : 'danger'" size="small" disable-transitions>
          {{ item.enable === 1 ? "启用" : "停用" }}
        </el-tag>
      </div>
      <div class="card-body">
        <div class="permission-row">
          <span class="label">公开文件：</span>
          <el-tag :type="item.publicAuth === 1 ? 'success' : 'info'" size="small" disable-transitions>
            {{ item.publicAuth === 1 ? "有权限" : "无权限" }}
          </el-tag>
        </div>
        <div class="permission-row">
          <span class="label">保护文件：</span>
          <el-tag :type="item.protectedAuth === 1 ? 'success' : 'info'" size="small" disable-transitions>
            {{ item.protectedAuth === 1 ? "有权限" : "无权限" }}
          </el-tag>
        </div>
        <div class="permission-row">
          <span class="label">私有文件：</span>
          <el-tag :type="item.privateAuth === 1 ? 'success' : 'info'" size="small" disable-transitions>
            {{ item.privateAuth === 1 ? "有权限" : "无权限" }}
          </el-tag>
        </div>
      </div>
      <div class="card-footer">
        <el-button type="primary" size="small" @click="$emit('edit', item)">修改</el-button>
        <el-button type="success" size="small" @click="$emit('reset-password', item)">重置密码</el-button>
        <el-button type="danger" size="small" @click="$emit('delete', item)">删除</el-button>
      </div>
    </div>
    <div v-if="userData.length === 0" class="empty-tip">
      <el-empty description="暂无用户" :image-size="80" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { User } from "@element-plus/icons-vue";

interface User {
  id: string;
  username: string;
  enable: number;
  publicAuth: number;
  protectedAuth: number;
  privateAuth: number;
}

interface Props {
  userData: User[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "edit", row: User): void;
  (e: "reset-password", row: User): void;
  (e: "delete", row: User): void;
  (e: "add"): void;
}>();
</script>

<style lang="scss" scoped>
.user-card-list {
  width: calc(100% - 20px);
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
    }
  }
  .user-card {
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
      .username-section {
        display: flex;
        align-items: center;
        gap: 10px;
        flex: 1;
        min-width: 0;
        .user-icon {
          font-size: 24px;
          color: #3d5eb9;
          flex-shrink: 0;
        }
        .username-text {
          font-size: 16px;
          font-weight: 500;
          color: #333;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
      .el-tag {
        flex-shrink: 0;
        margin-left: 10px;
      }
    }
    .card-body {
      padding: 12px 15px;
      .permission-row {
        margin-bottom: 10px;
        display: flex;
        align-items: center;
        gap: 10px;
        &:last-child {
          margin-bottom: 0;
        }
        .label {
          font-size: 13px;
          color: #666;
          font-weight: 500;
          width: 70px;
          flex-shrink: 0;
        }
        .el-tag {
          flex: 1;
        }
      }
    }
    .card-footer {
      display: flex;
      gap: 8px;
      padding: 10px 15px;
      background: #fafafa;
      border-top: 1px solid #eee;
      .el-button {
        flex: 1;
      }
    }
  }
  .empty-tip {
    padding: 60px 0;
    text-align: center;
  }
}
</style>
