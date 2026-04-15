<script setup lang="ts">
import { ref, computed } from "vue";
import { useGoodsStore, categories } from "@/stores/goods";

const goodsStore = useGoodsStore();
const activeIndex = ref(0);

const currentCategory = computed(() => categories[activeIndex.value]);
const list = computed(() => goodsStore.filterList({ category: currentCategory.value }));

function goDetail(id: string) {
  goodsStore.view(id);
  uni.navigateTo({ url: `/pages/goods/detail/index?id=${id}` });
}
</script>

<template>
  <view>
    <wd-navbar title="分类" safe-area-inset-top fixed placeholder />

    <wd-sidebar
      v-model="activeIndex"
      style="height: calc(100vh - 44px - env(safe-area-inset-bottom))"
    >
      <wd-sidebar-item v-for="(cat, idx) in categories" :key="cat" :label="cat" :value="idx" />
    </wd-sidebar>

    <view style="padding: 12px">
      <wd-empty v-if="list.length === 0" description="该分类暂无商品" />

      <wd-card
        v-for="item in list"
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
          </view>
        </view>
      </wd-card>
    </view>
  </view>
</template>
