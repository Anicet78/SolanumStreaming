import { Route, Router } from '@solidjs/router'
import './App.css'
import Home from './Home'
import Search from './Search'
import { AuthProvider } from './contexts/AuthContext'
import Movie from './Movie'

function App() {
  return (
  <Router root={(props) => <AuthProvider>{props.children}</AuthProvider>}>
    <Route path="/" component={Home} />
    <Route path="/search" component={Search} />
    <Route path="/:title" component={Movie} />
  </Router>
)}

export default App
