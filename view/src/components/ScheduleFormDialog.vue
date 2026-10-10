<template>
  <el-dialog
    v-model="open"
    :title="editing ? t('schedule.editSchedule') : t('schedule.addSchedule')"
    width="600px"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item :label="t('schedule.taskName')" prop="name">
        <el-input v-model="form.name" :placeholder="t('schedule.taskName')" />
      </el-form-item>
      <el-form-item :label="t('schedule.associatedToken')" prop="tokenId">
        <TokenSelector v-model="form.tokenId" @update:model-value="formRef?.validateField('tokenId')" />
      </el-form-item>
      <el-form-item :label="t('schedule.cronExpression')" prop="cron">
        <el-input v-model="form.cron" placeholder="* * * * * *" />
        <div class="help-text">
          {{ t('schedule.cronHelp') }}<br>
          {{ t('schedule.exampleCron') }}
        </div>
      </el-form-item>
      <el-form-item :label="t('schedule.requestUrl')" prop="url">
        <el-input v-model="form.url" :placeholder="urlPlaceholder" />
      </el-form-item>
      <el-form-item :label="t('schedule.requestMethod')" prop="method">
        <el-radio-group v-model="form.method">
          <el-radio value="GET">GET</el-radio>
          <el-radio value="POST">POST</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="t('schedule.requestHeaders')">
        <div class="header-list">
          <div v-for="(header, index) in headerList" :key="index" class="header-row">
            <el-input v-model="header.key" placeholder="Key" class="header-key" />
            <el-input v-model="header.value" placeholder="Value" class="header-value" />
            <el-button type="danger" :icon="Delete" circle @click="removeHeader(index)" />
          </div>
          <el-button type="primary" plain :icon="Plus" @click="addHeader">
            {{ t('common.add') }}
          </el-button>
        </div>
      </el-form-item>
      <el-form-item v-if="form.method === 'POST'" :label="t('schedule.requestBody')" prop="body">
        <el-input
          v-model="form.body"
          type="textarea"
          :rows="4"
          placeholder="{}"
        />
      </el-form-item>
      <el-form-item :label="t('schedule.useTemplate')" prop="useTemplate">
        <el-switch v-model="form.useTemplate" />
      </el-form-item>
      <el-form-item v-if="form.useTemplate" :label="t('schedule.template')" prop="template">
        <el-input
          v-model="form.template"
          type="textarea"
          :rows="6"
          :placeholder="t('schedule.templatePlaceholder')"
        />
        <div class="help-text">
          {{ t('schedule.templateHelp') }}
          <a href="https://pkg.go.dev/text/template" target="_blank" rel="noopener">Go template →</a>
        </div>
      </el-form-item>
      <el-form-item :label="t('schedule.onlySendByError')" prop="isSendOnlyByError">
        <el-switch v-model="form.isSendOnlyByError" />
        <div class="help-text">
          {{ t('schedule.onlySendByErrorHelp') }}
        </div>
      </el-form-item>
      <el-form-item :label="t('common.status')" prop="isUse">
        <el-switch v-model="form.isUse" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="open = false">
        {{ t('common.cancel') }}
      </el-button>
      <el-button type="warning" @click="handleTest" :loading="testing" :disabled="!form.cron || !form.url || !form.tokenId">
        {{ t('schedule.testSchedule') }}
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
import { Plus, Delete } from '@element-plus/icons-vue'

const { t } = useI18n()
import { createSchedule, updateSchedule, testSchedule } from '@/api/schedule'
import { ElMessage } from 'element-plus'
import TokenSelector from '@/components/TokenSelector.vue'

interface HeaderItem {
  key: string
  value: string
}

interface Props {
  open: boolean
  editing: ScheduleConfig | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)
const testing = ref(false)
const headerList = ref<HeaderItem[]>([])

