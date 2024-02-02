import React from 'react';
import { useNavigate } from 'react-router-dom';
import ChatList from './ChatList'; 
import PrivateChat from './PrivateChat';
import './styles/MainPage.css'; // Ensure you have the corresponding CSS file

const MainPage = () => {
  const navigate = useNavigate();
  const selectedContact = {
    image: 'path_to_image', // Replace with actual image path
    name: 'Contact Name',
    status: 'Online', // Replace with actual status
  };

  return (
    <div className="main-page">
      <header className="main-header">
        <div className="profile-icon" onClick={() => navigate('/profile')}>
          <span className="material-icons">account_circle</span>
        </div>
      </header>
      {/* Chat list and other main page content go here */}
      <ChatList />
      <div className="main-content">
        <PrivateChat contactInfo={selectedContact} />
      </div>
    </div>
  );
};

export default MainPage;
