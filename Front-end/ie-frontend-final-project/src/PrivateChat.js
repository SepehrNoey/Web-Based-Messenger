import React, { useState } from 'react';
import './styles/PrivateChat.css'; // Make sure to create and import your CSS file

const PrivateChat = ({ contactInfo }) => {
  const [newMessage, setNewMessage] = useState('');
  const [messages, setMessages] = useState([
    // Placeholder for message data
  ]);

  const handleSendMessage = () => {
    // Logic to send a new message
    // This would involve updating the `messages` state and sending the message to the backend
    console.log(newMessage);
    setNewMessage(''); // Clear input after sending
  };

  return (
    <div className="private-chat-container">
      <div className="contact-info">
        <img src={contactInfo.image} alt={contactInfo.name} className="contact-image" />
        <div className="contact-name">{contactInfo.name}</div>
        <div className="contact-status">{contactInfo.status}</div>
      </div>
      <div className="message-list">
        {messages.map((message) => (
          <div
            key={message.id}
            className={`message-item ${message.sent ? 'sent' : 'received'}`}
          >
            {message.content}
          </div>
        ))}
      </div>
      <div className="message-input-area">
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Type a message..."
          className="message-input"
        />
        <button onClick={handleSendMessage} className="send-button">Send</button>
      </div>
    </div>
  );
};

export default PrivateChat;
