import { Route, Router } from '@solidjs/router'
import './App.css'
import Home from './Home'
import Search from './Search'
import { AuthProvider } from './contexts/AuthContext'
import Movie from './Movie'
import Layout from './components/Layout'
import Collection from './Collection'
import CollectionMovie from './components/CollectionMovie'

function App() {
  return (
    <Router
      root={(props) => (
        <AuthProvider>
          <Layout>
            {props.children}
          </Layout>
        </AuthProvider>
      )}
    >
      <Route path="/" component={Home} />
      <Route path="/search" component={Search} />
      <Route path="/collection" component={Collection} />
      <Route path="/collection/:movie_id" component={CollectionMovie} />
      <Route path="/:movie_id" component={Movie} />
    </Router>
  );
}


export default App
