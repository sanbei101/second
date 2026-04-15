<script setup lang="ts">
import { ref, computed } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { useGoodsStore } from "@/stores/goods";
import { useUserStore } from "@/stores/user";
import { useOrderStore } from "@/stores/order";

const goodsStore = useGoodsStore();
const userStore = useUserStore();
const orderStore = useOrderStore();

const goodsId = ref("");
const remark = ref("");
const showBuyPopup = ref(false);

const goods = computed(() => goodsStore.getById(goodsId.value));
const isOwner = computed(() => goods.value?.sellerId === userStore.currentUser?.id);
const canBuy = computed(
  () => userStore.isLogin && !isOwner.value && goods.value?.status === "on_sale",
);

onLoad((query) => {
  goodsId.value = query?.id || "";
  if (goodsId.value) goodsStore.view(goodsId.value);
});

function openBuy() {
  if (!userStore.isLogin) {
    uni.showToast({ title: "请先登录", icon: "none" });
    return;
  }
  showBuyPopup.value = true;
}

function confirmBuy() {
  if (!goods.value) return;
  orderStore.create(goods.value.id, userStore.currentUser!.id, goods.value.sellerId, remark.value);
  showBuyPopup.value = false;
  remark.value = "";
  uni.showToast({ title: "购买请求已发送", icon: "success" });
}

function goChat() {
  uni.showToast({ title: "聊天功能开发中", icon: "none" });
}
function goBack() {
  uni.navigateBack();
}
</script>

<template>
  <view v-if="goods">
    <wd-navbar
      :title="goods.title"
      safe-area-inset-top
      fixed
      placeholder
      left-arrow
      @click-left="goBack"
    />

    <wd-swiper
      :list="goods.images.length ? goods.images : ['https://img.yzcdn.cn/vant/defaultpic.png']"
      autoplay
      height="240px"
      style="margin-top: 44px"
    />

    <view style="padding: 16px; background: #fff">
      <view style="display: flex; align-items: baseline; gap: 8px">
        <text style="color: #f44; font-size: 24px; font-weight: bold">¥{{ goods.price }}</text>
        <view style="font-size: 13px; color: #999; text-decoration: line-through">
          原价 ¥{{ goods.originalPrice }}
        </view>
      </view>
      <view style="font-size: 18px; font-weight: bold; margin-top: 8px">{{ goods.title }}</view>
      <view style="display: flex; gap: 8px; margin-top: 8px">
        <wd-tag type="primary">{{ goods.category }}</wd-tag>
        <wd-tag type="success">{{ goods.condition }}</wd-tag>
      </view>
      <view style="font-size: 14px; color: #666; margin-top: 12px; line-height: 1.6">
        {{
          goods.description
        }}
      </view>
      <view style="font-size: 12px; color: #999; margin-top: 12px">
        浏览 {{ goods.viewCount }} · 发布于
        {{ new Date(goods.createdAt).toLocaleDateString() }}
      </view>
    </view>

    <wd-divider />

    <view style="padding: 16px; background: #fff">
      <wd-text text="卖家信息" style="font-size: 15px; font-weight: bold; margin-bottom: 12px" />
      <view style="display: flex; align-items: center; gap: 12px">
        <wd-avatar src="https://img.yzcdn.cn/vant/cat.jpeg" />
        <view>
          <view style="font-size: 15px">校园卖家</view>
          <view style="font-size: 12px; color: #999">信用良好</view>
        </view>
      </view>
    </view>

    <view style="height: 60px" />

    <view
      style="
        position: fixed;
        bottom: 0;
        left: 0;
        right: 0;
        background: #fff;
        border-top: 1px solid #eee;
        padding: 8px 16px;
        display: flex;
        gap: 12px;
        padding-bottom: calc(8px + env(safe-area-inset-bottom));
      "
    >
      <wd-button size="large" style="flex: 1" @click="goChat">联系卖家</wd-button>
      <wd-button v-if="canBuy" type="primary" size="large" style="flex: 2" @click="openBuy">
        立即购买
      </wd-button>
      <wd-button v-else-if="isOwner" size="large" style="flex: 2" disabled>
        自己发布的商品
      </wd-button>
      <wd-button v-else size="large" style="flex: 2" disabled>不可购买</wd-button>
    </view>

    <wd-popup
      v-model="showBuyPopup"
      position="bottom"
      :style="{ padding: '16px', paddingBottom: 'calc(16px + env(safe-area-inset-bottom))' }"
    >
      <view style="font-size: 16px; font-weight: bold; margin-bottom: 12px; text-align: center">
        确认购买
      </view>
      <wd-textarea v-model="remark" placeholder="给卖家留言（可选）" />
      <view style="margin-top: 16px; display: flex; gap: 12px">
        <wd-button block @click="showBuyPopup = false">取消</wd-button>
        <wd-button type="primary" block @click="confirmBuy">确认</wd-button>
      </view>
    </wd-popup>
  </view>

  <wd-empty v-else description="商品不存在" />
</template>
