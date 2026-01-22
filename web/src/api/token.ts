import { authRequest } from "@/utils/request";
import sha256 from "crypto-js/sha256";
import { useTokenStore } from "@/store/token";
import { clone } from "@/utils";

class TokenApi {
  /**
   * 登录
   * @param signIn
   * @returns
   */
  signIn(signIn: SignIn) {
    let condition = clone(signIn);
    condition.password = sha256(condition.password).toString();
    return authRequest<TokenResult>({
      method: "post",
      url: "/sign-in",
      data: condition,
    });
  }

  /**
   * 退出登录
   * @returns
   */
  signOut() {
    const tokenStore = useTokenStore();
    return authRequest({
      method: "post",
      url: "/sign-out",
      data: { refreshToken: tokenStore.refreshToken },
    });
  }

  /**
   * 获取验证码
   * @param captchaId
   * @returns
   */
  captcha(captchaId: string) {
    return authRequest<CaptchaResult>({
      method: "post",
      url: "/captcha",
      data: { captchaId: captchaId },
    });
  }

  /**
   * 校验验证码
   * @param signIn
   * @returns
   */
  validateCaptcha(signIn: SignIn) {
    return authRequest({
      method: "post",
      url: "/captcha/validate",
      data: signIn,
    });
  }
}

export default new TokenApi();
