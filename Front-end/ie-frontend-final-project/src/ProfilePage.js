import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './styles/ProfilePage.css';

const ProfilePage = () => {
  const navigate = useNavigate();
  const [userInfo, setUserInfo] = useState({
    firstname: '',
    lastname: '',
    phone: '',
    username: '',
    password: '',
    image: null,
    bio: '',
  });
  const [error, setError] = useState('');

  const handleChange = (e) => {
    const { name, value, files } = e.target;
    if (name === 'image') {
      setUserInfo({ ...userInfo, image: files[0] });
    } else {
      setUserInfo({ ...userInfo, [name]: value });
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log(userInfo);
  };

  const handleDeleteAccount = async () => {
    const confirmDeletion = window.confirm('Are you sure you want to delete your account? This action cannot be undone.');
    if (confirmDeletion) {
      try {
        const response = await fetch('/api/user/delete', {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('jwtToken')}`,
            'Content-Type': 'application/json',
          },
        });
  
        if (!response.ok) {
          throw new Error('Failed to delete account');
        }
  
        localStorage.removeItem('jwtToken');
        navigate('/login');
      } catch (error) {
        console.error('Error deleting account:', error);
        alert('An error occurred while deleting your account. Please try again.');
      }
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('jwtToken');
    navigate('/login');
  };
  

  return (
    <div className="profile-page">
      <form onSubmit={handleSubmit} className="profile-form">
        {Object.keys(userInfo).map((key) => (
          key !== 'image' ? (
            <div className="form-group" key={key}>
              <label>{key.charAt(0).toUpperCase() + key.slice(1)}</label>
              <input 
                type={key === 'password' ? 'password' : 'text'}
                name={key}
                value={userInfo[key]}
                onChange={handleChange}
              />
            </div>
          ) : (
            <div className="form-group" key={key}>
              <label>Profile Picture</label>
              <input 
                type="file"
                name={key}
                onChange={handleChange}
              />
            </div>
          )
        ))}
        <div className="form-actions">
          <button type="submit" id='save'>Save Changes</button>
          <button type="button" id='logout' onClick={handleLogout}>Log Out</button>
          <button type="button" id='delete' onClick={handleDeleteAccount} className="delete-account">Delete Account</button>
        </div>
      </form>
      {error && <p className="error-message">{error}</p>}
    </div>
  );
};

export default ProfilePage;
