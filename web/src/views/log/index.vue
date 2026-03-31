<template>
  <div class="page-log" :style="isMobile ? { 'overflow-y': 'auto' } : {}">
    <div class="filter-view">
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
            v-if="!isMobile"
            style="width: 240px"
            v-model="dateRange"
            type="daterange"
            value-format="x"
            range-separator="至"
            start-placeholder="起始时间"
            end-placeholder="截止时间"
            @change="dateRangeChange"
            @clear="dateClear"
          />
          <el-date-picker
            v-else
            style="width: 180px"
            v-model="dateMonth"
            type="month"
            value-format="x"
            placeholder="选择月份"
            @change="dateMonthChange"
            @clear="dateClear"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="tablePageCurrentChange(1)" :loading="tableLoading">查询</el-button>
          <el-button @click="resetCondition">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-table v-if="!isMobile" class="table-view" ref="tableRef" :data="tableData" height="100%" stripe border size="small" v-loading="tableLoading">
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
    <log-card-list v-else :log-data="tableData" />
    <el-pagination
      background
      :size="isMobile ? 'small' : 'default'"
      :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper'"
      v-model:pageSize="tableCondition.page.size"
      v-model:currentPage="tableCondition.page.current"
      :total="tableTotal"
      @size-change="tablePageSizeChange"
      @current-change="tablePageCurrentChange"
    ></el-pagination>
  </div>
</template>

<script lang="ts" setup>
import { ref, Ref, onMounted, nextTick } from "vue";
import { ElTable } from "element-plus";
import { Search } from "@element-plus/icons-vue";
import LogApi from "@/api/log";
import { formatTime } from "@/utils";
import { useTokenStore } from "@/store/token";
import { storeToRefs } from "pinia";
import { useIsMobile } from "@/utils/useIsMobile";
import LogCardList from "./log-card-list.vue";
import { getMonthLastMilliSecond } from "@/utils/time";

const isMobile = useIsMobile();
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
const dateRange = ref([]);
const dateMonth = ref("");

onMounted(() => {
  if (isMobile.value) {
    tableCondition.value.page.size = 10;
  }
  queryTableData();
});

/**
 * 查询表格数据
 */
const queryTableData = () => {
  tableLoading.value = true;
  LogApi.page(tableCondition.value)
    .then((res) => {
      tableData.value = res.data.records;
      tableTotal.value = res.data.total;
      nextTick(() => {
        if (!isMobile.value) {
          tableRef.value?.setScrollTop(0);
        } else {
          document.getElementsByClassName("page-log")[0]?.scrollTo(0, 0);
        }
      });
    })
    .finally(() => {
      tableLoading.value = false;
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
 * 日期范围选择变化
 * @param date
 */
const dateRangeChange = (date: number[]) => {
  dateMonth.value = "";
  if (date && date.length === 2) {
    tableCondition.value.condition.startTime = date[0];
    tableCondition.value.condition.endTime = date[1] + 86400000 - 1;
  } else {
    tableCondition.value.condition.startTime = 0;
    tableCondition.value.condition.endTime = 0;
  }
};

/**
 * 月份选择变化
 * @param date
 */
const dateMonthChange = (date: number) => {
  dateRange.value = [];
  tableCondition.value.condition.startTime = date;
  tableCondition.value.condition.endTime = getMonthLastMilliSecond(date);
};

/**
 * 日期清空
 */
const dateClear = () => {
  tableCondition.value.condition.startTime = 0;
  tableCondition.value.condition.endTime = 0;
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

/**
 * 重置查询条件
 */
const resetCondition = () => {
  tableCondition.value = {
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
  };
  dateRange.value = [];
  dateMonth.value = "";
  if (isMobile.value) {
    tableCondition.value.page.size = 10;
  }
  queryTableData();
};
</script>

<style lang="scss" scoped>
.page-log {
  display: flex;
  flex-direction: column;
  align-items: center;
  .filter-view {
    width: 100%;
    margin-top: 10px;
    .el-form {
      display: flex;
      flex-wrap: wrap;
      justify-content: center;
      .el-form-item {
        margin-bottom: 10px;
      }
    }
    .filter-input {
      width: 180px;
    }
  }
  .table-view {
    .el-scrollbar__wrap {
      display: flex;
    }
  }
}
</style>
