import request from "@/utils/request";
import sha256 from "crypto-js/sha256";
import { clone } from "@/utils";

class UserApi {
  /**
   * 分页查询用户记录
   * @param pageCondition
   * @returns
   */
  page(pageCondition: PageCondition<string>) {
    return request<PageResult<User>>({
      method: "post",
      url: "/user/page",
      data: pageCondition,
    });
  }

  /**
   * 保存用户记录
   * @param user
   * @returns
   */
  save(user: User) {
    let userCondition = clone(user);
    if (userCondition.password) {
      userCondition.password = sha256(userCondition.password).toString();
    }
    return request({
      method: "post",
      url: "/user/save",
      data: userCondition,
    });
  }

  /**
   * 删除用户记录
   * @param id
   * @returns
   */
  delete(id: string) {
    return request({
      method: "post",
      url: "/user/delete",
      data: { id: id },
    });
  }

  /**
   * 重置用户密码
   * @param user
   * @returns
   */
  resetPassword(user: User) {
    let userCondition = clone(user);
    userCondition.password = sha256(userCondition.password).toString();
    return request({
      method: "post",
      url: "/user/reset-password",
      data: userCondition,
    });
  }

  /**
   * 更新用户密码
   * @param password
   * @param newPassword
   * @returns
   */
  updatePassword(password: string, newPassword: string) {
    return request({
      method: "post",
      url: "/user/update-password",
      data: { password: sha256(password).toString(), newPassword: sha256(newPassword).toString() },
    });
  }
}

export default new UserApi();
