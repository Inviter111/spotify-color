import React, { useState } from "react";
import "./Connect.css";

interface IProps {
    connect(token: string, event: React.FormEvent): void;
}

function Connect({ connect }: IProps) {
    const [token, setToken] = useState("");

    return (
        <>
            <h1 className="header">Please enter access token</h1>
            <form className="form" onSubmit={event => connect(token, event)}>
                <input
                    type="text"
                    className="input"
                    placeholder="Enter a Spotify Token"
                    value={token}
                    onChange={event => setToken(event.target.value)}
                />
                <button className="button">✓</button>
            </form>
        </>
    );
}

export default Connect;
