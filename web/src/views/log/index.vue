<template>
  <el-dialog
    :model-value="props.visible"
    title="操作日志"
    append-to-body
    destroy-on-close
    :close-on-click-modal="false"
    center
    fullscreen
    :before-close="beforeClose"
  >
    <div class="page-log">
      <div class="filter-view">
        <el-button type="primary" :icon="Filter" @click="drawerVisible = true" :loading="tableLoading">条件查询</el-button>
      </div>
      <el-table class="table-view" ref="tableRef" :data="tableData" height="100%" stripe border size="small" v-loading="tableLoading">
        <el-table-column prop="" label="序号" align="center" width="60">
          <template #default="scope">
            {{ (tableCondition.page.current - 1) * tableCondition.page.size + scope.$index + 1 }}
          </template>
        </el-table-column>
        <el-table-column prop="level" label="级别" align="center" width="80">
          <template #default="scope">
            <el-tag :type="levelToTag(scope.row.level)" disable-transitions>{{ scope.row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" align="center" width="150" />
        <el-table-column prop="content" label="内容" align="left" header-align="center" />
        <el-table-column prop="username" label="操作用户" align="center" width="150" />
        <el-table-column prop="createTime" label="时间" align="center" width="160">
          <template #default="scope"> {{ formatTime(scope.row.createTime, "YYYY-MM-DD HH:mm:ss") }} </template>
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
      header-class="log-filter-drawer-header"
      body-class="log-filter-drawer-body"
    >
      <el-form inline>
        <el-form-item label="级别">
          <el-select class="filter-input" v-model="tableCondition.condition.level" placeholder="级别筛选" clearable>
            <el-option label="Debug" value="Debug" />
            <el-option label="Info" value="Info" />
            <el-option label="Warn" value="Warn" />
            <el-option label="Error" value="Error" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题">
          <el-input class="filter-input" v-model="tableCondition.condition.title" placeholder="标题筛选" clearable />
        </el-form-item>
        <el-form-item label="内容">
          <el-input class="filter-input" v-model="tableCondition.condition.content" placeholder="内容筛选" clearable />
        </el-form-item>
        <el-form-item label="操作用户" v-if="isAdmin">
          <el-input class="filter-input" v-model="tableCondition.condition.username" placeholder="用户筛选" clearable />
        </el-form-item>
        <el-form-item label="时间">
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
import { ElTable } from "element-plus";
import { Search, Filter } from "@element-plus/icons-vue";
import LogApi from "@/api/log";
import { formatTime } from "@/utils";
import { useTokenStore } from "@/store/token";
import { storeToRefs } from "pinia";

const emit = defineEmits<{ "update:visible": [visible: boolean] }>();

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
});

const tokenStore = useTokenStore();
const { isAdmin } = storeToRefs(tokenStore);
const tableData: Ref<Log[]> = ref([]);
const tableTotal = ref(0);
const tableLoading = ref(false);
const tableCondition: Ref<PageCondition<LogCondition>> = ref({
  page: {
    current: 1,
    size: 100,
  },
  condition: {
    title: "",
    content: "",
    level: "",
    username: "",
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
  LogApi.page(tableCondition.value)
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
 * 日志等级转tag颜色
 * @param level
 */
const levelToTag = (level: string) => {
  if (level === "Info") {
    return "primary";
  } else if (level === "Warn") {
    return "warning";
  } else if (level === "Error") {
    return "danger";
  }
  return "info";
};
</script>

<style lang="scss">
.page-log {
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
  }
  .table-view {
    .el-scrollbar__wrap {
      display: flex;
    }
  }
}
.log-filter-drawer-header {
  margin-bottom: 0;
  .el-drawer__title {
    text-align: center;
  }
}
.log-filter-drawer-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  .filter-input {
    width: 180px;
  }
  .search-button {
    margin-top: 10px;
  }
}
</style>
