import React, { useState, useEffect } from 'react';
import './styles/ChatList.css'; // Make sure to create and import your CSS file

const ChatList = () => {
  // Dummy data for chat items
  const [chats, setChats] = useState([
    // Populate with data fetched from your backend
  ]);
  const [searchTerm, setSearchTerm] = useState('');

  // Handle search bar changes
  const handleSearchChange = (event) => {
    setSearchTerm(event.target.value);
    // You might want to implement debouncing here
  };

  useEffect(() => {
    // TODO: Fetch the chat list from the backend and update the state
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
            // Filter chats based on search term
            // Check if chat name or last message includes the search term
            // This is a simple client-side filter. For large datasets, consider server-side searching
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
