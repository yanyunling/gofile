import { defineStore } from "pinia";

interface TokenStore {
  accessToken: string;
  refreshToken: string;
  isAdmin: boolean;
  username: string;
  hasUpdate: number;
}

export const useTokenStore = defineStore("token", {
  state: (): TokenStore => ({
    accessToken: "",
    refreshToken: "",
    isAdmin: false,
    username: "",
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
      this.username = tokenResult.username;
    },
    /**
     * 清空token
     */
    removeToken() {
      this.accessToken = "";
      this.refreshToken = "";
      this.isAdmin = false;
      this.hasUpdate = 0;
    },
  },
  persist: true,
});
