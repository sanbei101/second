import { defineStore } from "pinia";
import { ref, computed } from "vue";
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
}

const STORAGE_KEY = "campus_secondhand_goods";

export const categories = ["电子产品", "书籍教材", "生活用品", "服装鞋帽", "交通工具", "其他"];
export const conditions = ["全新", "99新", "95新", "9成新", "8成新"];

const mockGoods: Goods[] = [
  {
    id: "g_1",
    title: "iPad Air 5 64G",
    description: "去年买的，用得很少，保护的很好，带原装充电器",
    price: 3200,
    originalPrice: 4399,
    category: "电子产品",
    condition: "99新",
    images: ["https://img.yzcdn.cn/vant/ipad.jpeg"],
    sellerId: "user_002",
    status: "on_sale",
    viewCount: 128,
    createdAt: "2024-01-10T10:00:00Z",
  },
  {
    id: "g_2",
    title: "高等数学同济版",
    description: "考研用书，有少量笔记，不影响阅读",
    price: 15,
    originalPrice: 45,
    category: "书籍教材",
    condition: "9成新",
    images: ["https://img.yzcdn.cn/vant/book.jpeg"],
    sellerId: "user_002",
    status: "on_sale",
    viewCount: 56,
    createdAt: "2024-01-12T14:00:00Z",
  },
  {
    id: "g_3",
    title: "捷安特山地车",
    description: "毕业出，骑了两年，变速器正常，送锁",
    price: 450,
    originalPrice: 1200,
    category: "交通工具",
    condition: "8成新",
    images: ["https://img.yzcdn.cn/vant/bike.jpeg"],
    sellerId: "user_002",
    status: "on_sale",
    viewCount: 210,
    createdAt: "2024-01-15T09:00:00Z",
  },
  {
    id: "g_4",
    title: "宿舍小风扇",
    description: "USB供电，三档风速，几乎全新",
    price: 25,
    originalPrice: 59,
    category: "生活用品",
    condition: "99新",
    images: ["https://img.yzcdn.cn/vant/fan.jpeg"],
    sellerId: "user_002",
    status: "on_sale",
    viewCount: 34,
    createdAt: "2024-01-18T16:00:00Z",
  },
];

export const useGoodsStore = defineStore("goods", () => {
  const goodsList = ref<Goods[]>([]);

  function init() {
    const raw = uni.getStorageSync(STORAGE_KEY);
    goodsList.value = raw ? JSON.parse(raw) : [...mockGoods];
  }

  function save() {
    uni.setStorageSync(STORAGE_KEY, JSON.stringify(goodsList.value));
  }

  const onSaleList = computed(() => goodsList.value.filter((g) => g.status === "on_sale"));

  function getById(id: string) {
    return goodsList.value.find((g) => g.id === id);
  }

  function getBySeller(sellerId: string) {
    return goodsList.value.filter((g) => g.sellerId === sellerId);
  }

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

  function add(goods: Omit<Goods, "id" | "createdAt" | "viewCount" | "status">) {
    const newGoods: Goods = {
      ...goods,
      id: `g_${  Date.now()}`,
      status: "on_sale",
      viewCount: 0,
      createdAt: new Date().toISOString(),
    };
    goodsList.value.unshift(newGoods);
    save();
    return newGoods.id;
  }

  function update(id: string, data: Partial<Goods>) {
    const idx = goodsList.value.findIndex((g) => g.id === id);
    if (idx === -1) return false;
    Object.assign(goodsList.value[idx], data);
    save();
    return true;
  }

  function remove(id: string) {
    const idx = goodsList.value.findIndex((g) => g.id === id);
    if (idx === -1) return false;
    goodsList.value.splice(idx, 1);
    save();
    return true;
  }

  function view(id: string) {
    const idx = goodsList.value.findIndex((g) => g.id === id);
    if (idx > -1) {
      goodsList.value[idx].viewCount++;
      save();
    }
  }

  init();

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
  };
});
