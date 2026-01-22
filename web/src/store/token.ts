import { defineStore } from "pinia";

interface TokenStore {
  accessToken: string;
  refreshToken: string;
  isAdmin: boolean;
  username: string;
  nickName: string;
  hasUpdate: number;
}

export const useTokenStore = defineStore("token", {
  state: (): TokenStore => ({
    accessToken: "",
    refreshToken: "",
    isAdmin: false,
    username: "",
    nickName: "",
    hasUpdate: 0,
  }),
  actions: {
    /**
     * 缓存token
     * @param tokenResult
     */
    setToken(tokenResult: TokenResult) {
      this.accessToken = tokenResult.accessToken;
      this.refreshToken = tokenResult.refreshToken;
      this.isAdmin = tokenResult.isAdmin;
      this.hasUpdate = tokenResult.hasUpdate;
      this.setUserInfo(tokenResult.username, tokenResult.nickName);
    },
    /**
     * 缓存用户信息
     * @param username
     * @param nickName
     */
    setUserInfo(username: string, nickName: string) {
      this.username = username;
      this.nickName = nickName;
    },
    /**
     * 清空token
     */
    removeToken() {
      this.accessToken = "";
      this.refreshToken = "";
      this.isAdmin = false;
      this.hasUpdate = 0;
      this.nickName = "";
    },
  },
  persist: true,
});
