import request from "@/utils/request";

class ShareApi {
  /**
   * 分页查询分享记录
   * @param pageCondition
   * @returns
   */
  page(pageCondition: PageCondition<Share>) {
    return request<PageResult<Share>>({
      method: "post",
      url: "/share/page",
      data: pageCondition,
    });
  }

  /**
   * 删除分享记录
   * @param id
   * @returns
   */
  delete(id: string) {
    return request({
      method: "post",
      url: "/share/delete",
      data: { id: id },
    });
  }
}

export default new ShareApi();
