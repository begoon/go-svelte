<script lang="ts">
    import Counter from "../lib/Counter.svelte";
    import IP from "../lib/IP.svelte";
    import "./app.css";

    const data = window.__DATA__;
</script>

<main>
    <h1>ROOT</h1>
    <Counter />
    <IP />
    <br />
    ▶︎ <a href="about/" class="underline">about</a>
    ▶︎ <a href="about/ABC" class="underline">about/id</a>
    <div>page data: {JSON.stringify(data)}</div>
    <img src="/images/vite.svg" alt="Vite" class="fixed left-0 bottom-0" />
</main>

<style>
</style>
