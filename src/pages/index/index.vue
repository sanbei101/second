<script setup lang="ts">
import { ref, computed } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { useGoodsStore, categories } from "@/stores/goods";
import { useUserStore } from "@/stores/user";

const goodsStore = useGoodsStore();
const userStore = useUserStore();

const keyword = ref("");
const activeCategory = ref("全部");
const priceRange = ref<[number, number]>([0, 10000]);
const showFilter = ref(false);
const loading = ref(false);

const filteredGoods = computed(() => {
  return goodsStore.filterList({
    keyword: keyword.value,
    category: activeCategory.value,
    minPrice: priceRange.value[0],
    maxPrice: priceRange.value[1],
  });
});

function onSearch() {
  loading.value = true;
  setTimeout(() => {
    loading.value = false;
  }, 300);
}

function goDetail(id: string) {
  goodsStore.view(id);
  uni.navigateTo({ url: `/pages/goods/detail/index?id=${id}` });
}

function onCategoryChange(val: string | number) {
  activeCategory.value = val as string;
}

onShow(() => {
  loading.value = false;
});
</script>

<template>
  <view>
    <wd-navbar title="校园二手交易" safe-area-inset-top fixed placeholder />

    <wd-search v-model="keyword" placeholder="搜索商品" @search="onSearch" @change="onSearch" />

    <wd-tabs v-model="activeCategory" @change="onCategoryChange">
      <wd-tab title="全部" name="全部" />
      <wd-tab v-for="cat in categories" :key="cat" :title="cat" :name="cat" />
    </wd-tabs>

    <view style="padding: 12px">
      <view v-if="filteredGoods.length === 0" style="margin-top: 40px">
        <wd-empty description="暂无商品" />
      </view>

      <wd-card
        v-for="item in filteredGoods"
        :key="item.id"
        :title="item.title"
        style="margin-bottom: 12px"
        @click="goDetail(item.id)"
      >
        <view style="display: flex; gap: 12px">
          <wd-img
            :src="item.images[0] || 'https://img.yzcdn.cn/vant/defaultpic.png'"
            width="80"
            height="80"
            radius="4"
          />
          <view style="flex: 1">
            <wd-text :text="item.description" :lines="2" style="color: #666; font-size: 13px" />
            <view
              style="
                margin-top: 8px;
                display: flex;
                align-items: center;
                justify-content: space-between;
              "
            >
              <text style="color: #f44; font-size: 16px; font-weight: bold">¥{{ item.price }}</text>
              <wd-tag type="primary" size="small">{{ item.condition }}</wd-tag>
            </view>
            <view style="margin-top: 4px; font-size: 12px; color: #999">
              {{ item.category }} · 浏览 {{ item.viewCount }}
            </view>
          </view>
        </view>
      </wd-card>
    </view>

    <wd-toast />
  </view>
</template>
