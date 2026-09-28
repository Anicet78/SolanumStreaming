import LogIn from './components/auth/LogIn.tsx';
import SignUp from './components/auth/SignUp.tsx';
import SearchBar from './components/SearchBar.tsx';
import { createSignal, Show, Switch, Match } from "solid-js";

const Home = () => {
  return (
    <div class="flex flex-col min-h-screen items-center justify-center">
      <SearchBar/>
    </div>
  );
}

export default Home