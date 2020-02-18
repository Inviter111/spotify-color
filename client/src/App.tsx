import React from "react";

import "./App.css";

import Config from "./pages/Config";
import Index from "./pages/Index";

import { useRoutes } from "hookrouter";

const routes = {
    "/": () => <Index />,
    "/config": () => <Config />
};

function App() {
    const routeResult = useRoutes(routes);
    return <div className="App">{routeResult}</div>;
}

export default App;
