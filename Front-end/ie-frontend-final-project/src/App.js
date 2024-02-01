import logo from './logo.svg';
import './App.css';
import SignUp from './components/signUp/signUp';
import { BrowserRouter, Route, Router, Routes } from 'react-router-dom';

function App() {
  return (
    <div className='App'>
      <BrowserRouter>
        <Routes>
          <Route path="/api/register" Component={SignUp}></Route>
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
