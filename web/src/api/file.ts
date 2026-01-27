import request from "@/utils/request";
import { openRequest } from "@/utils/request";

class FileApi {
  /**
   * 查询公开文件列表
   * @param path
   * @returns
   */
  listPublic(path: string) {
    return openRequest<FileInfo[]>({
      method: "post",
      url: "/file/list",
      data: { path: path },
    });
  }

  /**
   * 查询保护文件列表
   * @param path
   * @returns
   */
  listProtected(path: string) {
    return request<FileInfo[]>({
      method: "post",
      url: "/file/list",
      data: { path: path },
    });
  }

  /**
   * 查询私有文件列表
   * @param path
   * @returns
   */
  listPrivate(path: string) {
    return request<FileInfo[]>({
      method: "post",
      url: "/file/list/private",
      data: { path: path },
    });
  }
}

export default new FileApi();
