import React from "react";
import './Config.css'
import { config } from "../config";
import { ISettings, ConfigMessage } from '../types/common';

interface IState {
    settings: ISettings;
    isLoading: boolean
}

export default class Config extends React.Component<any, IState> {
    conn?: WebSocket;

    constructor(props: any) {
        super(props);
        this.state = {
            isLoading: true,
            settings: {
                enableLyrics: false,
            }
        };
    }

    componentDidMount() {
        this.conn = new WebSocket(`ws://${config.host}/ws`);
        this.conn.onopen = () => {
            const eventData: ConfigMessage = {
                type: 'getState',
            };
            this.conn?.send(JSON.stringify(eventData));
        }
        this.conn.onmessage = (event) => {
            console.log(event.data);
            const parsedEvent: ConfigMessage = JSON.parse(event.data);
            if (parsedEvent.type === 'sendState') {
                this.setState({
                    isLoading: false,
                    settings: {
                        enableLyrics: parsedEvent.data?.enableLyrics,
                    },
                } as IState);
            }
        }
    }

    sendMessage({ target: { checked } }: React.ChangeEvent<HTMLInputElement>) {
        console.log('send message')
        this.setState({
            settings: {
                enableLyrics: checked,
            }
        })
        const event: ConfigMessage = {
            type: 'change',
            data: {
                enableLyrics: checked,
            },
        };
        this.conn?.send(JSON.stringify(event));
    }

    render() {
        const { isLoading, settings: { enableLyrics } } = this.state;
        return (
            <>
                {!isLoading && 
                <div className="container">
                    <div className="config-header">
                        <h3 className="font-color">Config</h3>
                    </div>
                    <div className="settings-container">
                        <div className="settings-cell">
                            <span className="label font-color">Enable lyrics</span>
                            <input type="checkbox" name="lyrics" id="lyrics" className="check-lyrics" onChange={this.sendMessage.bind(this)} checked={enableLyrics}/>
                        </div>
                        <hr className="spacing-line"></hr>
                    </div>
                </div>}
            </>
        );
    }
}
