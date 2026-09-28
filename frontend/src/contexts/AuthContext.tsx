import { createContext, useContext, createEffect, createSignal, createResource, Switch, Match } from "solid-js";
import { usersApi, type User } from "../api/auth";
import type { Accessor, ParentProps } from "solid-js";
import { setAuthToken } from "../api/token";
import LogIn from "../components/auth/LogIn";
import SignUp from "../components/auth/SignUp";

interface AuthContextValue {
  user: Accessor<User | null>;
  login: (user: User) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue>();

export function AuthProvider(props: ParentProps) {
  const [user, setUser] = createSignal<User | null>(null);
  const [view, setView] = createSignal<"login" | "signup">("login");

  const [refresh] = createResource(async () => {
    try {
      const user = await usersApi.refresh();
      setUser(user);
      return user;
    } catch {
      return null;
    }
  });

  createEffect(() => {
    setAuthToken(user()?.jwt || null);
  });

  const login = (user: User) => setUser(user);
  const logout = () => setUser(null);

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      <Switch>
        <Match when={refresh.loading}>
          <div class="flex"><span>Loading</span></div>
        </Match>
        <Match when={!user()}>
          <Switch>
            <Match when={view() === "login"}>
              <LogIn onGoToSignup={() => setView("signup")} />
            </Match>
            <Match when={view() === "signup"}>
              <SignUp onGoToLogin={() => setView("login")} />
            </Match>
          </Switch>
        </Match>
        <Match when={user()}>{props.children}</Match>
      </Switch>
    </AuthContext.Provider>
  );
}

export const useAuthContext = () => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuthContext must be used within AuthProvider");
  return ctx;
};
