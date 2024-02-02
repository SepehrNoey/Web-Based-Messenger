import React, { useState, useEffect } from 'react';
import './styles/ChatList.css';

const ChatList = () => {
  
  const [chats, setChats] = useState([
    
  ]);
  const [searchTerm, setSearchTerm] = useState('');

 
  const handleSearchChange = (event) => {
    setSearchTerm(event.target.value);
    
  };

  useEffect(() => {
    
  }, []);

  return (
    <div className="chat-list-container">
      <div className="search-container">
        <input
          type="text"
          className="search-bar"
          placeholder="Search"
          value={searchTerm}
          onChange={handleSearchChange}
        />
      </div>
      <ul className="chat-items">
        {chats
          .filter((chat) => {
            
            return (
              chat.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
              chat.lastMessage.toLowerCase().includes(searchTerm.toLowerCase())
            );
          })
          .map((chat) => (
            <li key={chat.id} className="chat-item">
              <img src={chat.image} alt={chat.name} className="chat-image" />
              <div className="chat-info">
                <div className="chat-name">{chat.name}</div>
                <div className="chat-status">{chat.status}</div>
                <div className="chat-last-message">{chat.lastMessage}</div>
              </div>
              <div className="chat-meta">
                <span className="chat-time">{chat.lastTime}</span>
                {chat.unreadCount > 0 && (
                  <span className="chat-unread-count">{chat.unreadCount}</span>
                )}
              </div>
            </li>
          ))}
      </ul>
    </div>
  );
};

export default ChatList;
