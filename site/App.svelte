<script lang="ts">
    import Counter from "../lib/Counter.svelte";
    import "./app.css";

    const data = window.__DATA__;
</script>

<main>
    <h1>ROOT</h1>
    <Counter />
    <a href="about/">▶︎ about</a>
    <a href="about/ABC">▶︎ parameterised</a>
    <div>{JSON.stringify(data)}</div>
    <img src="/images/vite.svg" alt="Vite" class="fixed left-0 bottom-0" />
</main>

<style>
</style>
