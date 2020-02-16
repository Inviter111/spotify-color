import React from "react";
import axios from 'axios'

import "./App.css";
import { config } from './config'

import { loadScript } from "./utils/loadScript";

import TrackInfo from "./components/TrackInfo/TrackInfo";
import Cover from "./components/Cover/Cover";
import Connect from "./components/Connect/Connect";

import { IWebPlaybackPlayer, IWebPlaybackState } from "./types/spotify";

interface IState {
  imageURL: string;
  songName: string;
  artistName: string;
  connected: boolean;
}

class App extends React.Component<any, IState> {
  private player?: IWebPlaybackPlayer;

  constructor(props: any) {
    super(props);
    this.state = {
      imageURL: "",
      songName: "",
      artistName: "",
      connected: false
    };
  }

  private handleStateChange = async (state: IWebPlaybackState | null) => {
    try {
      if (state) {
        const imageURL = state.track_window.current_track.album.images[2].url;
        const songName = state.track_window.current_track.name;
        const artistName = state.track_window.current_track.artists[0].name;
        if (artistName !== this.state.artistName || songName !== this.state.songName) {
          this.setState({
            songName,
            artistName,
          });
          const smallImageURL = state.track_window.current_track.album.images[0].url;
          const res = await axios.get(`${config.baseURL}/getColor?url=${smallImageURL}`);
          this.setState({
            imageURL,
          });
  
          const { hex } = res.data;
  
          const div = document.querySelector(".App");
          if (div) {
            if (hex) {
              div.className = 'App'
              // @ts-ignore
              div.style.backgroundColor = `#${hex}`
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

      this.player.addListener("player_state_changed", this.handleStateChange);

      this.player.connect().then((connected) => {
        if (connected) {
          this.setState({ connected: true })
        }
      });
    };
  }

  render() {
    const { imageURL, songName, artistName, connected } = this.state;
    return (
      <div className="App">
        {connected ? (
          <>
            {imageURL && <Cover url={imageURL} />}
            <TrackInfo songName={songName} artistName={artistName} />
          </>
        ) : (
          <Connect connect={this.connect.bind(this)} />
        )}
      </div>
    );
  }
}

export default App;
