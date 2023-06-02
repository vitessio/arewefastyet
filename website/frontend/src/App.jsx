
import { BrowserRouter,Route,Routes } from 'react-router-dom'

import PublicRoute from './pages/PublicRoute';
import './App.css'

function App() {
  

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path='/*' element={<PublicRoute/>}/>
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
