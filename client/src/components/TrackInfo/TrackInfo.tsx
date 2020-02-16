import React from 'react'
import { IProps } from './types'
import './TrackInfo.css'

const trackInfo: React.FC<IProps> = (props: IProps) => {
    const { songName, artistName } = props;

    return (
        <div className="trackInfo">
            <p className="songName">{songName}</p>
            <p className="artistName">{artistName}</p>
        </div>
    )
}

export default trackInfo;
