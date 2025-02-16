import { mount } from "svelte";
import App from "./Page.svelte";

const app = mount(App, {
    target: document.getElementById("app")!,
});

export default app;
