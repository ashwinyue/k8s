import { useState } from 'react'
import './App.css'

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
}

const mockUsers: User[] = [
  { id: 1, name: '张三', email: 'zhangsan@example.com', role: '管理员' },
  { id: 2, name: '李四', email: 'lisi@example.com', role: '用户' },
  { id: 3, name: '王五', email: 'wangwu@example.com', role: '编辑' },
];

function App() {
  const [users, setUsers] = useState<User[]>(mockUsers);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [newUserName, setNewUserName] = useState('');
  const [newUserEmail, setNewUserEmail] = useState('');
  const [showForm, setShowForm] = useState(false);

  const addUser = () => {
    if (newUserName && newUserEmail) {
      const newUser: User = {
        id: Date.now(),
        name: newUserName,
        email: newUserEmail,
        role: '用户'
      };
      setUsers([...users, newUser]);
      setNewUserName('');
      setNewUserEmail('');
      setShowForm(false);
    }
  };

  const deleteUser = (id: number) => {
    setUsers(users.filter(user => user.id !== id));
    if (selectedUser?.id === id) {
      setSelectedUser(null);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <h1 className="text-2xl font-bold text-gray-900">React 速成示例</h1>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-500">后端工程师前端速成</span>
              <div className="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
                <span className="text-white text-sm font-medium">U</span>
              </div>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* 用户列表 */}
          <div className="lg:col-span-2">
            <div className="card">
              <div className="flex justify-between items-center mb-6">
                <h2 className="text-xl font-semibold text-gray-900">用户管理</h2>
                <button
                  onClick={() => setShowForm(!showForm)}
                  className="btn btn-primary"
                >
                  {showForm ? '取消' : '添加用户'}
                </button>
              </div>

              {/* 添加用户表单 */}
              {showForm && (
                <div className="mb-6 p-4 bg-gray-50 rounded-lg">
                  <h3 className="text-lg font-medium mb-4">添加新用户</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <input
                      type="text"
                      placeholder="用户名"
                      value={newUserName}
                      onChange={(e) => setNewUserName(e.target.value)}
                      className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <input
                      type="email"
                      placeholder="邮箱"
                      value={newUserEmail}
                      onChange={(e) => setNewUserEmail(e.target.value)}
                      className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                  <div className="mt-4 flex space-x-2">
                    <button onClick={addUser} className="btn btn-primary">
                      确认添加
                    </button>
                    <button
                      onClick={() => setShowForm(false)}
                      className="btn btn-secondary"
                    >
                      取消
                    </button>
                  </div>
                </div>
              )}

              {/* 用户列表 */}
              <div className="space-y-3">
                {users.map((user) => (
                  <div
                    key={user.id}
                    className={`p-4 border rounded-lg cursor-pointer transition-all ${
                      selectedUser?.id === user.id
                        ? 'border-blue-500 bg-blue-50'
                        : 'border-gray-200 hover:border-gray-300 hover:shadow-sm'
                    }`}
                    onClick={() => setSelectedUser(user)}
                  >
                    <div className="flex justify-between items-center">
                      <div>
                        <h3 className="font-medium text-gray-900">{user.name}</h3>
                        <p className="text-sm text-gray-500">{user.email}</p>
                        <span className="inline-block mt-1 px-2 py-1 text-xs bg-gray-100 text-gray-700 rounded">
                          {user.role}
                        </span>
                      </div>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          deleteUser(user.id);
                        }}
                        className="text-red-600 hover:text-red-800 text-sm font-medium"
                      >
                        删除
                      </button>
                    </div>
                  </div>
                ))}
              </div>

              {users.length === 0 && (
                <div className="text-center py-8 text-gray-500">
                  暂无用户数据
                </div>
              )}
            </div>
          </div>

          {/* 用户详情 */}
          <div className="lg:col-span-1">
            <div className="card">
              <h2 className="text-xl font-semibold text-gray-900 mb-6">用户详情</h2>
              {selectedUser ? (
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      用户名
                    </label>
                    <p className="text-gray-900">{selectedUser.name}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      邮箱
                    </label>
                    <p className="text-gray-900">{selectedUser.email}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      角色
                    </label>
                    <p className="text-gray-900">{selectedUser.role}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      用户ID
                    </label>
                    <p className="text-gray-500 text-sm">{selectedUser.id}</p>
                  </div>
                </div>
              ) : (
                <div className="text-center py-8 text-gray-500">
                  请选择一个用户查看详情
                </div>
              )}
            </div>

            {/* 统计信息 */}
            <div className="card mt-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">统计信息</h3>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-600">总用户数</span>
                  <span className="font-medium">{users.length}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">管理员</span>
                  <span className="font-medium">
                    {users.filter(u => u.role === '管理员').length}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">普通用户</span>
                  <span className="font-medium">
                    {users.filter(u => u.role === '用户').length}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 技术栈说明 */}
        <div className="mt-8 card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">🚀 使用的技术栈</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div className="text-center p-3 bg-blue-50 rounded-lg">
              <div className="font-medium text-blue-900">React 18</div>
              <div className="text-blue-600">UI框架</div>
            </div>
            <div className="text-center p-3 bg-green-50 rounded-lg">
              <div className="font-medium text-green-900">TypeScript</div>
              <div className="text-green-600">类型安全</div>
            </div>
            <div className="text-center p-3 bg-purple-50 rounded-lg">
              <div className="font-medium text-purple-900">Vite</div>
              <div className="text-purple-600">构建工具</div>
            </div>
            <div className="text-center p-3 bg-cyan-50 rounded-lg">
              <div className="font-medium text-cyan-900">Tailwind CSS</div>
              <div className="text-cyan-600">样式框架</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App
