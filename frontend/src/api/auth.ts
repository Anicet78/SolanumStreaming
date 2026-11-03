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

export const usersApi = {
  create: (data: Credentials) => api.post<User>(":8081/register", data),
  login: (data: Credentials) => api.post<User>(":8081/login", data),
  refresh: () => api.get<User>(":8081/refresh"),
  delete: () => api.delete<void>(":8081/profile"),
  update: (data: {new_username: string, new_language: Languages}) => api.patch<void>(":8081/profile", data),
};
