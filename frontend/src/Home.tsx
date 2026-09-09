import Login from './components/Login.tsx';
import SearchBar from './components/SearchBar.tsx';

const Home = () => {
	const logged: boolean = false

	if(!logged)
		return <Login/>

	return (
		<div class="flex flex-col min-h-screen items-center justify-center">
			<SearchBar/>
		</div>
	)
}

export default Home