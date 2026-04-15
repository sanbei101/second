import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { request } from "@/utils/request";
import type { Goods } from "./goods";

export type OrderStatus = "pending" | "confirmed" | "cancelled" | "completed";

export type Order = {
  id: string;
  goodsId: string;
  goods?: Goods;
  buyerId: string;
  sellerId: string;
  status: OrderStatus;
  remark: string;
  createdAt: string;
  updatedAt: string;
};

function normalizeGoods(item: any): Goods {
  let images: string[] = [];
  if (Array.isArray(item.images)) {
    images = item.images;
  } else if (typeof item.images === "string") {
    try {
      images = JSON.parse(item.images);
    } catch {
      images = item.images ? [item.images] : [];
    }
  }
  return {
    ...item,
    id: String(item.id),
    sellerId: String(item.sellerId),
    images,
  };
}

function normalizeOrder(item: any): Order {
  const order: Order = {
    ...item,
    id: String(item.id),
    goodsId: String(item.goodsId),
    buyerId: String(item.buyerId),
    sellerId: String(item.sellerId),
  };
  if (item.goods) {
    order.goods = normalizeGoods(item.goods);
  }
  return order;
}

export const useOrderStore = defineStore("order", () => {
  const orders = ref<Order[]>([]);

  async function fetchList(role: "buyer" | "seller" = "buyer") {
    const data = await request<{ orders: any[] }>({
      url: "/orders",
      method: "GET",
      data: { role },
    });
    orders.value = data.orders.map(normalizeOrder);
    return orders.value;
  }

  async function fetchById(id: string) {
    const existing = orders.value.find((o) => o.id === id);
    if (existing) return existing;
    const data = await request<{ order: any }>({
      url: `/orders/${id}`,
      method: "GET",
    });
    const order = normalizeOrder(data.order);
    orders.value.push(order);
    return order;
  }

  async function create(goodsId: string, _buyerId: string, _sellerId: string, remark: string = "") {
    const data = await request<{ order: any }>({
      url: "/orders",
      method: "POST",
      data: { goodsId: Number(goodsId), remark },
    });
    const order = normalizeOrder(data.order);
    orders.value.unshift(order);
    return order.id;
  }

  async function updateStatus(id: string, status: OrderStatus) {
    await request({
      url: `/orders/${id}/status`,
      method: "PUT",
      data: { status },
    });
    const idx = orders.value.findIndex((o) => o.id === id);
    if (idx === -1) return false;
    orders.value[idx].status = status;
    orders.value[idx].updatedAt = new Date().toISOString();
    return true;
  }

  function getById(id: string) {
    return orders.value.find((o) => o.id === id);
  }

  function getByBuyer(buyerId: string) {
    return orders.value
      .filter((o) => o.buyerId === buyerId)
      .sort((a, b) => +new Date(b.createdAt) - +new Date(a.createdAt));
  }

  function getBySeller(sellerId: string) {
    return orders.value
      .filter((o) => o.sellerId === sellerId)
      .sort((a, b) => +new Date(b.createdAt) - +new Date(a.createdAt));
  }

  const pendingOrders = computed(() => orders.value.filter((o) => o.status === "pending"));

  return {
    orders,
    pendingOrders,
    create,
    updateStatus,
    getById,
    getByBuyer,
    getBySeller,
    fetchList,
    fetchById,
  };
});
