import { PhoneInput } from "react-international-phone";
import { useState } from "react";
import "react-international-phone/style.css";

export default function PhoneNumberBox() {
    const [phone, setPhone] = useState("");

    return (
        <div>
            <PhoneInput
                defaultCountry="ir"
                value="phone"
                onChange={(phone) => setPhone(phone)}
            />
        </div>
    );
};