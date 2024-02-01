import React, { useState } from "react";
import TextBox from "../textBox/textBox";
import PasswordBox from "../passwordBox/passwordBox";
import UploadPhoto from "../uploadPhoto/uploadPhoto";
import PhoneNumberBox from "../phoneNumberBox/phoneNumberBox";

function SignUp() {
    

    return (
        <div className='SignUp'>
            
            <PhoneNumberBox></PhoneNumberBox>
            <UploadPhoto></UploadPhoto>
            <TextBox placeholder="First Name"></TextBox>
            <TextBox placeholder="Last Name"></TextBox>
            <TextBox placeholder="Username"></TextBox>
            <PasswordBox placeholder="Password"></PasswordBox>
            <TextBox placeholder="Bio"></TextBox>

        </div>
    )
}

export default SignUp;