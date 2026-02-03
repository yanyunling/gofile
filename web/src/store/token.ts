import { defineStore } from "pinia";

interface TokenStore {
  accessToken: string;
  refreshToken: string;
  isAdmin: boolean;
  username: string;
  publicAuth: number;
  protectedAuth: number;
  privateAuth: number;
}

export const useTokenStore = defineStore("token", {
  state: (): TokenStore => ({
    accessToken: "",
    refreshToken: "",
    isAdmin: false,
    username: "",
    publicAuth: 0,
    protectedAuth: 0,
    privateAuth: 0,
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
      this.publicAuth = tokenResult.publicAuth;
      this.protectedAuth = tokenResult.protectedAuth;
      this.privateAuth = tokenResult.privateAuth;
      this.username = tokenResult.username;
    },
    /**
     * 清空token
     */
    removeToken() {
      this.accessToken = "";
      this.refreshToken = "";
      this.isAdmin = false;
      this.publicAuth = 0;
      this.protectedAuth = 0;
      this.privateAuth = 0;
    },
  },
  persist: true,
});
