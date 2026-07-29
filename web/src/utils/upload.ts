interface HTMLInputEvent extends Event {
  target: HTMLInputElement & EventTarget;
}

/**
 * 文件上传
 */
export class Upload {
  /**
   * 上传文件类型
   */
  static InputAccept = {
    xlsx: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    image: "image/*",
    uploadImage: ".apng,.bmp,.gif,.ico,.jfif,.jpeg,.jpg,.png,.webp",
  };

  /**
   * 上传文件
   * @param multiple 多选
   * @param accept 支持类型
   * @param isDirectory 是否选择文件夹
   * @returns
   */
  static openFiles(multiple?: boolean, accept?: string, isDirectory?: boolean) {
    return new Promise<FileList>((resolve, reject) => {
      const input = document.createElement("input");
      input.type = "file";
      if (accept) {
        input.accept = accept;
      }
      if (multiple) {
        input.multiple = true;
      }
      if (isDirectory) {
        input.setAttribute("webkitdirectory", "");
      }
      input.onchange = (e: Event) => {
        let event = e as HTMLInputEvent;
        const elem = event.target;
        resolve(elem.files ?? new FileList());
      };
      input.oncancel = (e: Event) => {
        reject(e);
      };
      input.click();
    });
  }

  /**
   * 从 File 对象中提取纯目录路径
   * @param file 包含 webkitRelativePath 的 File 对象
   * @returns 目录路径字符串，如果文件在根目录则返回空字符串
   */
  static getDirectoryPath(file: File): string {
    const relativePath = file.webkitRelativePath;
    if (!relativePath) {
      return ""
    }
    const lastSlashIndex = relativePath.lastIndexOf("/");
    if (lastSlashIndex === -1) {
      return "";
    }
    return relativePath.substring(0, lastSlashIndex);
  }

  /**
   * 获取文件名和后缀
   * @param file
   * @returns
   */
  static getFileName(file: File) {
    let splits = file.name.split(".");
    let name = splits.slice(0, splits.length - 1).join(".");
    let ext = "";
    if (splits.length > 1) {
      ext = splits[splits.length - 1].toLowerCase();
    }
    return { name, ext };
  }

  /**
   * 获取文件后缀名
   * @param fileName
   * @returns
   */
  static getFileExtension(fileName: string) {
    let splits = fileName.split(".");
    let ext = "";
    if (splits.length > 1) {
      ext = splits[splits.length - 1].toLowerCase();
    }
    return ext;
  }

  /**
   * 将Blob读取为ArrayBuffer
   * @param blob
   * @returns
   */
  static readBlobToArrayBuffer(blob: Blob) {
    return new Promise<ArrayBuffer>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const data = e.target?.result;
        if (data) {
          resolve(data as ArrayBuffer);
        } else {
          reject();
        }
      };
      reader.onerror = () => {
        reject();
      };
      reader.readAsArrayBuffer(blob);
    });
  }

  /**
   * 将Blob读取为String
   * @param blob
   * @returns
   */
  static readBlobToString(blob: Blob) {
    return new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const text = (e.target?.result as string) || "";
        resolve(text);
      };
      reader.onerror = () => {
        reject();
      };
      reader.readAsText(blob);
    });
  }

  /**
   * 将Blob读取为Base64
   * @param blob
   * @returns
   */
  static readBlobToBase64(blob: Blob) {
    return new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const text = (e.target?.result as string) || "";
        resolve(text);
      };
      reader.onerror = () => {
        reject();
      };
      reader.readAsDataURL(blob);
    });
  }

  /**
   * 格式化文件大小
   * @param size
   */
  static formatFileSize(size: number) {
    if (size < 1024) {
      return size + "B";
    }
    if (size < 1024 * 1024) {
      return (size / 1024).toFixed(2) + "KB";
    }
    if (size < 1024 * 1024 * 1024) {
      return (size / 1024 / 1024).toFixed(2) + "MB";
    }
    return (size / 1024 / 1024 / 1024).toFixed(2) + "GB";
  }
}
