import { api } from "./client";

export interface Credentials {
  username: string;
  password: string;
}

export interface User {
  uuid: string;
  name: string;
  role: string;
  jwt: string;
}

export const usersApi = {
  create: (data: Credentials) => api.post<User>(":8081/register", data),
  login: (data: Credentials) => api.post<User>(":8081/login", data),
  delete: () => api.delete<void>(":8081/profile"),
  rename: (data: Omit<Credentials, "password">) => api.patch<void>(":8081/profile", data),
};
