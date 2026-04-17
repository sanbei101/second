<script setup lang="ts">
import { ref, computed } from "vue";
import { useGoodsStore, categories, conditions } from "@/stores/goods";
import { useUserStore } from "@/stores/user";

const goodsStore = useGoodsStore();
const userStore = useUserStore();

const isEdit = ref(false);
const goodsId = ref(0);
const title = ref("");
const description = ref("");
const price = ref("");
const originalPrice = ref("");
const category = ref("");
const condition = ref("");
const categoryPickerVisible = ref(false);
const conditionPickerVisible = ref(false);
const images = ref<string[]>([]);

const categoryColumns = computed(() => categories.map((c) => ({ value: c, label: c })));
const conditionColumns = computed(() => conditions.map((c) => ({ value: c, label: c })));

function onCategoryConfirm({ value }: any) {
  category.value = value[0];
}

function onConditionConfirm({ value }: any) {
  condition.value = value[0];
}

function openCategoryPicker() {
  categoryPickerVisible.value = true;
}

function openConditionPicker() {
  conditionPickerVisible.value = true;
}

function onUpload() {
  uni.chooseImage({
    count: 1,
    success: (res: any) => {
      images.value.push(res.tempFilePaths[0]);
    },
  });
}

function removeImage(idx: number) {
  images.value.splice(idx, 1);
}

async function submit() {
  if (!title.value || !price.value || !category.value) {
    uni.showToast({ title: "请填写完整信息", icon: "none" });
    return;
  }
  if (!userStore.currentUser) {
    uni.showToast({ title: "请先登录", icon: "none" });
    return;
  }
  const data = {
    title: title.value,
    description: description.value,
    price: Number(price.value),
    originalPrice: Number(originalPrice.value) || Number(price.value),
    category: category.value,
    condition: condition.value,
    images: images.value,
    sellerId: userStore.currentUser.id,
  };
  if (isEdit.value) {
    await goodsStore.update(goodsId.value, data);
    uni.showToast({ title: "修改成功", icon: "success" });
  } else {
    await goodsStore.add(data);
    uni.showToast({ title: "发布成功", icon: "success" });
    title.value = "";
    description.value = "";
    price.value = "";
    originalPrice.value = "";
    category.value = "";
    condition.value = "";
    images.value = [];
  }
  setTimeout(() => {
    uni.switchTab({ url: "/pages/index/index" });
  }, 800);
}
function goBack() {
  uni.navigateBack();
}
</script>

<template>
  <view>
    <wd-navbar
      :title="isEdit ? '编辑商品' : '发布商品'"
      safe-area-inset-top
      fixed
      placeholder
      left-arrow
      @click-left="goBack"
    />

    <view style="padding: 16px">
      <wd-cell-group border>
        <wd-input
          v-model="title"
          label="商品名称"
          label-width="90px"
          placeholder="请输入商品名称"
        />
        <wd-cell
          title="分类"
          :value="category"
          placeholder="请选择分类"
          clickable
          is-link
          @click="openCategoryPicker"
        />
        <wd-cell
          title="成色"
          :value="condition"
          placeholder="请选择成色"
          clickable
          is-link
          @click="openConditionPicker"
        />
        <wd-input
          v-model="price"
          label="售价"
          label-width="90px"
          placeholder="请输入售价"
          type="digit"
        />
        <wd-input
          v-model="originalPrice"
          label="原价"
          label-width="90px"
          placeholder="请输入原价"
          type="digit"
        />
      </wd-cell-group>

      <wd-picker
        v-model:visible="categoryPickerVisible"
        :columns="categoryColumns"
        :model-value="category ? [category] : []"
        @confirm="onCategoryConfirm"
      ></wd-picker>
      <wd-picker
        v-model:visible="conditionPickerVisible"
        :columns="conditionColumns"
        :model-value="condition ? [condition] : []"
        @confirm="onConditionConfirm"
      ></wd-picker>

      <wd-textarea v-model="description" placeholder="请描述商品详情..." style="margin-top: 12px" />

      <view style="margin-top: 16px">
        <wd-text text="商品图片" style="font-size: 14px; color: #333; margin-bottom: 8px" />
        <view style="display: flex; flex-wrap: wrap; gap: 8px">
          <view v-for="(img, idx) in images" :key="idx" style="position: relative">
            <wd-img :src="img" width="80" height="80" radius="4" />
            <wd-icon
              name="close"
              style="
                position: absolute;
                top: -4px;
                right: -4px;
                background: #f44;
                color: #fff;
                border-radius: 50%;
                padding: 2px;
              "
              @click="removeImage(idx)"
            />
          </view>
          <view
            style="
              width: 80px;
              height: 80px;
              border: 1px dashed #ccc;
              display: flex;
              align-items: center;
              justify-content: center;
              border-radius: 4px;
            "
            @click="onUpload"
          >
            <wd-icon name="add" size="24px" style="color: #999" />
          </view>
        </view>
      </view>

      <view style="margin-top: 24px">
        <wd-button type="primary" block @click="submit">
          {{ isEdit ? "保存修改" : "立即发布" }}
        </wd-button>
      </view>
    </view>
  </view>
</template>
