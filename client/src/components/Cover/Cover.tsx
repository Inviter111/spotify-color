import React from "react";
import "./Cover.css";
import { IProps } from "./types";

const cover: React.FC<IProps> = (props: IProps) => {
    const { url } = props;

    return <img src={url} alt="cover" className="cover" />;
};

export default cover;
