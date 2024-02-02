import React, { useState } from "react";

function TextBox({placeholder}) {
    const [input, setInput] = useState("");

    const onChange = (event) => {
        setInput(event.target.value);
    };

    return (
        <div>
            <label htmlFor="username"></label>
            <input 
                type="text"
                id="username"
                value={input}
                onChange={onChange}
                placeholder={placeholder}
            />
        </div>
    )
}

export default TextBox;