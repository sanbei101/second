<script setup lang="ts">
import { computed } from "vue";
import { useUserStore } from "@/stores/user";

const userStore = useUserStore();
const isLogin = computed(() => userStore.isLogin);
const user = computed(() => userStore.currentUser);

function goLogin() {
  uni.navigateTo({ url: "/pages/user/login/index" });
}

function goProfile() {
  uni.navigateTo({ url: "/pages/user/profile/index" });
}

function goManageUsers() {
  uni.navigateTo({ url: "/pages/user/manage/index" });
}

function goManageGoods() {
  uni.navigateTo({ url: "/pages/goods/manage/index" });
}

function logout() {
  userStore.logout();
  uni.showToast({ title: "已退出", icon: "none" });
}
</script>

<template>
  <view>
    <wd-navbar title="我的" safe-area-inset-top fixed placeholder />

    <view
      style="
        padding: 16px;
        background: linear-gradient(180deg, #4d80f0 0%, #ffffff 100%);
        margin-bottom: 12px;
      "
    >
      <view v-if="!isLogin" style="text-align: center; padding: 24px 0" @click="goLogin">
        <wd-button type="primary">登录 / 注册</wd-button>
      </view>

      <view v-else style="display: flex; align-items: center; gap: 16px" @click="goProfile">
        <wd-avatar :src="user?.avatar" size="large" />
        <view>
          <view style="font-size: 18px; font-weight: bold; color: #fff">{{ user?.nickname }}</view>
          <view style="font-size: 13px; color: rgba(255, 255, 255, 0.9); margin-top: 4px">
            {{ user?.role === "admin" ? "管理员" : user?.role === "seller" ? "卖家" : "买家" }}
          </view>
        </view>
      </view>
    </view>

    <wd-cell-group>
      <wd-cell title="我发布的商品" icon="goods" is-link @click="goManageGoods" />
      <wd-cell title="修改资料" icon="edit" is-link @click="goProfile" />
      <wd-cell
        v-if="userStore.isAdmin"
        title="用户管理"
        icon="user"
        is-link
        @click="goManageUsers"
      />
      <wd-cell title="退出登录" icon="logout" is-link @click="logout" />
    </wd-cell-group>
  </view>
</template>
