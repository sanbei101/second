<script setup lang="ts">
import { computed } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { useGoodsStore } from "@/stores/goods";
import type { Goods } from "@/stores/goods";
import { useUserStore } from "@/stores/user";

const goodsStore = useGoodsStore();
const userStore = useUserStore();
const list = computed(() => goodsStore.getBySeller(String(userStore.currentUser?.id)));

onShow(() => {
  goodsStore.fetchList();
});

function editGoods(item: Goods) {
  uni.navigateTo({ url: `/pages/publish/index?id=${item.id}` });
}

async function deleteGoods(id: string) {
  uni.showModal({
    title: "确认删除",
    content: "删除后不可恢复",
    success: async (res) => {
      if (res.confirm) {
        await goodsStore.remove(id);
        uni.showToast({ title: "已删除", icon: "success" });
      }
    },
  });
}

async function toggleStatus(item: Goods) {
  const next = item.status === "on_sale" ? "off_shelf" : "on_sale";
  await goodsStore.update(item.id, { status: next });
  uni.showToast({ title: next === "on_sale" ? "已上架" : "已下架", icon: "none" });
}
function goBack() {
  uni.navigateBack();
}
</script>

<template>
  <view>
    <wd-navbar
      title="我的发布"
      safe-area-inset-top
      fixed
      placeholder
      left-arrow
      @click-left="goBack"
    />

    <view style="padding: 12px">
      <wd-empty v-if="list.length === 0" description="还没有发布商品" />

      <wd-card v-for="item in list" :key="item.id" :title="item.title" style="margin-bottom: 12px">
        <view style="display: flex; gap: 12px">
          <wd-img
            :src="item.images[0] || 'https://img.yzcdn.cn/vant/defaultpic.png'"
            width="80"
            height="80"
            radius="4"
          />
          <view style="flex: 1">
            <text style="color: #f44; font-size: 16px; font-weight: bold">¥{{ item.price }}</text>
            <view style="margin-top: 4px; font-size: 12px; color: #999">
              {{ item.category }} · {{ item.condition }} · 浏览 {{ item.viewCount }}
            </view>
            <view style="margin-top: 4px">
              <wd-tag
                :type="(item.status === 'on_sale' ? 'success' : 'primary') as any"
                size="small"
              >
                {{ item.status === "on_sale" ? "在售" : "已下架" }}
              </wd-tag>
            </view>
          </view>
        </view>

        <view style="margin-top: 12px; display: flex; gap: 8px; justify-content: flex-end">
          <wd-button size="small" @click="editGoods(item)">编辑</wd-button>
          <wd-button
            size="small"
            :type="(item.status === 'on_sale' ? 'warning' : 'success') as any"
            @click="toggleStatus(item)"
          >
            {{ item.status === "on_sale" ? "下架" : "上架" }}
          </wd-button>
          <wd-button size="small" type="danger" @click="deleteGoods(item.id)">删除</wd-button>
        </view>
      </wd-card>
    </view>
  </view>
</template>
