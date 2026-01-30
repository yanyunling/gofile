<template>
  <el-dialog
    :model-value="props.visible"
    title="用户管理"
    append-to-body
    destroy-on-close
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    center
    fullscreen
    :before-close="beforeClose"
  >
    <div class="page-user">
      <div class="filter-view">
        <el-input class="filter-input" v-model="tableCondition.condition" placeholder="全局搜索" clearable></el-input>
        <el-button type="primary" :icon="Search" @click="tablePageCurrentChange(1)" :loading="tableLoading">查询</el-button>
      </div>
      <el-table class="table-view" ref="tableRef" :data="tableData" height="100%" stripe border v-loading="tableLoading">
        <el-table-column prop="" label="序号" align="center" width="60">
          <template #default="scope">
            {{ (tableCondition.page.current - 1) * tableCondition.page.size + scope.$index + 1 }}
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" align="center" />
        <el-table-column prop="enable" label="状态" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.enable === 1" type="success" disable-transitions>启用</el-tag>
            <el-tag v-else type="danger" disable-transitions>停用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="hasUpdate" label="可以更新" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.hasUpdate === 1" type="success" disable-transitions>是</el-tag>
            <el-tag v-else type="danger" disable-transitions>否</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="" label="" align="center" width="180">
          <template #header>
            <el-button type="primary" link size="small" @click="addClick">新增</el-button>
          </template>
          <template #default="scope">
            <el-button type="primary" link size="small" @click="editClick(scope.row)">修改</el-button>
            <el-button type="danger" link size="small" @click="deleteClick(scope.row)">删除</el-button>
            <el-button type="success" link size="small" @click="resetPasswordClick(scope.row)">重置密码</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="total, sizes, prev, pager, next, jumper"
        :pageSizes="[10, 20, 50, 100, 200, 500]"
        v-model:pageSize="tableCondition.page.size"
        v-model:currentPage="tableCondition.page.current"
        :total="tableTotal"
        @size-change="tablePageSizeChange"
        @current-change="tablePageCurrentChange"
      ></el-pagination>
      <save-dialog v-model:visible="dialogVisible" :form-data="dialogFormData" @save="queryTableData" />
      <reset-password-dialog v-model:visible="resetPasswordDialogVisible" :form-data="dialogFormData" />
    </div>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, Ref, watch, nextTick } from "vue";
import { ElMessage, ElMessageBox, ElTable } from "element-plus";
import { Search } from "@element-plus/icons-vue";
import UserApi from "@/api/user";
import saveDialog from "./save-dialog.vue";
import resetPasswordDialog from "./reset-password-dialog.vue";

const emit = defineEmits<{ "update:visible": [visible: boolean] }>();

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
});

const tableData: Ref<User[]> = ref([]);
const tableTotal = ref(0);
const tableLoading = ref(false);
const tableCondition: Ref<PageCondition<string>> = ref({
  page: {
    current: 1,
    size: 10,
  },
  condition: "",
});
const tableRef = ref<InstanceType<typeof ElTable>>();
const dialogVisible = ref(false);
const resetPasswordDialogVisible = ref(false);
const dialogFormData: Ref<User> = ref();

watch(
  () => props.visible,
  (val) => {
    if (val) {
      queryTableData();
    }
  },
);

/**
 * 查询表格数据
 */
const queryTableData = () => {
  tableLoading.value = true;
  UserApi.page(tableCondition.value)
    .then((res) => {
      tableData.value = res.data.records;
      tableTotal.value = res.data.total;
      nextTick(() => {
        tableRef.value?.setScrollTop(0);
      });
    })
    .finally(() => {
      tableLoading.value = false;
    });
};

/**
 * 点击新增
 */
const addClick = () => {
  dialogFormData.value = {
    id: "",
    username: "",
    password: "",
    enable: 1,
    hasUpdate: 0,
  };
  dialogVisible.value = true;
};

/**
 * 点击修改
 * @param user
 */
const editClick = (user: User) => {
  dialogFormData.value = user;
  dialogVisible.value = true;
};

/**
 * 点击重置密码
 * @param user
 */
const resetPasswordClick = (user: User) => {
  dialogFormData.value = user;
  resetPasswordDialogVisible.value = true;
};

/**
 * 删除记录
 * @param row
 */
const deleteClick = (row: User) => {
  ElMessageBox.confirm("是否删除用户：" + row.username + "？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    tableLoading.value = true;
    UserApi.delete(row.id)
      .then(() => {
        ElMessage.success("删除成功");
        queryTableData();
      })
      .catch(() => {
        tableLoading.value = false;
      });
  });
};

/**
 * 每页显示条数变化
 * @param size
 */
const tablePageSizeChange = (size: number) => {
  tableCondition.value.page.current = 1;
  tableCondition.value.page.size = size;
  queryTableData();
};

/**
 * 当前页码变化
 * @param current
 */
const tablePageCurrentChange = (current: number) => {
  tableCondition.value.page.current = current;
  queryTableData();
};

/**
 * 弹窗关闭
 */
const beforeClose = () => {
  emit("update:visible", false);
};
</script>

<style lang="scss" scoped>
.page-user {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: fixed;
  width: 100%;
  height: calc(100% - 50px);
  left: 0;
  right: 0;
  top: 50px;
  bottom: 0;
  .filter-view {
    max-width: 100%;
    margin: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    .filter-input {
      width: 200px;
    }
  }
  .table-view {
    .el-scrollbar__wrap {
      display: flex;
    }
  }
}
</style>
