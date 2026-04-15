<script setup lang="ts">
function goBack() {
  uni.navigateBack();
}
</script>

<template>
  <view>
    <wd-navbar title="用户管理" safe-area-inset-top fixed placeholder left-arrow @click-left="goBack" />
    <wd-empty description="功能已下线" />
  </view>
</template>
