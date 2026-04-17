import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { request } from "@/utils/request";
import type { User } from "./user";

export type GoodsStatus = "on_sale" | "sold" | "off_shelf";
export type GoodsCondition = "全新" | "99新" | "95新" | "9成新" | "8成新";
export type GoodsCategory = "电子产品" | "书籍教材" | "生活用品" | "服装鞋帽" | "交通工具" | "其他";

export type Goods = {
  id: number;
  title: string;
  description: string;
  price: number;
  originalPrice: number;
  category: string;
  condition: string;
  images: string[];
  sellerId: number;
  seller?: User;
  status: GoodsStatus;
  viewCount: number;
  createdAt: string;
};

export const categories: GoodsCategory[] = [
  "电子产品",
  "书籍教材",
  "生活用品",
  "服装鞋帽",
  "交通工具",
  "其他",
];

export const conditions: GoodsCondition[] = ["全新", "99新", "95新", "9成新", "8成新"];

type ListGoodsParams = {
  keyword?: string;
  category?: string;
  minPrice?: number;
  maxPrice?: number;
  page?: number;
  pageSize?: number;
};

type ListGoodsResponse = {
  goods: Goods[];
  total: number;
  page: number;
  pageSize: number;
};

type GoodsResponse = {
  goods: Goods;
};

type MessageResponse = {
  message: string;
};

export const useGoodsStore = defineStore("goods", () => {
  const goodsList = ref<Goods[]>([]);
  const total = ref(0);
  const page = ref(1);
  const pageSize = ref(10);

  async function fetchList(params?: ListGoodsParams) {
    const query: ListGoodsParams = {
      page: params?.page ?? 1,
      pageSize: params?.pageSize ?? 10,
    };
    if (params?.keyword) query.keyword = params.keyword;
    if (params?.category && params.category !== "全部") query.category = params.category;
    if (params?.minPrice !== undefined) query.minPrice = params.minPrice;
    if (params?.maxPrice !== undefined) query.maxPrice = params.maxPrice;

    const data = await request<ListGoodsResponse>({
      url: "/goods",
      method: "GET",
      data: query as Record<string, string | number | undefined>,
    });
    goodsList.value = data.goods;
    total.value = data.total;
    page.value = data.page;
    pageSize.value = data.pageSize;
    return goodsList.value;
  }

  async function getById(id: number) {
    const existing = goodsList.value.find((g) => g.id === id);
    if (existing) return existing;
    const data = await request<GoodsResponse>({
      url: `/goods/${id}`,
      method: "GET",
    });
    return data.goods;
  }

  async function add(goods: Omit<Goods, "id" | "createdAt" | "viewCount" | "status">) {
    const data = await request<GoodsResponse>({
      url: "/goods",
      method: "POST",
      data: goods,
    });
    goodsList.value.unshift(data.goods);
    return data.goods.id;
  }

  async function update(id: number, data: Partial<Goods>) {
    await request<MessageResponse>({
      url: `/goods/${id}`,
      method: "PUT",
      data,
    });
    const idx = goodsList.value.findIndex((g) => g.id === id);
    if (idx > -1) {
      Object.assign(goodsList.value[idx], data);
    }
    return true;
  }

  async function remove(id: number) {
    await request<MessageResponse>({
      url: `/goods/${id}`,
      method: "DELETE",
    });
    const idx = goodsList.value.findIndex((g) => g.id === id);
    if (idx > -1) goodsList.value.splice(idx, 1);
    return true;
  }

  function view(id: number) {
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

  function getBySeller(sellerId: number) {
    return goodsList.value.filter((g) => g.sellerId === sellerId);
  }

  const hasMore = computed(() => goodsList.value.length < total.value);

  return {
    goodsList,
    total,
    page,
    pageSize,
    hasMore,
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
