import { api } from "./client";

export interface Credentials {
  username: string;
  password: string;
}

export type Languages = "fr-FR" | "en-US"

export interface User {
  uuid: string;
  username: string;
  role: string;
  lang: Languages;
  jwt: string;
}

const BASE = import.meta.env.VITE_AUTH_PATH;

export const usersApi = {
  create: (data: Credentials) => api.post<User>(`${BASE}/register`, data),
  login: (data: Credentials) => api.post<User>(`${BASE}/login`, data),
  refresh: () => api.get<User>(`${BASE}/refresh`),
  delete: () => api.delete<void>(`${BASE}/profile`),
  update: (data: {new_username: string, new_language: Languages}) => api.patch<void>(`${BASE}/profile`, data),
};
