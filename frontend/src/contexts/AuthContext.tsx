import { createContext, useContext, createEffect, createSignal } from "solid-js";
import type { User } from "../api/auth";
import type { ParentProps } from "solid-js";
import { setAuthToken } from "../api/token";

interface AuthContextValue {
  login: (user: User) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue>();

export function AuthProvider(props: ParentProps) {
  const [auth, setAuth] = createSignal<User | null>(null);

  createEffect(() => {
    setAuthToken(auth()?.jwt || null);
  });

  const login = (user: User) => {
    setAuth(user);
  };

  const logout = () => {
    setAuth(null);
  };

  return (
    <AuthContext.Provider value={{ login, logout }}>
      {props.children}
    </AuthContext.Provider>
  );
}

export const useAuthContext = () => useContext(AuthContext);
