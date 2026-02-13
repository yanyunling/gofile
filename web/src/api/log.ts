import request from "@/utils/request";

class LogApi {
  /**
   * 分页查询日志记录
   * @param pageCondition
   * @returns
   */
  page(pageCondition: PageCondition<LogCondition>) {
    return request<PageResult<Log>>({
      method: "post",
      url: "/log/page",
      data: pageCondition,
    });
  }
}

export default new LogApi();
