import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Goods } from './goods'

export type OrderStatus = 'pending' | 'confirmed' | 'cancelled' | 'completed'

export interface Order {
  id: string
  goodsId: string
  goods?: Goods
  buyerId: string
  sellerId: string
  status: OrderStatus
  remark: string
  createdAt: string
  updatedAt: string
}

const STORAGE_KEY = 'campus_secondhand_orders'

export const useOrderStore = defineStore('order', () => {
  const orders = ref<Order[]>([])

  function init() {
    const raw = uni.getStorageSync(STORAGE_KEY)
    orders.value = raw ? JSON.parse(raw) : []
  }

  function save() {
    uni.setStorageSync(STORAGE_KEY, JSON.stringify(orders.value))
  }

  function create(goodsId: string, buyerId: string, sellerId: string, remark: string = '') {
    const order: Order = {
      id: 'o_' + Date.now(),
      goodsId,
      buyerId,
      sellerId,
      status: 'pending',
      remark,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }
    orders.value.unshift(order)
    save()
    return order.id
  }

  function updateStatus(id: string, status: OrderStatus) {
    const idx = orders.value.findIndex(o => o.id === id)
    if (idx === -1) return false
    orders.value[idx].status = status
    orders.value[idx].updatedAt = new Date().toISOString()
    save()
    return true
  }

  function getById(id: string) {
    return orders.value.find(o => o.id === id)
  }

  function getByBuyer(buyerId: string) {
    return orders.value.filter(o => o.buyerId === buyerId).sort((a, b) => +new Date(b.createdAt) - +new Date(a.createdAt))
  }

  function getBySeller(sellerId: string) {
    return orders.value.filter(o => o.sellerId === sellerId).sort((a, b) => +new Date(b.createdAt) - +new Date(a.createdAt))
  }

  const pendingOrders = computed(() => orders.value.filter(o => o.status === 'pending'))

  init()

  return {
    orders,
    pendingOrders,
    create,
    updateStatus,
    getById,
    getByBuyer,
    getBySeller
  }
})
