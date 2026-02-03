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
      url: "/file/public/list",
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
      url: "/file/protected/list",
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
   * @param onProgress
   * @param id
   * @returns
   */
  downloadPublic(path: string, name: string, onProgress: Function, id: string) {
    return openRequest<Blob>({
      method: "post",
      url: "/file/public/download",
      responseType: "blob",
      timeout: 0,
      data: { path: path, name: name },
      onDownloadProgress: (progressEvent) => {
        onProgress(progressEvent, name, false, id);
      },
    });
  }

  /**
   * 下载保护文件
   * @param path
   * @param name
   * @param onProgress
   * @param id
   * @returns
   */
  downloadProtected(path: string, name: string, onProgress: Function, id: string) {
    return request<Blob>({
      method: "post",
      url: "/file/protected/download",
      responseType: "blob",
      timeout: 0,
      data: { path: path, name: name },
      onDownloadProgress: (progressEvent) => {
        onProgress(progressEvent, name, false, id);
      },
    });
  }

  /**
   * 下载私有文件
   * @param path
   * @param name
   * @param onProgress
   * @param id
   * @returns
   */
  downloadPrivate(path: string, name: string, onProgress: Function, id: string) {
    return request<Blob>({
      method: "post",
      url: "/file/private/download",
      responseType: "blob",
      timeout: 0,
      data: { path: path, name: name },
      onDownloadProgress: (progressEvent) => {
        onProgress(progressEvent, name, false, id);
      },
    });
  }

  /**
   * 上传公开文件
   * @param path
   * @param file
   * @param onProgress
   * @param id
   * @returns
   */
  uploadPublic(path: string, file: File, onProgress: Function, id: string) {
    const formData = new FormData();
    formData.append("path", path);
    formData.append("file", file);
    return request<null>({
      method: "post",
      url: "/file/public/upload",
      timeout: 0,
      data: formData,
      onUploadProgress(progressEvent) {
        onProgress(progressEvent, file.name, true, id);
      },
    });
  }

  /**
   * 上传保护文件
   * @param path
   * @param file
   * @param onProgress
   * @param id
   * @returns
   */
  uploadProtected(path: string, file: File, onProgress: Function, id: string) {
    const formData = new FormData();
    formData.append("path", path);
    formData.append("file", file);
    return request<null>({
      method: "post",
      url: "/file/protected/upload",
      timeout: 0,
      data: formData,
      onUploadProgress(progressEvent) {
        onProgress(progressEvent, file.name, true, id);
      },
    });
  }

  /**
   * 上传私有文件
   * @param path
   * @param file
   * @param onProgress
   * @param id
   * @returns
   */
  uploadPrivate(path: string, file: File, onProgress: Function, id: string) {
    const formData = new FormData();
    formData.append("path", path);
    formData.append("file", file);
    return request<null>({
      method: "post",
      url: "/file/private/upload",
      timeout: 0,
      data: formData,
      onUploadProgress(progressEvent) {
        onProgress(progressEvent, file.name, true, id);
      },
    });
  }

  /**
   * 创建公开目录
   * @param path
   * @param name
   * @returns
   */
  folderPublic(path: string, name: string) {
    return request<null>({
      method: "post",
      url: "/file/public/folder",
      data: { path: path, name: name },
    });
  }

  /**
   * 创建保护目录
   * @param path
   * @param name
   * @returns
   */
  folderProtected(path: string, name: string) {
    return request<null>({
      method: "post",
      url: "/file/protected/folder",
      data: { path: path, name: name },
    });
  }

  /**
   * 创建私有目录
   * @param path
   * @param name
   * @returns
   */
  folderPrivate(path: string, name: string) {
    return request<null>({
      method: "post",
      url: "/file/private/folder",
      data: { path: path, name: name },
    });
  }

  /**
   * 删除公开文件
   * @param path
   * @param name
   * @returns
   */
  deletePublic(path: string, name: string) {
    return request<null>({
      method: "post",
      url: "/file/public/delete",
      data: { path: path, name: name },
    });
  }

  /**
   * 删除保护文件
   * @param path
   * @param name
   * @returns
   */
  deleteProtected(path: string, name: string) {
    return request<null>({
      method: "post",
      url: "/file/protected/delete",
      data: { path: path, name: name },
    });
  }

  /**
   * 删除私有文件
   * @param path
   * @param name
   * @returns
   */
  deletePrivate(path: string, name: string) {
    return request<null>({
      method: "post",
      url: "/file/private/delete",
      data: { path: path, name: name },
    });
  }
}

export default new FileApi();
