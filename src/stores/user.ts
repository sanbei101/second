import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { request } from "@/utils/request";

export type User = {
  id: number;
  openid: string;
  nickname: string;
  avatar: string;
  phone: string;
  createdAt: string;
};

const TOKEN_KEY = "token";
const CURRENT_USER_KEY = "campus_secondhand_current_user";

export const useUserStore = defineStore("user", () => {
  const currentUser = ref<User | null>(null);
  const isLogin = computed(() => !!currentUser.value);

  function setAuth(token: string, user: User) {
    uni.setStorageSync(TOKEN_KEY, token);
    currentUser.value = user;
    uni.setStorageSync(CURRENT_USER_KEY, JSON.stringify(user));
  }

  async function fetchProfile() {
    const data = await request<{ user: User }>({
      url: "/users/me",
      method: "GET",
    });
    currentUser.value = data.user;
    uni.setStorageSync(CURRENT_USER_KEY, JSON.stringify(data.user));
    return data.user;
  }

  function init() {
    const token = uni.getStorageSync(TOKEN_KEY);
    const cached = uni.getStorageSync(CURRENT_USER_KEY);
    if (cached) {
      currentUser.value = JSON.parse(cached);
    }
    if (token && currentUser.value) {
      fetchProfile().catch(() => logout());
    }
  }

  async function login(phone: string, password: string) {
    const data = await request<{ token: string; user: User }>({
      url: "/auth/login",
      method: "POST",
      data: { phone, password },
    });
    setAuth(data.token, data.user);
    return true;
  }

  async function wxLogin() {
    const openid = `wx_${Date.now()}`;
    const data = await request<{ token: string; user: User }>({
      url: "/auth/wx-login",
      method: "POST",
      data: { openid },
    });
    setAuth(data.token, data.user);
    return true;
  }

  async function register(phone: string, password: string, nickname: string) {
    const data = await request<{ token: string; user: User }>({
      url: "/auth/register",
      method: "POST",
      data: { phone, password, nickname },
    });
    setAuth(data.token, data.user);
    return true;
  }

  function logout() {
    currentUser.value = null;
    uni.removeStorageSync(TOKEN_KEY);
    uni.removeStorageSync(CURRENT_USER_KEY);
  }

  async function updateProfile(data: Partial<Pick<User, "nickname" | "avatar" | "phone">>) {
    const res = await request<{ user: User }>({
      url: "/users/me",
      method: "PUT",
      data,
    });
    currentUser.value = res.user;
    uni.setStorageSync(CURRENT_USER_KEY, JSON.stringify(res.user));
    return true;
  }

  init();

  return {
    currentUser,
    isLogin,
    login,
    wxLogin,
    register,
    logout,
    updateProfile,
    fetchProfile,
  };
});
