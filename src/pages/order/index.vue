<script setup lang="ts">
import { ref, computed } from "vue";
import { useOrderStore } from "@/stores/order";
import { useUserStore } from "@/stores/user";
import { useGoodsStore } from "@/stores/goods";

const orderStore = useOrderStore();
const userStore = useUserStore();
const goodsStore = useGoodsStore();
const activeTab = ref(0);

const tabs = ["我购买的", "我卖出的"];

const list = computed(() => {
  if (!userStore.currentUser) return [];
  const orders =
    activeTab.value === 0
      ? orderStore.getByBuyer(userStore.currentUser.id)
      : orderStore.getBySeller(userStore.currentUser.id);
  return orders.map((o) => ({
    ...o,
    goods: goodsStore.getById(o.goodsId),
  }));
});

const statusText: Record<string, string> = {
  pending: "待确认",
  confirmed: "已确认",
  cancelled: "已取消",
  completed: "已完成",
};

const statusType: Record<string, string> = {
  pending: "warning",
  confirmed: "success",
  cancelled: "default",
  completed: "primary",
};

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/order/detail/index?id=${id}` });
}
</script>

<template>
  <view>
    <wd-navbar title="我的订单" safe-area-inset-top fixed placeholder />

    <wd-tabs v-model="activeTab">
      <wd-tab v-for="(t, i) in tabs" :key="i" :title="t" />
    </wd-tabs>

    <view style="padding: 12px">
      <wd-empty v-if="list.length === 0" description="暂无订单" />

      <wd-card
        v-for="item in list"
        :key="item.id"
        style="margin-bottom: 12px"
        @click="goDetail(item.id)"
      >
        <view
          style="
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 8px;
          "
        >
          <wd-tag :type="statusType[item.status] as any" size="small">
            {{
              statusText[item.status]
            }}
          </wd-tag>
          <view style="font-size: 12px; color: #999">
            {{
              new Date(item.createdAt).toLocaleString()
            }}
          </view>
        </view>

        <view style="display: flex; gap: 12px">
          <wd-img
            :src="item.goods?.images[0] || 'https://img.yzcdn.cn/vant/defaultpic.png'"
            width="70"
            height="70"
            radius="4"
          />
          <view style="flex: 1">
            <view style="font-size: 15px; font-weight: bold">{{ item.goods?.title }}</view>
            <text style="color: #f44; font-size: 15px; margin-top: 4px">
              ¥{{ item.goods?.price || 0 }}
            </text>
          </view>
        </view>
      </wd-card>
    </view>
  </view>
</template>
