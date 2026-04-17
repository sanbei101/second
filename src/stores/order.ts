import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { request } from "@/utils/request";
import type { Goods } from "./goods";

export type OrderStatus = "pending" | "confirmed" | "cancelled" | "completed";

export type Order = {
  id: number;
  goodsId: number;
  goods?: Goods;
  buyerId: number;
  sellerId: number;
  status: OrderStatus;
  remark: string;
  createdAt: string;
  updatedAt: string;
};

type ListOrdersResponse = {
  orders: Order[];
};

type OrderResponse = {
  order: Order;
};

type MessageResponse = {
  message: string;
};

function parseOrderImages(order: Order): Order {
  if (order.goods) {
    const raw = order.goods.images;
    if (Array.isArray(raw)) {
      order.goods.images = raw;
    } else if (typeof raw === "string") {
      try {
        order.goods.images = JSON.parse(raw) as string[];
      } catch {
        order.goods.images = [];
      }
    }
  }
  return order;
}

export const useOrderStore = defineStore("order", () => {
  const orders = ref<Order[]>([]);

  async function fetchList(role: "buyer" | "seller" = "buyer") {
    const data = await request<ListOrdersResponse>({
      url: "/orders",
      method: "GET",
      data: { role },
    });
    orders.value = data.orders.map(parseOrderImages);
    return orders.value;
  }

  async function fetchById(id: number) {
    const existing = orders.value.find((o) => o.id === id);
    if (existing) return existing;
    const data = await request<OrderResponse>({
      url: `/orders/${id}`,
      method: "GET",
    });
    const order = parseOrderImages(data.order);
    orders.value.push(order);
    return order;
  }

  async function create(goodsId: number, _sellerId: number, remark: string = "") {
    const data = await request<OrderResponse>({
      url: "/orders",
      method: "POST",
      data: { goodsId, remark },
    });
    orders.value.unshift(data.order);
    return data.order.id;
  }

  async function updateStatus(id: number, status: OrderStatus) {
    await request<MessageResponse>({
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

  function getById(id: number) {
    return orders.value.find((o) => o.id === id);
  }

  function getByBuyer(buyerId: number) {
    return orders.value
      .filter((o) => o.buyerId === buyerId)
      .sort((a, b) => +new Date(b.createdAt) - +new Date(a.createdAt));
  }

  function getBySeller(sellerId: number) {
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
