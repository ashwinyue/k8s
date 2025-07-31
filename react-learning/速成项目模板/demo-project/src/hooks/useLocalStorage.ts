import { useState, useEffect } from 'react';

/**
 * 自定义Hook：使用localStorage进行数据持久化
 * 这是后端工程师常用的数据存储模式的前端实现
 */
function useLocalStorage<T>(key: string, initialValue: T) {
  // 从localStorage读取初始值
  const [storedValue, setStoredValue] = useState<T>(() => {
    try {
      const item = window.localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      console.error(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  });

  // 返回一个包装过的setState函数，会同时更新localStorage
  const setValue = (value: T | ((val: T) => T)) => {
    try {
      // 允许value是一个函数，类似于useState
      const valueToStore = value instanceof Function ? value(storedValue) : value;
      setStoredValue(valueToStore);
      
      // 保存到localStorage
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(key, JSON.stringify(valueToStore));
      }
    } catch (error) {
      console.error(`Error setting localStorage key "${key}":`, error);
    }
  };

  return [storedValue, setValue] as const;
}

export default useLocalStorage;

/**
 * 使用示例：
 * 
 * const [user, setUser] = useLocalStorage('user', { name: '', email: '' });
 * 
 * // 更新用户信息（会自动保存到localStorage）
 * setUser({ name: '张三', email: 'zhangsan@example.com' });
 * 
 * // 使用函数式更新
 * setUser(prevUser => ({ ...prevUser, name: '李四' }));
 */