import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Signup from './signup'; // Adjust the path based on your file structure
import LoginForm from './Login'; // Adjust the path based on your file structure
import MainPage from './MainPage'; // Adjust the path based on your file structure
import ProfilePage from './ProfilePage'; // Make sure to import the ProfilePage component

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/api/register" element={<Signup />} />
        <Route path="/api/login" element={<LoginForm />} />
        <Route path="/main" element={<MainPage />} />
        <Route path="/profile" element={<ProfilePage />} /> // Define the route for /profile
        <Route path="/" element={<Navigate replace to="api/login" />} />
      </Routes>
    </Router>
  );
};

export default App;
