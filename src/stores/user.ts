import { defineStore } from "pinia";
import { ref, computed } from "vue";

export type UserRole = "buyer" | "seller" | "admin";

export type User = {
  id: string;
  openid: string;
  nickname: string;
  avatar: string;
  phone: string;
  role: UserRole;
  password: string;
  createdAt: string;
};

const STORAGE_KEY = "campus_secondhand_users";
const CURRENT_USER_KEY = "campus_secondhand_current_user";

const defaultUsers: User[] = [
  {
    id: "admin_001",
    openid: "wx_admin_001",
    nickname: "管理员",
    avatar: "https://img.yzcdn.cn/vant/cat.jpeg",
    phone: "13800000000",
    role: "admin",
    password: "admin123",
    createdAt: "2024-01-01T00:00:00Z",
  },
  {
    id: "user_001",
    openid: "wx_user_001",
    nickname: "测试买家",
    avatar: "https://img.yzcdn.cn/vant/cat.jpeg",
    phone: "13800138000",
    role: "buyer",
    password: "123456",
    createdAt: "2024-01-02T00:00:00Z",
  },
  {
    id: "user_002",
    openid: "wx_user_002",
    nickname: "测试卖家",
    avatar: "https://img.yzcdn.cn/vant/cat.jpeg",
    phone: "13800138111",
    role: "seller",
    password: "123456",
    createdAt: "2024-01-03T00:00:00Z",
  },
];

export const useUserStore = defineStore("user", () => {
  const users = ref<User[]>([]);
  const currentUser = ref<User | null>(null);

  const isLogin = computed(() => !!currentUser.value);
  const isAdmin = computed(() => currentUser.value?.role === "admin");
  const isSeller = computed(() => currentUser.value?.role === "seller");

  function init() {
    const raw = uni.getStorageSync(STORAGE_KEY);
    users.value = raw ? JSON.parse(raw) : [...defaultUsers];
    const cur = uni.getStorageSync(CURRENT_USER_KEY);
    currentUser.value = cur ? JSON.parse(cur) : null;
  }

  function saveUsers() {
    uni.setStorageSync(STORAGE_KEY, JSON.stringify(users.value));
  }

  function saveCurrent() {
    if (currentUser.value) {
      uni.setStorageSync(CURRENT_USER_KEY, JSON.stringify(currentUser.value));
    } else {
      uni.removeStorageSync(CURRENT_USER_KEY);
    }
  }

  function login(phone: string, password: string) {
    const u = users.value.find((item) => item.phone === phone && item.password === password);
    if (!u) return false;
    currentUser.value = u;
    saveCurrent();
    return true;
  }

  function wxLogin(role: UserRole = "buyer") {
    const openid = `wx_${Date.now()}`;
    const exist = users.value.find((u) => u.openid === openid);
    if (exist) {
      currentUser.value = exist;
      saveCurrent();
      return true;
    }
    const newUser: User = {
      id: `u_${Date.now()}`,
      openid,
      nickname: `微信用户${Math.floor(Math.random() * 10000)}`,
      avatar: "https://img.yzcdn.cn/vant/cat.jpeg",
      phone: "",
      role,
      password: "123456",
      createdAt: new Date().toISOString(),
    };
    users.value.push(newUser);
    currentUser.value = newUser;
    saveUsers();
    saveCurrent();
    return true;
  }

  function register(phone: string, password: string, role: UserRole, nickname: string) {
    if (users.value.some((u) => u.phone === phone)) return false;
    const newUser: User = {
      id: `u_${Date.now()}`,
      openid: `wx_${Date.now()}`,
      nickname: nickname || `用户${Math.floor(Math.random() * 10000)}`,
      avatar: "https://img.yzcdn.cn/vant/cat.jpeg",
      phone,
      role,
      password,
      createdAt: new Date().toISOString(),
    };
    users.value.push(newUser);
    currentUser.value = newUser;
    saveUsers();
    saveCurrent();
    return true;
  }

  function logout() {
    currentUser.value = null;
    saveCurrent();
  }

  function updatePassword(oldPwd: string, newPwd: string) {
    if (!currentUser.value) return false;
    if (currentUser.value.password !== oldPwd) return false;
    currentUser.value.password = newPwd;
    const idx = users.value.findIndex((u) => u.id === currentUser.value!.id);
    if (idx > -1) users.value[idx].password = newPwd;
    saveUsers();
    saveCurrent();
    return true;
  }

  function updateProfile(data: Partial<User>) {
    if (!currentUser.value) return false;
    Object.assign(currentUser.value, data);
    const idx = users.value.findIndex((u) => u.id === currentUser.value!.id);
    if (idx > -1) Object.assign(users.value[idx], data);
    saveUsers();
    saveCurrent();
    return true;
  }

  function updateUserRole(id: string, role: UserRole) {
    const idx = users.value.findIndex((u) => u.id === id);
    if (idx === -1) return false;
    users.value[idx].role = role;
    saveUsers();
    return true;
  }

  init();

  return {
    users,
    currentUser,
    isLogin,
    isAdmin,
    isSeller,
    login,
    wxLogin,
    register,
    logout,
    updatePassword,
    updateProfile,
    updateUserRole,
  };
});
