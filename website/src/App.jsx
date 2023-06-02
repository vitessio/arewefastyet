
import { BrowserRouter,Route,Routes } from 'react-router-dom'

import PublicRoute from './pages/PublicRoute';


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
