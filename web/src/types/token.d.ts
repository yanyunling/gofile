interface TokenResult {
  accessToken: string;
  refreshToken: string;
  isAdmin: boolean;
  username: string;
  hasUpdate: number;
}

interface CaptchaResult {
  captchaId: string;
  image: string;
  tile: string;
  y: number;
  width: number;
}

interface SignIn {
  username?: string;
  password?: string;
  captchaId: string;
  captchaX: number;
  captchaY: number;
}
