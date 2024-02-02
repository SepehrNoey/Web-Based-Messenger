import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Signup from './signup'; 
import LoginForm from './Login'; 
import MainPage from './MainPage'; 
import ProfilePage from './ProfilePage'; 

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/api/register" element={<Signup />} />
        <Route path="/api/login" element={<LoginForm />} />
        <Route path="/main" element={<MainPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/" element={<Navigate replace to="/api/login" />} />
      </Routes>
    </Router>
  );
};

export default App;
