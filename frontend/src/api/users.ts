import { api } from "./client";

export interface Credentials {
  username: string;
  password: string;
}

export interface User {
  uuid: string;
  name: string;
  role: string;
}

export const usersApi = {
  create: (data: Credentials) => api.post<User>("/users/register", data),
  login: (data: Credentials) => api.post<User>("/users/login", data),
  delete: () => api.delete<void>("/users/profile"),
  rename: (data: Omit<Credentials, "password">) => api.patch<void>("/users/profile", data),
};
