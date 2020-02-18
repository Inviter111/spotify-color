import React from "react";
import { config } from "../config";

export default class Config extends React.Component {
    conn?: WebSocket;

    constructor(props: any) {
        super(props);
    }

    componentDidMount() {
        this.conn = new WebSocket(`ws://${config.host}/ws`);
    }

    sendMessage() {
        this.conn?.send(JSON.stringify({ test: true }));
    }

    render() {
        return (
            <>
                <p>Config</p>
                <button onClick={this.sendMessage.bind(this)}>
                    Send message
                </button>
            </>
        );
    }
}
