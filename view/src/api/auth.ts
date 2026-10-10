import request from './request'

export function getSystemInfo(): Promise<ApiResponse<SystemInfo>> {
  return request.get('/set').then(res => {
    return {
      code: res.code,
      msg: res.msg,
      data: {
        initialized: res.data?.hasInit ?? false,
        dbInitialized: res.data?.hasDbInit ?? false,
        hasAdmin: res.data?.hasAdmin ?? false,
        hasLogin: res.data?.hasLogin,
        isDocker: res.data?.isDocker,
        // Present only when the session cookie is valid; the cookie itself is HttpOnly
        username: res.data?.username ?? '',
        isAdmin: res.data?.isAdmin ?? false,
        version: ''
      }
    }
  })
}

export function login(username: string, password: string): Promise<ApiResponse<any>> {
  return request.post('/signIn', {
    username,
    password
  })
}

export function logout(): Promise<ApiResponse<any>> {
  return request.post('/logout')
}

// Self-service password change: the account page's reset is admin-only, so this is the
// only route a normal user has. Rotates the session salt like any other change, so the
// caller has to sign in again afterwards.
export function changePassword(oldPassword: string, newPassword: string): Promise<ApiResponse<any>> {
  return request.post('/changePassword', { oldPassword, newPassword })
}
