<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useOrderStore, type OrderStatus } from '@/stores/order'
import { useUserStore } from '@/stores/user'
import { useGoodsStore } from '@/stores/goods'

const orderStore = useOrderStore()
const userStore = useUserStore()
const goodsStore = useGoodsStore()

const orderId = ref('')
const order = computed(() => orderStore.getById(orderId.value))
const goods = computed(() => order.value ? goodsStore.getById(order.value.goodsId) : null)
const isSeller = computed(() => order.value?.sellerId === userStore.currentUser?.id)
const isBuyer = computed(() => order.value?.buyerId === userStore.currentUser?.id)

const statusText: Record<string, string> = {
  pending: '待卖家确认',
  confirmed: '交易已确认',
  cancelled: '交易已取消',
  completed: '交易已完成'
}

onLoad((query) => {
  orderId.value = query?.id || ''
})

function changeStatus(status: OrderStatus) {
  orderStore.updateStatus(orderId.value, status)
  if (status === 'confirmed' || status === 'completed') {
    goodsStore.update(goods.value!.id, { status: 'sold' })
  }
  uni.showToast({ title: '操作成功', icon: 'success' })
}
function goBack() {
  uni.navigateBack()
}
</script>

<template>
  <view v-if="order && goods">
    <wd-navbar title="订单详情" safe-area-inset-top fixed placeholder left-arrow @click-left="goBack" />

    <view style="padding: 16px; background: #4D80F0; color: #fff; text-align: center; margin-top: 44px;">
      <view style="font-size: 20px; font-weight: bold;">{{ statusText[order.status] }}</view>
    </view>

    <view style="padding: 16px; background: #fff;">
      <view style="font-size: 15px; font-weight: bold; margin-bottom: 12px;">商品信息</view>
      <view style="display: flex; gap: 12px;">
        <wd-img :src="goods.images[0] || 'https://img.yzcdn.cn/vant/defaultpic.png'" width="80" height="80" radius="4" />
        <view>
          <view style="font-size: 15px;">{{ goods.title }}</view>
          <text style="color: #f44; font-size: 16px; font-weight: bold; margin-top: 4px;">¥{{ goods.price }}</text>
        </view>
      </view>
    </view>

    <wd-divider />

    <view style="padding: 16px; background: #fff;">
      <wd-cell title="订单编号" :value="order.id" />
      <wd-cell title="下单时间" :value="new Date(order.createdAt).toLocaleString()" />
      <wd-cell title="买家留言" :value="order.remark || '无'" />
    </view>

    <view v-if="order.status === 'pending' && isSeller" style="padding: 16px; display: flex; gap: 12px;">
      <wd-button type="danger" block @click="changeStatus('cancelled')">拒绝交易</wd-button>
      <wd-button type="primary" block @click="changeStatus('confirmed')">确认交易</wd-button>
    </view>

    <view v-if="order.status === 'confirmed' && isBuyer" style="padding: 16px;">
      <wd-button type="primary" block @click="changeStatus('completed')">确认收货</wd-button>
    </view>
  </view>

  <wd-empty v-else description="订单不存在" />
</template>
