import React, {useState} from "react";

function PasswordBox({placeholder}) {
    const [password, setPassword] = useState("");

    const onChange = (event) => {
        setPassword(event.target.value);
    };

    return (
        <div>
            <label htmlFor="password"></label>
            <input 
                type="password"
                id="password"
                value={password}
                onChange={onChange}
                placeholder={placeholder}
            />
        </div>
    )
}

export default PasswordBox;