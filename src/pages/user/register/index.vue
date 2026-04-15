<script setup lang="ts">
import { ref } from "vue";
import { useUserStore, type UserRole } from "@/stores/user";

const userStore = useUserStore();
const phone = ref("");
const password = ref("");
const confirmPassword = ref("");
const nickname = ref("");
const role = ref<UserRole>("buyer");
const loading = ref(false);

const roleOptions = [
  { label: "买家", value: "buyer" },
  { label: "卖家", value: "seller" },
];

function onRegister() {
  if (!phone.value || !password.value || !nickname.value) {
    uni.showToast({ title: "请填写完整信息", icon: "none" });
    return;
  }
  if (password.value !== confirmPassword.value) {
    uni.showToast({ title: "两次密码不一致", icon: "none" });
    return;
  }
  loading.value = true;
  const ok = userStore.register(phone.value, password.value, role.value, nickname.value);
  loading.value = false;
  if (ok) {
    uni.showToast({ title: "注册成功", icon: "success" });
    setTimeout(() => {
      uni.switchTab({ url: "/pages/my/index" });
    }, 800);
  } else {
    uni.showToast({ title: "手机号已注册", icon: "none" });
  }
}
function goBack() {
  uni.navigateBack();
}
</script>

<template>
  <view>
    <wd-navbar title="注册" safe-area-inset-top fixed placeholder left-arrow @click-left="goBack" />

    <view style="padding: 24px 16px">
      <view>
        <wd-cell-group border>
          <wd-input
            v-model="nickname"
            label="昵称"
            label-width="80px"
            placeholder="请输入昵称"
            clearable
          />
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
          <wd-input
            v-model="confirmPassword"
            label="确认密码"
            label-width="80px"
            placeholder="请再次输入密码"
            show-password
            clearable
          />
          <wd-picker
            :value="role"
            label="身份"
            label-width="80px"
            placeholder="请选择身份"
            :columns="roleOptions"
            @confirm="role = $event.value"
          />
        </wd-cell-group>

        <view style="margin-top: 24px">
          <wd-button type="primary" block :loading="loading" @click="onRegister">注册</wd-button>
        </view>
      </view>
    </view>
  </view>
</template>
