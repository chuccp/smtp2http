<template>
  <el-dialog
    v-model="open"
    :title="t('user.changePassword')"
    width="450px"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item :label="t('auth.currentPassword')" prop="oldPassword">
        <el-input
          v-model="form.oldPassword"
          type="password"
          show-password
          :placeholder="t('auth.currentPassword')"
        />
      </el-form-item>
      <el-form-item :label="t('user.newPassword')" prop="newPassword">
        <el-input
          v-model="form.newPassword"
          type="password"
          show-password
          :placeholder="t('user.newPassword')"
        />
      </el-form-item>
      <el-form-item :label="t('settings.confirmPassword')" prop="confirmPassword">
        <el-input
          v-model="form.confirmPassword"
          type="password"
          show-password
          :placeholder="t('settings.confirmPassword')"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="open = false">
        {{ t('common.cancel') }}
      </el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitting">
        {{ t('common.save') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { changePassword } from '@/api/auth'
import { useAuthStore } from '@/store/auth'

const { t } = useI18n()
const authStore = useAuthStore()

interface Props {
  open: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateMinLength = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (!value || value.length >= 8) {
    callback()
  } else {
    callback(new Error(t('auth.passwordMinLength')))
  }
}

const validateConfirmPassword = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (value !== form.value.newPassword) {
    callback(new Error(t('settings.passwordsNotMatch')))
  } else {
    callback()
  }
}

const rules = computed<FormRules>(() => ({
  oldPassword: [{ required: true, message: t('auth.currentPassword'), trigger: 'blur' }],
  newPassword: [
    { required: true, message: t('user.newPassword'), trigger: 'blur' },
    { validator: validateMinLength, trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('settings.confirmPassword'), trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}))

watch(() => props.open, (newVal) => {
  if (newVal) {
    form.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
  }
})

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const res = await changePassword(form.value.oldPassword, form.value.newPassword)
      if (res.code === 0 || res.code === 200) {
        emit('update:open', false)
        ElMessage.success(t('auth.passwordChanged'))
        // The change rotates this user's session salt, so the session this page is using
        // is already gone — leave for the login page instead of failing the next request.
        authStore.logoutAction()
      }
    } finally {
      submitting.value = false
    }
  })
}

const open = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val)
})
</script>
