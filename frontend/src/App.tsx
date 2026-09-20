import { NavLink, Route, Routes } from 'react-router-dom'
import ActivityLog from './pages/ActivityLog'
import Dashboard from './pages/Dashboard'
import Incidents from './pages/Incidents'
import Responders from './pages/Responders'
import Shelters from './pages/Shelters'

const sections = [
  { to: '/', label: 'Dashboard', end: true },
  { to: '/incidents', label: 'Incidents', end: false },
  { to: '/shelters', label: 'Shelters', end: false },
  { to: '/responders', label: 'Responders', end: false },
  { to: '/activity', label: 'Activity log', end: false },
]

export default function App() {
  return (
    <div className="app">
      <header className="app-header">
        <div>
          <p className="app-client">CoastLink Council</p>
          <h1>Emergency Support Coordination</h1>
        </div>
      </header>

      <div className="app-body">
        <nav className="app-nav">
          {sections.map((section) => (
            <NavLink key={section.to} to={section.to} end={section.end}>
              {section.label}
            </NavLink>
          ))}
        </nav>

        <main className="app-main">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/incidents" element={<Incidents />} />
            <Route path="/shelters" element={<Shelters />} />
            <Route path="/responders" element={<Responders />} />
            <Route path="/activity" element={<ActivityLog />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}

function NotFound() {
  return (
    <section className="panel">
      <h2>Page not found</h2>
      <p>Pick a section from the navigation on the left.</p>
    </section>
  )
}
