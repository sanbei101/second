<script setup lang="ts">
import { ref, computed } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { useGoodsStore, categories } from "@/stores/goods";

const goodsStore = useGoodsStore();

const keyword = ref("");
const activeCategory = ref("全部");
const loading = ref(false);
const refreshing = ref(false);

const filteredGoods = computed(() => goodsStore.goodsList);
const hasMore = computed<boolean>(() => goodsStore.hasMore);

async function loadGoods(isRefresh = false) {
  if (isRefresh) {
    refreshing.value = true;
  } else {
    loading.value = true;
  }
  await goodsStore.fetchList({
    keyword: keyword.value,
    category: activeCategory.value,
    page: isRefresh ? 1 : goodsStore.page,
    pageSize: goodsStore.pageSize,
  });
  if (isRefresh) {
    refreshing.value = false;
  }
  loading.value = false;
}

async function onSearch() {
  await loadGoods(true);
}

function goDetail(id: string) {
  goodsStore.view(id);
  uni.navigateTo({ url: `/pages/goods/detail/index?id=${id}` });
}

function onCategoryChange(val: string | number) {
  activeCategory.value = val as string;
  onSearch();
}

async function onLoadMore() {
  if (!hasMore.value || loading.value) return;
  await goodsStore.fetchList({
    keyword: keyword.value,
    category: activeCategory.value,
    page: goodsStore.page + 1,
    pageSize: goodsStore.pageSize,
  });
}

onShow(() => {
  loadGoods(true);
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
      <view v-if="filteredGoods.length === 0 && !loading" style="margin-top: 40px">
        <wd-empty description="暂无商品" />
      </view>

      <view v-for="item in filteredGoods" :key="item.id" @click="goDetail(item.id)">
        <wd-card :title="item.title" style="margin-bottom: 12px">
          <view style="display: flex; gap: 12px">
            <wd-img :src="item.images[0] || 'https://img.yzcdn.cn/vant/defaultpic.png'" width="80" height="80"
              radius="4" />
            <view style="flex: 1">
              <wd-text :text="item.description" :lines="2" style="color: #666; font-size: 13px" />
              <view style="
                margin-top: 8px;
                display: flex;
                align-items: center;
                justify-content: space-between;
              ">
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

      <view v-if="hasMore" style="text-align: center; padding: 16px">
        <wd-button size="small" :loading="loading" @click="onLoadMore">加载更多</wd-button>
      </view>
      <view v-else-if="filteredGoods.length > 0"
        style="text-align: center; padding: 16px; color: #999; font-size: 12px">
        没有更多了
      </view>

      <wd-toast />
    </view>
  </view>
</template>
