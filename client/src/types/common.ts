export interface ISettings {
    enableLyrics: boolean;
}

export interface ConfigMessage {
    type: 'change' | 'getState' | 'sendState';
    data?: ISettings;
}