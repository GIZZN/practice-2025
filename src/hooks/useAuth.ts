'use client';

import { useState, useEffect } from 'react';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Обновляем конфигурацию для fetch запросов
const fetchConfig = {
  credentials: 'include' as const,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
};

// Добавляем функцию для получения заголовков с токеном
const getAuthHeaders = () => {
  const token = localStorage.getItem('authToken');
  return {
    ...fetchConfig.headers,
    'Authorization': `Bearer ${token}`
  };
};

export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  notifications: number;
  phone?: string;
  birthDate: string;
  address?: string;
  city?: string;
  country?: string;
  postalCode: string;
  telegram?: string;
  whatsapp?: string;
  preferredContact: string;
  language?: string;
  createdAt: string;
  updatedAt: string;
}

export interface LoginData {
  email: string;
  password: string;
}

export interface RegisterData extends LoginData {
  name: string;
  confirmPassword: string;
}

export interface AuthContextType {
  user: User | null;
  login: (data: LoginData) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

const isBrowser = typeof window !== 'undefined';

const getStorageItem = (key: string) => {
  if (!isBrowser) return null;
  try {
    const item = localStorage.getItem(key);
    return item ? JSON.parse(item) : null;
  } catch (error) {
    console.error('Error accessing localStorage:', error);
    return null;
  }
};

const setStorageItem = (key: string, value: unknown) => {
  if (!isBrowser) return;
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (error) {
    console.error('Error setting localStorage:', error);
  }
};

const removeStorageItem = (key: string) => {
  if (!isBrowser) return;
  try {
    localStorage.removeItem(key);
  } catch (error) {
    console.error('Error removing from localStorage:', error);
  }
};

export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('authToken');
    if (token) {
      fetchUser(token);
    } else {
      setIsLoading(false);
    }
  }, []);

  const fetchUser = async (token: string) => {
    try {
      const response = await fetch(`${API_URL}/profile`, {
        ...fetchConfig,
        headers: {
          ...fetchConfig.headers,
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error('Ошибка при получении данных пользователя');
      }

      const data = await response.json();
      setUser({
        id: data.id,
        name: data.name,
        email: data.email,
        avatar: data.avatar,
        notifications: data.notifications,
        phone: data.phone,
        birthDate: data.birth_date,
        address: data.address,
        city: data.city,
        country: data.country,
        postalCode: data.postal_code,
        telegram: data.telegram,
        whatsapp: data.whatsapp,
        preferredContact: data.preferred_contact,
        language: data.language,
        createdAt: data.created_at,
        updatedAt: data.updated_at
      });
      setIsAuthenticated(true);
    } catch (error) {
      console.error('Error fetching user:', error);
      localStorage.removeItem('authToken');
      setIsAuthenticated(false);
      setUser(null);
    }
  };

  const validateEmail = (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const validatePassword = (password: string): boolean => {
    return password.length >= 8;
  };

  const login = async (data: LoginData): Promise<void> => {
    try {
      setIsLoading(true);
      setError(null);

      if (!validateEmail(data.email)) {
        throw new Error('Некорректный email');
      }

      if (!validatePassword(data.password)) {
        throw new Error('Пароль должен содержать минимум 8 символов');
      }

      const response = await fetch(`${API_URL}/auth/login`, {
        ...fetchConfig,
        method: 'POST',
        body: JSON.stringify(data)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Ошибка при входе');
      }

      const { token, user } = await response.json();
      
      localStorage.setItem('authToken', token);
      setUser(user);
      setIsAuthenticated(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Произошла ошибка при входе');
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (data: RegisterData): Promise<void> => {
    try {
      setIsLoading(true);
      setError(null);

      if (!validateEmail(data.email)) {
        throw new Error('Некорректный email');
      }

      if (!validatePassword(data.password)) {
        throw new Error('Пароль должен содержать минимум 8 символов');
      }

      if (data.password !== data.confirmPassword) {
        throw new Error('Пароли не совпадают');
      }

      if (!data.name.trim()) {
        throw new Error('Имя обязательно для заполнения');
      }

      const response = await fetch(`${API_URL}/auth/register`, {
        ...fetchConfig,
        method: 'POST',
        body: JSON.stringify(data)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Ошибка при регистрации');
      }

      const { token, user } = await response.json();
      
      localStorage.setItem('authToken', token);
      setUser(user);
      setIsAuthenticated(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Произошла ошибка при регистрации');
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem('authToken');
    setUser(null);
    setIsAuthenticated(false);
  };

  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    login,
    register,
    logout
  };
}; 