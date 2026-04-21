import { NavLink } from 'react-router-dom'

export default function Navbar() {
    return (
        <header className="site-nav">
            <div className="site-nav-inner">
                <h1 className="site-title">Trophy Collector</h1>
                <nav aria-label="Primary navigation">
                    <ul className="site-nav-links">
                        <li>
                            <NavLink to="/">Home</NavLink>
                        </li>
                        <li>
                            <NavLink to="/heroes">Heroes</NavLink>
                        </li>
                    </ul>
                </nav>
            </div>
        </header>
    )
}