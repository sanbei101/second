<script setup lang="ts">
import { ref, computed } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { useUserStore } from "@/stores/user";

const userStore = useUserStore();
const user = computed(() => userStore.currentUser);

const nickname = ref("");
const phone = ref("");

onShow(() => {
  if (user.value) {
    nickname.value = user.value.nickname;
    phone.value = user.value.phone;
  }
});

async function saveProfile() {
  await userStore.updateProfile({ nickname: nickname.value, phone: phone.value });
  uni.showToast({ title: "资料已保存", icon: "success" });
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
    </view>
  </view>
</template>
