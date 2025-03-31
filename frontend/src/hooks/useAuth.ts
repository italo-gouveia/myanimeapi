import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { auth, users } from '@/lib/api';
import type { User, UserCredentials, UserCreateRequest } from '@/types/models';

interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
}

export function useAuth() {
  const [state, setState] = useState<AuthState>({
    user: null,
    token: null,
    isLoading: true,
    error: null,
  });
  const router = useRouter();

  useEffect(() => {
    // Check for existing token on mount
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    if (token && user) {
      setState({
        user: JSON.parse(user),
        token,
        isLoading: false,
        error: null,
      });
    } else {
      setState(prev => ({ ...prev, isLoading: false }));
    }
  }, []);

  const login = async (credentials: UserCredentials) => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }));
      const { data } = await auth.login(credentials);
      
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
      
      setState({
        user: data.user,
        token: data.token,
        isLoading: false,
        error: null,
      });
      
      router.push('/dashboard');
    } catch (error) {
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'An error occurred during login',
      }));
      throw error;
    }
  };

  const register = async (userData: UserCreateRequest) => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }));
      await users.create(userData);
      router.push('/login');
    } catch (error) {
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'An error occurred during registration',
      }));
      throw error;
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setState({
      user: null,
      token: null,
      isLoading: false,
      error: null,
    });
    router.push('/login');
  };

  const updateProfile = async (data: {
    username: string;
    email: string;
    currentPassword?: string;
    newPassword?: string;
  }) => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }));
      const { data: updatedUser } = await users.update(state.user!.id, data);
      
      localStorage.setItem('user', JSON.stringify(updatedUser));
      setState(prev => ({
        ...prev,
        user: updatedUser,
        isLoading: false,
        error: null,
      }));
    } catch (error) {
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'An error occurred while updating profile',
      }));
      throw error;
    }
  };

  return {
    ...state,
    isAuthenticated: !!state.token,
    login,
    register,
    logout,
    updateProfile,
  };
} 