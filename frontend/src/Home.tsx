import Login from './components/auth/Login.tsx';
import Register from './components/auth/Register.tsx';
import SearchBar from './components/SearchBar.tsx';
import { createSignal, Show, Switch, Match } from "solid-js";

const Home = () => {
  const [isConnected, setIsConnected] = createSignal(false);
  const [view, setView] = createSignal("login");

  return (
    <Show when={isConnected()} fallback={
      <Switch>
        <Match when={view() === "login"}>
          <Login
            // onSuccess={() => setIsConnected(true)}
            onGoToRegister={() => setView("register")}
          />
        </Match>
        <Match when={view() === "register"}>
          <Register
            // onSuccess={() => setIsConnected(true)}
            onGoToLogin={() => setView("login")}
          />
        </Match>
      </Switch>
    }>
      <div class="flex flex-col min-h-screen items-center justify-center">
        <SearchBar/>
      </div>
    </Show>
  );
}

export default Home