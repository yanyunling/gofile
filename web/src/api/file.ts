import request from "@/utils/request";
import { openRequest } from "@/utils/request";

class FileApi {
  /**
   * 查询文件列表
   * @param parentDir
   * @param path
   * @returns
   */
  list(parentDir: string, path: string) {
    let listFunc = request;
    if (parentDir === "public") {
      listFunc = openRequest;
    }
    return listFunc<FileInfo[]>({
      method: "post",
      url: `/file/${parentDir}/list`,
      data: { path: path },
    });
  }

  /**
   * 分享文件
   * @param parentDir
   * @param path
   * @param name
   * @param shareHours
   * @returns
   */
  share(parentDir: string, path: string, name: string, shareHours: number) {
    return request<string>({
      method: "post",
      url: `/file/${parentDir}/share`,
      data: { path: path, name: name, shareHours: shareHours },
    });
  }

  /**
   * 下载文件
   * @param parentDir
   * @param path
   * @param name
   * @param onProgress
   * @param id
   * @returns
   */
  download(parentDir: string, path: string, name: string, onProgress: Function, id: string) {
    let listFunc = request;
    if (parentDir === "public") {
      listFunc = openRequest;
    }
    return listFunc<Blob>({
      method: "post",
      url: `/file/${parentDir}/download`,
      responseType: "blob",
      timeout: 0,
      data: { path: path, name: name },
      onDownloadProgress: (progressEvent) => {
        onProgress(progressEvent, name, false, id);
      },
    });
  }

  /**
   * 上传文件
   * @param parentDir
   * @param path
   * @param file
   * @param onProgress
   * @param id
   * @returns
   */
  upload(parentDir: string, path: string, file: File, onProgress: Function, id: string) {
    const formData = new FormData();
    formData.append("path", path);
    formData.append("file", file);
    return request<null>({
      method: "post",
      url: `/file/${parentDir}/upload`,
      timeout: 0,
      data: formData,
      onUploadProgress(progressEvent) {
        onProgress(progressEvent, file.name, true, id);
      },
    });
  }

  /**
   * 创建目录
   * @param parentDir
   * @param path
   * @param name
   * @returns
   */
  folder(parentDir: string, path: string, name: string) {
    return request<null>({
      method: "post",
      url: `/file/${parentDir}/folder`,
      data: { path: path, name: name },
    });
  }

  /**
   * 删除文件
   * @param parentDir
   * @param path
   * @param name
   * @returns
   */
  delete(parentDir: string, path: string, name: string) {
    return request<null>({
      method: "post",
      url: `/file/${parentDir}/delete`,
      data: { path: path, name: name },
    });
  }
}

export default new FileApi();
