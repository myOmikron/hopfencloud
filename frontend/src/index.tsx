import React from "react";
import ReactDOM from "react-dom/client";
import { ToastContainer } from "react-toastify";
import "./index.css";
import "react-toastify/dist/ReactToastify.css";

type RouterProps = {};
type RouterState = {};

class Router extends React.Component<RouterProps, RouterState> {
    constructor(props: RouterProps) {
        super(props);

        this.state = {};
    }

    render() {
        return <div></div>;
    }
}

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);
root.render(
    <>
        <Router />
        <ToastContainer />
    </>
);
