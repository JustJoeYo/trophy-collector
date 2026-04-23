import { NavLink } from 'react-router-dom'
import './Navbar.css'

export default function Navbar() {
  return (
    <header className="nav">
      <div className="nav-inner">
        <nav aria-label="Primary navigation">
          <ul className="nav-links">
            <li><NavLink to="/">Home</NavLink></li>
            <li><NavLink to="/heroes">Heroes</NavLink></li>
            <li><NavLink to="/items">Items</NavLink></li>
          </ul>
        </nav>
      </div>
    </header>
  )
}