const defaultForm: Partial<ScheduleConfig> = {
  name: '',
  tokenId: 0,
  cron: '',
  url: `${window.location.origin}/api/set`,
  method: 'GET',
  headers: [] as any,
  body: '',
  useTemplate: false,
  isUse: true,
  isSendOnlyByError: false
}

// A schedule is called by the server, without the browser session, so the default target
// has to be a route that answers unauthenticated requests: /api/set reports the instance
// status and is reachable without a cookie. /api/scheduleTestApi echoes the request back
// but requires login, so a schedule pointed at it would only ever mail the 403 body.
const form = ref<Partial<ScheduleConfig>>({ ...defaultForm })
const urlPlaceholder = `${window.location.origin}/api/set`

const rules = computed<FormRules<Partial<ScheduleConfig>>>(() => ({
  name: [{ required: true, message: t('schedule.taskName'), trigger: 'blur' }],
  tokenId: [{ required: true, message: t('schedule.associatedToken'), trigger: 'change' }],
  cron: [{ required: true, message: t('schedule.cronExpression'), trigger: 'blur' }],
  url: [{ required: true, message: t('schedule.requestUrl'), trigger: 'blur' }]
}))

const syncHeadersToForm = () => {
  const filtered = headerList.value.filter(h => h.key || h.value)
  form.value.headers = filtered.map(h => ({ name: h.key, value: h.value }))
}

const addHeader = () => {
  headerList.value.push({ key: '', value: '' })
  syncHeadersToForm()
}

const removeHeader = (index: number) => {
  headerList.value.splice(index, 1)
  syncHeadersToForm()
}

const syncHeaderList = () => {
  const h = form.value.headers
  if (h) {
    if (Array.isArray(h)) {
      headerList.value = h.map(item => ({ key: item.name || '', value: item.value || '' }))
    } else if (typeof h === 'string') {
      // backward compat: handle old string format
      try {
        const parsed = JSON.parse(h)
        headerList.value = parsed.map((item: any) => ({ key: item.key || item.name || '', value: item.value || '' }))
      } catch {
        headerList.value = []
      }
    }
  } else {
    headerList.value = []
  }
}

watch(() => props.editing, () => {
  if (props.editing) {
    form.value = { ...props.editing }
    syncHeaderList()
  } else {
    form.value = { ...defaultForm }
    headerList.value = []
  }
}, { immediate: true })

watch(() => props.open, (newVal) => {
  if (newVal) {
    if (!props.editing) {
      form.value = { ...defaultForm }
    }
  }
})

const handleSubmit = async () => {
  if (!formRef.value) return
  syncHeadersToForm()
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (props.editing) {
          await updateSchedule(form.value as ScheduleConfig)
        } else {
          await createSchedule(form.value)
        }
        ElMessage.success(t('common.success'))
        emit('success')
        emit('update:open', false)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleTest = async () => {
  syncHeadersToForm()
  testing.value = true
  try {
    const headersArr = headerList.value
      .filter(h => h.key || h.value)
      .map(h => ({ name: h.key, value: h.value }))
    const payload: any = {
      name: form.value.name,
      cron: form.value.cron,
      url: form.value.url,
      method: form.value.method,
      headers: headersArr,
      body: form.value.body,
      useTemplate: form.value.useTemplate,
      tokenId: form.value.tokenId
    }
    await testSchedule(payload)
    ElMessage.success(t('schedule.testSuccess'))
  } catch (e: any) {
    // The error envelope carries `msg`, not `message`
    ElMessage.error(e?.response?.data?.msg || e?.message || 'Test failed')
  } finally {
    testing.value = false
  }
}

const open = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val)
})
</script>

<style scoped>
.help-text {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.5;
}
.help-text a {
  color: #409eff;
  text-decoration: none;
}
.help-text a:hover {
  text-decoration: underline;
}
.header-list {
  width: 100%;
}
.header-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.header-key {
  flex: 1;
}
.header-value {
  flex: 1;
}
</style>
