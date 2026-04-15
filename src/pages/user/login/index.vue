<script setup lang="ts">
import { ref } from "vue";
import { useUserStore } from "@/stores/user";

const userStore = useUserStore();
const phone = ref("");
const password = ref("");
const loading = ref(false);

async function onLogin() {
  if (!phone.value || !password.value) {
    uni.showToast({ title: "请输入手机号和密码", icon: "none" });
    return;
  }
  loading.value = true;
  try {
    await userStore.login(phone.value, password.value);
    uni.showToast({ title: "登录成功", icon: "success" });
    setTimeout(() => {
      uni.switchTab({ url: "/pages/my/index" });
    }, 800);
  } catch {
    uni.showToast({ title: "手机号或密码错误", icon: "none" });
  } finally {
    loading.value = false;
  }
}

async function onWxLogin() {
  await userStore.wxLogin();
  uni.showToast({ title: "微信登录成功", icon: "success" });
  setTimeout(() => {
    uni.switchTab({ url: "/pages/my/index" });
  }, 800);
}

function goRegister() {
  uni.navigateTo({ url: "/pages/user/register/index" });
}
function goBack() {
  uni.navigateBack();
}
</script>

<template>
  <view>
    <wd-navbar title="登录" safe-area-inset-top fixed placeholder left-arrow @click-left="goBack" />

    <view style="padding: 24px 16px">
      <view style="text-align: center; margin-bottom: 32px">
        <wd-text text="校园二手交易平台" style="font-size: 22px; font-weight: bold" />
        <wd-text text="安全便捷的校园交易" style="font-size: 14px; color: #999; margin-top: 8px" />
      </view>

      <view>
        <wd-cell-group border>
          <wd-input
            v-model="phone"
            label="手机号"
            label-width="80px"
            placeholder="请输入手机号"
            clearable
          />
          <wd-input
            v-model="password"
            label="密码"
            label-width="80px"
            placeholder="请输入密码"
            show-password
            clearable
          />
        </wd-cell-group>

        <view style="margin-top: 24px">
          <wd-button type="primary" block :loading="loading" @click="onLogin">登录</wd-button>
        </view>

        <view style="margin-top: 16px; text-align: center">
          <wd-button type="success" block @click="onWxLogin">微信一键登录</wd-button>
        </view>

        <view
          style="margin-top: 16px; text-align: center; color: #4d80f0; font-size: 14px"
          @click="goRegister"
        >
          还没有账号？立即注册
        </view>
      </view>
    </view>
  </view>
</template>
