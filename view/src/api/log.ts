import request from './request'

export function getLogs(
  page: number = 1,
  pageSize: number = 10,
  searchKey: string = '',
  adminOnly: boolean = false
): Promise<ApiResponse<PageResponse<LogEntry>>> {
  const params: Record<string, any> = { pageNo: page, pageSize, searchKey }
  // Ignored by the backend for a non-admin, who only ever sees their own entries
  if (adminOnly) params.adminOnly = true
  return request.get('/log', { params })
}

export function getLog(id: number): Promise<ApiResponse<LogEntry>> {
  return request.get(`/log/${id}`)
}

export function deleteLog(id: number): Promise<ApiResponse<any>> {
  return request.delete(`/log/${id}`)
}

// Deletes the caller's own entries older than `days` (0 removes all of them) together
// with their attachments. Admins clean up their own entries here too, like everyone else.
export function cleanupLogs(days: number): Promise<ApiResponse<{ deleted: number }>> {
  return request.post('/log/cleanup', { days })
}

// Downloads an attachment recorded on a log entry. The path is resolved server-side
// from the log record, so only files the caller's own logs registered can be fetched.
export function downloadAttachment(logId: number, file: string): Promise<Blob> {
  return request.get('/download', {
    params: { logId, file },
    responseType: 'blob'
  })
}
