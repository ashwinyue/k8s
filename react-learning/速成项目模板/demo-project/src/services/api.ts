/**
 * API服务层 - 后端工程师熟悉的模式
 * 统一管理所有API请求，类似于后端的Service层
 */

// 基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api';

// 通用请求配置
const defaultHeaders = {
  'Content-Type': 'application/json',
};

// 请求拦截器 - 类似于后端的中间件
class ApiClient {
  private baseURL: string;
  private headers: Record<string, string>;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
    this.headers = { ...defaultHeaders };
  }

  // 设置认证token
  setAuthToken(token: string) {
    this.headers['Authorization'] = `Bearer ${token}`;
  }

  // 移除认证token
  removeAuthToken() {
    delete this.headers['Authorization'];
  }

  // 通用请求方法
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const config: RequestInit = {
      headers: { ...this.headers, ...options.headers },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // GET请求
  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  // POST请求
  async post<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  // PUT请求
  async put<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  // DELETE请求
  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

// 创建API客户端实例
const apiClient = new ApiClient();

// 用户相关API - 类似于后端的UserController
export const userApi = {
  // 获取用户列表
  getUsers: () => apiClient.get<User[]>('/users'),
  
  // 获取单个用户
  getUser: (id: number) => apiClient.get<User>(`/users/${id}`),
  
  // 创建用户
  createUser: (userData: Omit<User, 'id'>) => 
    apiClient.post<User>('/users', userData),
  
  // 更新用户
  updateUser: (id: number, userData: Partial<User>) => 
    apiClient.put<User>(`/users/${id}`, userData),
  
  // 删除用户
  deleteUser: (id: number) => apiClient.delete(`/users/${id}`),
};

// 认证相关API
export const authApi = {
  // 登录
  login: (credentials: { email: string; password: string }) =>
    apiClient.post<{ token: string; user: User }>('/auth/login', credentials),
  
  // 注册
  register: (userData: { name: string; email: string; password: string }) =>
    apiClient.post<{ token: string; user: User }>('/auth/register', userData),
  
  // 登出
  logout: () => apiClient.post('/auth/logout'),
  
  // 刷新token
  refreshToken: () => apiClient.post<{ token: string }>('/auth/refresh'),
};

// 类型定义 - 与后端保持一致
export interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface ApiResponse<T> {
  data: T;
  message: string;
  success: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// 导出API客户端，供其他地方使用
export { apiClient };

/**
 * 使用示例：
 * 
 * // 在组件中使用
 * const fetchUsers = async () => {
 *   try {
 *     const users = await userApi.getUsers();
 *     setUsers(users);
 *   } catch (error) {
 *     console.error('Failed to fetch users:', error);
 *   }
 * };
 * 
 * // 设置认证token
 * apiClient.setAuthToken('your-jwt-token');
 * 
 * // 创建新用户
 * const newUser = await userApi.createUser({
 *   name: '张三',
 *   email: 'zhangsan@example.com',
 *   role: '用户'
 * });
 */