import request from './request'

export function getLogs(page: number = 1, pageSize: number = 10, searchKey: string = ''): Promise<ApiResponse<PageResponse<LogEntry>>> {
  return request.get('/log', {
    params: {
      pageNo: page,
      pageSize,
      searchKey
    }
  })
}

export function getLog(id: number): Promise<ApiResponse<LogEntry>> {
  return request.get(`/log/${id}`)
}

// Downloads an attachment recorded on a log entry. The path is resolved server-side
// from the log record, so only files the caller's own logs registered can be fetched.
export function downloadAttachment(logId: number, file: string): Promise<Blob> {
  return request.get('/download', {
    params: { logId, file },
    responseType: 'blob'
  })
}
