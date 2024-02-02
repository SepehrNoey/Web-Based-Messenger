import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './styles/signup.css';


const Signup = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        firstname: '',
        lastname: '',
        phone: '',
        username: '',
        password: '',
        image: null,
        bio: '',
    });
    const [errors, setErrors] = useState({});

    const handleChange = (e) => {
        const { name, value } = e.target;
        if (name === 'image') {
            setFormData({ ...formData, image: e.target.files[0] });
        } else {
            setFormData({ ...formData, [name]: value });
        }
        // Clear errors for a specific field when user starts correcting it
        if (errors[name]) {
            setErrors({ ...errors, [name]: null });
        }
    };

    const validateForm = () => {
        const newErrors = {};
        // Add validation checks as needed for each field
        if (!formData.firstname.trim()) newErrors.firstname = 'First name is required';
        if (!formData.lastname.trim()) newErrors.lastname = 'Last name is required';
        if (!formData.phone.trim()) newErrors.phone = 'Phone number is required';
        if (!formData.username.trim()) newErrors.username = 'Username is required';
        if (!formData.password) newErrors.password = 'Password is required';

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0; // Returns true if no errors
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        // Client-side validation
        if (!validateForm()) return;

        const formDataToSend = new FormData();
        for (const key in formData) {
            formDataToSend.append(key, formData[key]);
        }

        fetch('/api/register', {
            method: 'POST',
            body: formDataToSend,
        })
            .then((response) => response.json().then(data => ({ status: response.status, body: data })))
            .then(({ status, body }) => {
                if (status >= 400) {
                    // Handle server-side validation errors
                    for (const key in body.errors) {
                        setErrors(prev => ({ ...prev, [key]: body.errors[key].message }));
                    }
                    return;
                }
                // On successful signup
                if (body.token) {
                    localStorage.setItem('jwtToken', body.token);
                    navigate('/api/login');
                }
            })
            .catch((error) => {
                console.error('There was a problem with the fetch operation:', error);
                setErrors({ general: 'An unexpected error occurred. Please try again.' });
            });
    };

    return (
        <form onSubmit={handleSubmit}>
            {Object.keys(formData).map((key) => (
                <div key={key}>
                    <label>{key.charAt(0).toUpperCase() + key.slice(1)}:</label>
                    <input
                        type={key === 'password' ? 'password' : (key === 'image' ? 'file' : 'text')}
                        name={key}
                        value={key !== 'image' ? formData[key] : undefined}
                        onChange={handleChange}
                        style={errors[key] ? { borderColor: 'red' } : {}}
                    />
                    {errors[key] && <p style={{ color: 'red' }}>{errors[key]}</p>}
                </div>
            ))}
            {errors.general && <p style={{ color: 'red', textAlign: 'center' }}>{errors.general}</p>}
            <button type="submit">Sign Up</button>
        </form>
    );
};

export default Signup;