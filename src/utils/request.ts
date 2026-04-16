const BASE_URL = "https://second-api.sanbei.codes:8849/api";

export function request<T>(options: UniApp.RequestOptions): Promise<T> {
  const token = uni.getStorageSync("token");
  return new Promise((resolve, reject) => {
    uni.request({
      ...options,
      url: BASE_URL + options.url,
      header: {
        "Content-Type": "application/json",
        ...options.header,
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      success: (res) => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data as T);
        } else {
          const msg = (res.data as any)?.error || `请求失败: ${res.statusCode}`;
          uni.showToast({ title: msg, icon: "none" });
          reject(new Error(msg));
        }
      },
      fail: (err) => {
        uni.showToast({ title: "网络错误", icon: "none" });
        reject(err);
      },
    });
  });
}
