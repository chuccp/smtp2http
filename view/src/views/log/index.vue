<template>
  <div class="log-container app-container">
    <div class="search-bar">
      <el-form :inline="true" class="demo-form-inline">
        <el-form-item>
          <el-input
            v-model="searchKey"
            :placeholder="t('log.searchKeyword')"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item v-if="authStore.getIsAdmin">
          <el-checkbox v-model="adminOnly" :label="t('log.adminCreated')" @change="handleSearch" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            {{ t('common.search') }}
          </el-button>
          <el-button @click="handleReset">
            {{ t('common.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-container table-card-responsive">
      <el-table
        :data="tableData"
        border
        striped
        v-loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="token" :label="t('log.token')" width="120">
          <template #default="{ row }">
            {{ row.token.slice(0, 12) }}...
          </template>
        </el-table-column>
        <el-table-column v-if="authStore.getIsAdmin" prop="userName" :label="t('common.creator')" width="120" />
        <el-table-column prop="smtp" :label="t('log.smtpServer')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="mail" :label="t('log.recipients')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="subject" :label="t('log.subject')" min-width="150" />
        <el-table-column prop="status" :label="t('log.status')" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.statusStr)">
              {{ t('log.' + (row.statusStr || 'unknown')) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" :label="t('log.createTime')" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('common.action')" width="100" align="center" class="full-width-on-mobile">
          <template #default="{ row }">
            <el-button size="small" @click.stop="handleDetail(row)">
              {{ t('log.logDetail') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="loadData"
      @current-change="loadData"
    />

    <!-- Log Detail Dialog -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="t('log.logDetail')"
      width="600px"
    >
      <el-descriptions :column="1" border v-if="currentDetail">
        <el-descriptions-item :label="t('log.token')">
          {{ currentDetail.token }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.subject')">
          {{ currentDetail.subject }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.content')">
          <pre>{{ currentDetail.content }}</pre>
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.result')">
          {{ currentDetail.result }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.status')">
          <el-tag :type="statusTagType(currentDetail.statusStr)">
            {{ t('log.' + (currentDetail.statusStr || 'unknown')) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="attachmentsOf(currentDetail).length" :label="t('log.attachments')">
          <div class="attachment-list">
            <el-button
              v-for="file in attachmentsOf(currentDetail)"
              :key="file.name"
              size="small"
              link
              type="primary"
              @click="handleDownload(currentDetail.id, file.name)"
            >
              <el-icon><Download /></el-icon>
              {{ file.name }}
            </el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.createTime')">
          {{ formatTime(currentDetail.createTime) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/time'
import { getLogs, downloadAttachment } from '@/api/log'
import { useAuthStore } from '@/store/auth'

const { t } = useI18n()
const authStore = useAuthStore()

const loading = ref(false)
const tableData = ref<LogEntry[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKey = ref('')
const adminOnly = ref(false)
const detailDialogVisible = ref(false)
const currentDetail = ref<LogEntry | null>(null)

const loadData = async () => {
  loading.value = true
  try {
    const res = await getLogs(page.value, pageSize.value, searchKey.value, adminOnly.value)
    if (res.code === 0 || res.code === 200) {
      tableData.value = res.data.list
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  page.value = 1
  loadData()
}

const handleReset = () => {
  searchKey.value = ''
  adminOnly.value = false
  page.value = 1
  loadData()
}

const handleRowClick = (row: LogEntry) => {
  handleDetail(row)
}

const handleDetail = (row: LogEntry) => {
  currentDetail.value = row
  detailDialogVisible.value = true
}

const statusTagType = (status?: string) => {
  if (status === 'success') return 'success'
  if (status === 'error') return 'danger'
  return 'warning'
}

// The backend stores attachments as a JSON string on the log row
const attachmentsOf = (row: LogEntry | null): Array<{ name: string }> => {
  if (!row?.files) return []
  try {
    const parsed = JSON.parse(row.files)
    return Array.isArray(parsed) ? parsed.filter(f => f?.name) : []
  } catch {
    return []
  }
}

const handleDownload = async (logId: number, fileName: string) => {
  try {
    const blob = await downloadAttachment(logId, fileName)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fileName
    link.click()
    window.URL.revokeObjectURL(url)
  } catch {
    ElMessage.error(t('log.downloadFailed'))
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.log-container {
  // uses global app-container styles
}

pre {
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
  font-family: inherit;
}

.attachment-list {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}
</style>
