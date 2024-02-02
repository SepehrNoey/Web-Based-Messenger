import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './styles/login.css';

const Login = () => {
  const navigate = useNavigate();
  const [credentials, setCredentials] = useState({ username: '', password: '' });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setCredentials({ ...credentials, [name]: value });
  };

  const handleLogin = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(credentials),
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.message || `HTTP error! status: ${response.status}`);
        }

        localStorage.setItem('jwtToken', data.token);
        setIsLoading(false);
        navigate('/main'); // Assuming '/main' is the route for MainPage
    } catch (error) {
        console.error('Login failed:', error);
        setError(error.message || 'Login failed. Please try again.');
        setIsLoading(false);
    }
};


  return (
    <div className="login-container">
      <form onSubmit={handleLogin}>
        <div>
          <label>Username</label>
          <input type="text" name="username" value={credentials.username} onChange={handleChange} />
        </div>
        <div>
          <label>Password</label>
          <input type="password" name="password" value={credentials.password} onChange={handleChange} />
        </div>
        <button type="submit" disabled={isLoading}>{isLoading ? 'Logging in...' : 'Login'}</button>
        {error && <p className="error-message">{error}</p>}

        <p>Don't have an account? <Link to="/api/register">Sign up</Link></p>
      </form>
    </div>
  );
};

export default Login;