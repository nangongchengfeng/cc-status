import axios from 'axios';

export const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
  timeout: 10000,
});

http.interceptors.request.use((config) => {
  const nextConfig = config;
  nextConfig.headers.Authorization = `Bearer ${import.meta.env.VITE_AUTH_TOKEN ?? 'dev-token'}`;
  return nextConfig;
});
