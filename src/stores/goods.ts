import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { request } from "@/utils/request";
import type { User } from "./user";

export type GoodsStatus = "on_sale" | "sold" | "off_shelf";

export type Goods = {
  id: string;
  title: string;
  description: string;
  price: number;
  originalPrice: number;
  category: string;
  condition: string;
  images: string[];
  sellerId: string;
  seller?: User;
  status: GoodsStatus;
  viewCount: number;
  createdAt: string;
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

export const categories = ["电子产品", "书籍教材", "生活用品", "服装鞋帽", "交通工具", "其他"];
export const conditions = ["全新", "99新", "95新", "9成新", "8成新"];

export const useGoodsStore = defineStore("goods", () => {
  const goodsList = ref<Goods[]>([]);

  async function fetchList(params?: {
    keyword?: string;
    category?: string;
    minPrice?: number;
    maxPrice?: number;
  }) {
    const query: Record<string, any> = {};
    if (params?.keyword) query.keyword = params.keyword;
    if (params?.category && params.category !== "全部") query.category = params.category;
    if (params?.minPrice !== undefined) query.minPrice = params.minPrice;
    if (params?.maxPrice !== undefined) query.maxPrice = params.maxPrice;

    const data = await request<{ goods: any[] }>({
      url: "/goods",
      method: "GET",
      data: query,
    });
    goodsList.value = data.goods.map(normalizeGoods);
    return goodsList.value;
  }

  async function getById(id: string) {
    const existing = goodsList.value.find((g) => g.id === id);
    if (existing) return existing;
    const data = await request<{ goods: any }>({
      url: `/goods/${id}`,
      method: "GET",
    });
    return normalizeGoods(data.goods);
  }

  async function add(goods: Omit<Goods, "id" | "createdAt" | "viewCount" | "status">) {
    const data = await request<{ goods: any }>({
      url: "/goods",
      method: "POST",
      data: {
        ...goods,
        sellerId: Number(goods.sellerId),
      },
    });
    const g = normalizeGoods(data.goods);
    goodsList.value.unshift(g);
    return g.id;
  }

  async function update(id: string, data: Partial<Goods>) {
    const payload: any = { ...data };
    if (data.sellerId !== undefined) payload.sellerId = Number(data.sellerId);
    await request({
      url: `/goods/${id}`,
      method: "PUT",
      data: payload,
    });
    const idx = goodsList.value.findIndex((g) => g.id === id);
    if (idx > -1) {
      Object.assign(goodsList.value[idx], data);
    }
    return true;
  }

  async function remove(id: string) {
    await request({
      url: `/goods/${id}`,
      method: "DELETE",
    });
    const idx = goodsList.value.findIndex((g) => g.id === id);
    if (idx > -1) goodsList.value.splice(idx, 1);
    return true;
  }

  function view(id: string) {
    const idx = goodsList.value.findIndex((g) => g.id === id);
    if (idx > -1) {
      goodsList.value[idx].viewCount++;
    }
  }

  const onSaleList = computed(() => goodsList.value.filter((g) => g.status === "on_sale"));

  function filterList(params: {
    keyword?: string;
    category?: string;
    minPrice?: number;
    maxPrice?: number;
  }) {
    return onSaleList.value.filter((g) => {
      if (params.keyword && !g.title.includes(params.keyword)) return false;
      if (params.category && params.category !== "全部" && g.category !== params.category)
        return false;
      if (params.minPrice !== undefined && g.price < params.minPrice) return false;
      if (params.maxPrice !== undefined && g.price > params.maxPrice) return false;
      return true;
    });
  }

  function getBySeller(sellerId: string) {
    return goodsList.value.filter((g) => g.sellerId === sellerId);
  }

  return {
    goodsList,
    onSaleList,
    getById,
    getBySeller,
    filterList,
    add,
    update,
    remove,
    view,
    fetchList,
  };
});
