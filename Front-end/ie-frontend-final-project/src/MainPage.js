import React from 'react';
import { useNavigate } from 'react-router-dom';
import ChatList from './ChatList'; 
import PrivateChat from './PrivateChat';
import './styles/MainPage.css';

const MainPage = () => {
  const navigate = useNavigate();
  const selectedContact = {
    image: 'path_to_image',
    name: 'Contact Name',
    status: 'Online',
  };

  return (
    <div className="main-page">
      <header className="main-header">
        <div className="profile-icon" onClick={() => navigate('/profile')}>
          <span className="material-icons">account_circle</span>
        </div>
      </header>
      <ChatList />
      <div className="main-content">
        <PrivateChat contactInfo={selectedContact} />
      </div>
    </div>
  );
};

export default MainPage;
