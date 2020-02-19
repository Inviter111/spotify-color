import React from "react";
import axios from "axios";

import TrackInfo from "../components/TrackInfo/TrackInfo";
import Cover from "../components/Cover/Cover";
import Connect from "../components/Connect/Connect";

import { config } from "../config";

import { loadScript } from "../utils/loadScript";

import { IWebPlaybackPlayer, IWebPlaybackState } from "../types/spotify";
import { ISettings, ConfigMessage } from '../types/common';

interface IState {
    imageURL: string;
    songName: string;
    artistName: string;
    connected: boolean;
    settings: ISettings;
}

export default class Index extends React.Component<any, IState> {
    private player?: IWebPlaybackPlayer;
    conn?: WebSocket;

    constructor(props: any) {
        super(props);
        this.state = {
            imageURL: "",
            songName: "",
            artistName: "",
            connected: false,
            settings: {
                enableLyrics: false,
            },
        };
    }

    componentDidMount() {
        this.conn = new WebSocket(`ws://${config.host}/ws`);
        this.conn.onclose = evt => {
            console.log("Connection closed");
        };
        this.conn.onmessage = evt => {
            const parsedEvent: ConfigMessage = JSON.parse(evt.data);
            console.log(parsedEvent)
            if (parsedEvent.type === 'getState') {
                const message: ConfigMessage = {
                    type: 'sendState',
                    data: {
                        enableLyrics: this.state.settings.enableLyrics,
                    },
                };
                this.conn?.send(JSON.stringify(message));
            } else {
                this.setState({
                    settings: parsedEvent.data,
                } as IState);
            }
        };
    }

    private handleStateChange = async (state: IWebPlaybackState | null) => {
        try {
            if (state) {
                const imageURL =
                    state.track_window.current_track.album.images[2].url;
                const songName = state.track_window.current_track.name;
                const artistName =
                    state.track_window.current_track.artists[0].name;
                if (
                    artistName !== this.state.artistName ||
                    songName !== this.state.songName
                ) {
                    this.setState({
                        songName,
                        artistName
                    });
                    const smallImageURL =
                        state.track_window.current_track.album.images[0].url;
                    const res = await axios.get(
                        `${config.baseURL}/get-color?url=${smallImageURL}`
                    );
                    this.setState({
                        imageURL
                    });

                    const { hex } = res.data;

                    const div = document.querySelector(".App");
                    if (div) {
                        if (hex) {
                            div.className = "App";
                            // @ts-ignore
                            div.style.backgroundColor = `#${hex}`;
                        } else {
                            div.className = "App no-color";
                        }
                    }
                }
            }
        } catch (err) {
            console.log(err);
        }
    };

    connect(token: string, event: React.FormEvent) {
        event.preventDefault();
        loadScript({
            source: "https://sdk.scdn.co/spotify-player.js"
        });

        // @ts-ignore
        window.onSpotifyWebPlaybackSDKReady = () => {
            // @ts-ignore
            this.player = new window.Spotify.Player({
                getOAuthToken: (cb: (token: string) => void) => {
                    cb(token);
                },
                name: "Web Playback SDK"
            }) as IWebPlaybackPlayer;

            this.player.addListener(
                "player_state_changed",
                this.handleStateChange
            );

            this.player.connect().then(connected => {
                if (connected) {
                    this.setState({ connected: true });
                }
            });
        };
    }

    render() {
        const { imageURL, songName, artistName, connected, settings } = this.state;
        return (
            <>
                {connected ? (
                    <>
                        {imageURL && <Cover url={imageURL} />}
                        <TrackInfo
                            songName={songName}
                            artistName={artistName}
                        />
                    </>
                ) : (
                    <>
                    <p>Lyrics enabled {settings.enableLyrics.toString()}</p>
                    <Connect connect={this.connect.bind(this)} />
                    </>
                )}
            </>
        );
    }
}
