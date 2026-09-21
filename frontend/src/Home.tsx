import LogIn from './components/auth/LogIn.tsx';
import SignUp from './components/auth/SignUp.tsx';
import SearchBar from './components/SearchBar.tsx';
import { createSignal, Show, Switch, Match } from "solid-js";

const Home = () => {
  const [isConnected, setIsConnected] = createSignal(false);
  const [view, setView] = createSignal("login");

  return (
    <Show when={isConnected()} fallback={
      <Switch>
        <Match when={view() === "login"}>
          <LogIn
            // onSuccess={() => setIsConnected(true)}
            onGoToSignup={() => setView("signup")}
          />
        </Match>
        <Match when={view() === "signup"}>
          <SignUp
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