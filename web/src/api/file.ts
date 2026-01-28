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
      url: "/file/private/list",
      data: { path: path },
    });
  }

  /**
   * 下载公开文件
   * @param path
   * @param name
   * @returns
   */
  downloadPublic(path: string, name: string) {
    return openRequest<Blob>({
      method: "post",
      url: "/file/download",
      responseType: "blob",
      data: { path: path, name: name },
    });
  }

  /**
   * 下载保护文件
   * @param path
   * @param name
   * @returns
   */
  downloadProtected(path: string, name: string) {
    return request<Blob>({
      method: "post",
      url: "/file/download",
      responseType: "blob",
      data: { path: path, name: name },
    });
  }

  /**
   * 下载私有文件
   * @param path
   * @param name
   * @returns
   */
  downloadPrivate(path: string, name: string) {
    return request<Blob>({
      method: "post",
      url: "/file/private/download",
      responseType: "blob",
      data: { path: path, name: name },
    });
  }
}

export default new FileApi();
