import { Route, Router } from '@solidjs/router'
import './App.css'
import Home from './Home'
import Search from './Search'

function App() {
  return (
  <Router>
    <Route path="/" component={Home} />
    <Route path="/search" component={Search} />
  </Router>
)}

export default App
