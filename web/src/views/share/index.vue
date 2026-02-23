<template>
  <el-dialog
    :model-value="props.visible"
    title="分享记录"
    append-to-body
    destroy-on-close
    :close-on-click-modal="false"
    center
    fullscreen
    :before-close="beforeClose"
  >
    <div class="page-share">
      <el-button class="filter-button" type="primary" :icon="Filter" @click="drawerVisible = true" :loading="tableLoading">条件查询</el-button>
      <el-table class="table-view" ref="tableRef" :data="tableData" height="100%" stripe border size="small" v-loading="tableLoading">
        <el-table-column prop="" label="序号" align="center" width="60">
          <template #default="scope">
            {{ (tableCondition.page.current - 1) * tableCondition.page.size + scope.$index + 1 }}
          </template>
        </el-table-column>
        <el-table-column prop="parentDir" label="根目录" align="center" width="150" />
        <el-table-column prop="path" label="父级目录" align="center" />
        <el-table-column prop="name" label="文件名" align="center" />
        <el-table-column prop="id" label="下载链接" align="center">
          <template #default="scope">
            <span style="cursor: pointer; color: #3d5eb9" @click="copyLinkClick(scope.row.id)">{{ shareUrl + scope.row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="操作用户" align="center" width="150" />
        <el-table-column prop="shareHours" label="分享小时数" align="center" width="80" />
        <el-table-column prop="startTime" label="起始时间" align="center" width="160">
          <template #default="scope"> {{ formatTime(scope.row.startTime, "YYYY-MM-DD HH:mm:ss") }} </template>
        </el-table-column>
        <el-table-column prop="endTime" label="截止时间" align="center" width="160">
          <template #default="scope"> {{ formatTime(scope.row.endTime, "YYYY-MM-DD HH:mm:ss") }} </template>
        </el-table-column>
        <el-table-column prop="" label="操作" align="center" width="80">
          <template #default="scope">
            <el-button type="danger" link size="small" @click="deleteClick(scope.row)">删除</el-button>
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
    </div>
    <el-drawer
      v-model="drawerVisible"
      title="条件查询"
      direction="ttb"
      append-to-body
      size="30%"
      header-class="filter-drawer-header"
      body-class="filter-drawer-body"
    >
      <el-form inline>
        <el-form-item label="根目录">
          <el-select class="filter-input" v-model="tableCondition.condition.parentDir" placeholder="根目录筛选" clearable>
            <el-option label="public" value="public" />
            <el-option label="protected" value="protected" />
            <el-option label="private" value="private" />
          </el-select>
        </el-form-item>
        <el-form-item label="父级目录">
          <el-input class="filter-input" v-model="tableCondition.condition.path" placeholder="父级目录筛选" clearable />
        </el-form-item>
        <el-form-item label="文件名">
          <el-input class="filter-input" v-model="tableCondition.condition.name" placeholder="文件名筛选" clearable />
        </el-form-item>
        <el-form-item label="操作用户" v-if="isAdmin">
          <el-input class="filter-input" v-model="tableCondition.condition.username" placeholder="用户筛选" clearable />
        </el-form-item>
        <el-form-item label="分享小时数">
          <el-input-number class="filter-input" v-model="tableCondition.condition.shareHours" :min="1" :max="720" />
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            value-format="x"
            range-separator="至"
            start-placeholder="起始时间"
            end-placeholder="截止时间"
          />
        </el-form-item>
      </el-form>
      <el-button class="search-button" type="primary" :icon="Search" @click="tablePageCurrentChange(1)" :loading="tableLoading">查询</el-button>
    </el-drawer>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, Ref, watch, nextTick } from "vue";
import { ElTable, ElMessage, ElMessageBox } from "element-plus";
import { Search, Filter } from "@element-plus/icons-vue";
import ShareApi from "@/api/share";
import { formatTime } from "@/utils";
import { useTokenStore } from "@/store/token";
import { storeToRefs } from "pinia";
import { host, context } from "@/config";
import copy from "copy-to-clipboard";

const emit = defineEmits<{ "update:visible": [visible: boolean] }>();

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
});

const hostUrl = process.env.NODE_ENV === "production" ? location.origin : host;
const shareUrl = ref(`${hostUrl}${context}/open/file/share/`);
const tokenStore = useTokenStore();
const { isAdmin } = storeToRefs(tokenStore);
const tableData: Ref<Share[]> = ref([]);
const tableTotal = ref(0);
const tableLoading = ref(false);
const tableCondition: Ref<PageCondition<Share>> = ref({
  page: {
    current: 1,
    size: 100,
  },
  condition: {
    id: "",
    parentDir: "",
    path: "",
    name: "",
    username: "",
    shareHours: undefined,
    startTime: 0,
    endTime: 0,
  },
});
const tableRef = ref<InstanceType<typeof ElTable>>();
const drawerVisible = ref(false);
const dateRange = ref([]);

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
  if (dateRange.value && dateRange.value.length === 2) {
    tableCondition.value.condition.startTime = dateRange.value[0];
    tableCondition.value.condition.endTime = dateRange.value[1] + 86400000 - 1;
  } else {
    tableCondition.value.condition.startTime = 0;
    tableCondition.value.condition.endTime = 0;
  }
  tableLoading.value = true;
  ShareApi.page(tableCondition.value)
    .then((res) => {
      tableData.value = res.data.records;
      tableTotal.value = res.data.total;
      nextTick(() => {
        tableRef.value?.setScrollTop(0);
      });
    })
    .finally(() => {
      tableLoading.value = false;
      drawerVisible.value = false;
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

/**
 * 点击下载链接
 * @param id
 */
const copyLinkClick = (id: string) => {
  let url = shareUrl.value + id;
  const result = copy(url);
  if (result) {
    ElMessage.success("下载链接已复制到剪贴板");
  } else {
    ElMessage.error("复制到剪贴板失败，链接：" + url);
  }
};

/**
 * 删除记录
 * @param row
 */
const deleteClick = (row: Share) => {
  ElMessageBox.confirm("是否删除分享记录：" + row.name + "？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(() => {
      tableLoading.value = true;
      ShareApi.delete(row.id)
        .then(() => {
          ElMessage.success("删除成功");
          queryTableData();
        })
        .catch(() => {
          tableLoading.value = false;
        });
    })
    .catch(() => {});
};
</script>

<style lang="scss">
.page-share {
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
  .filter-button {
    margin: 10px;
    margin-left: -20px;
  }
  .table-view {
    .el-scrollbar__wrap {
      display: flex;
    }
  }
}
</style>
