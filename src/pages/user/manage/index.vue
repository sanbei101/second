<script setup lang="ts">
import { computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useUserStore, type UserRole } from '@/stores/user'

const userStore = useUserStore()
const users = computed(() => userStore.users)

const roleOptions = [
  { label: '买家', value: 'buyer' },
  { label: '卖家', value: 'seller' },
  { label: '管理员', value: 'admin' }
]

function changeRole(id: string, role: UserRole) {
  userStore.updateUserRole(id, role)
  uni.showToast({ title: '角色已更新', icon: 'success' })
}
function goBack() {
  uni.navigateBack()
}
</script>

<template>
  <view>
    <wd-navbar title="用户管理" safe-area-inset-top fixed placeholder left-arrow @click-left="goBack" />

    <view style="padding: 12px;">
      <wd-card v-for="u in users" :key="u.id" :title="u.nickname" style="margin-bottom: 12px;">
        <wd-cell-group>
          <wd-cell title="手机号" :value="u.phone" />
          <wd-cell title="注册时间" :value="new Date(u.createdAt).toLocaleDateString()" />
          <wd-cell title="当前角色">
            <template #value>
              <wd-tag type="primary">{{ u.role === 'admin' ? '管理员' : u.role === 'seller' ? '卖家' : '买家' }}</wd-tag>
            </template>
          </wd-cell>
        </wd-cell-group>

        <view style="margin-top: 12px; display: flex; gap: 8px;">
          <wd-button v-for="opt in roleOptions" :key="opt.value" size="small"
            :type="u.role === opt.value ? 'primary' : 'info'" @click="changeRole(u.id, opt.value as UserRole)">
            设为{{ opt.label }}
          </wd-button>
        </view>
      </wd-card>
    </view>
  </view>
</template>
