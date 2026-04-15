<script setup lang="ts">
import { ref, computed } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { useUserStore, type UserRole } from "@/stores/user";

const userStore = useUserStore();
const user = computed(() => userStore.currentUser);

const nickname = ref("");
const phone = ref("");
const oldPassword = ref("");
const newPassword = ref("");
const confirmNewPassword = ref("");
const role = ref<UserRole>("buyer");

onShow(() => {
  if (user.value) {
    nickname.value = user.value.nickname;
    phone.value = user.value.phone;
    role.value = user.value.role;
  }
});

function saveProfile() {
  userStore.updateProfile({ nickname: nickname.value, phone: phone.value });
  uni.showToast({ title: "资料已保存", icon: "success" });
}

function changePassword() {
  if (!oldPassword.value || !newPassword.value) {
    uni.showToast({ title: "请填写密码", icon: "none" });
    return;
  }
  if (newPassword.value !== confirmNewPassword.value) {
    uni.showToast({ title: "新密码不一致", icon: "none" });
    return;
  }
  const ok = userStore.updatePassword(oldPassword.value, newPassword.value);
  if (ok) {
    oldPassword.value = "";
    newPassword.value = "";
    confirmNewPassword.value = "";
    uni.showToast({ title: "密码修改成功", icon: "success" });
  } else {
    uni.showToast({ title: "原密码错误", icon: "none" });
  }
}
function goBack() {
  uni.navigateBack();
}
</script>

<template>
  <view>
    <wd-navbar
      title="个人资料"
      safe-area-inset-top
      fixed
      placeholder
      left-arrow
      @click-left="goBack"
    />

    <view style="padding: 16px">
      <wd-text text="基本信息" style="font-size: 14px; color: #999; margin-bottom: 8px" />
      <wd-cell-group border>
        <wd-input v-model="nickname" label="昵称" label-width="80px" placeholder="请输入昵称" />
        <wd-input v-model="phone" label="手机号" label-width="80px" placeholder="请输入手机号" />
      </wd-cell-group>
      <view style="margin-top: 16px">
        <wd-button type="primary" block @click="saveProfile">保存资料</wd-button>
      </view>

      <wd-divider />

      <wd-text text="修改密码" style="font-size: 14px; color: #999; margin-bottom: 8px" />
      <wd-cell-group border>
        <wd-input
          v-model="oldPassword"
          label="原密码"
          label-width="80px"
          placeholder="请输入原密码"
          show-password
        />
        <wd-input
          v-model="newPassword"
          label="新密码"
          label-width="80px"
          placeholder="请输入新密码"
          show-password
        />
        <wd-input
          v-model="confirmNewPassword"
          label="确认密码"
          label-width="80px"
          placeholder="请再次输入新密码"
          show-password
        />
      </wd-cell-group>
      <view style="margin-top: 16px">
        <wd-button type="primary" block @click="changePassword">修改密码</wd-button>
      </view>
    </view>
  </view>
</template>
