import { createContext, useContext } from "solid-js";
import { createStore } from "solid-js/store";
import type { User } from "../api/auth";
import type { ParentProps } from "solid-js";

const AuthContext = createContext();

export function AuthProvider(props: ParentProps) {
  const [auth, setAuth] = createStore<{ token: string; user: User | null }>({
    token: "",
    user: null,
  });

  const login = (token: string, user: User) => {
    setAuth({ token, user });
  };

  const logout = () => {
    setAuth({ token: "", user: null });
  };

  return (
    <AuthContext.Provider value={{ auth, login, logout }}>
      {props.children}
    </AuthContext.Provider>
  );
}

export const useAuthContext = () => useContext(AuthContext);
