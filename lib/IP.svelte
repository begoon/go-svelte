<script lang="ts">
    let ip: string = $state("");
    function getIP() {
        fetch("/ip")
            .then((res) => res.text())
            .then((data) => (ip = data));
    }
</script>

<button onclick={getIP} class="bg-green-500 text-white p-2 cursor-pointer hover:bg-green-700 rounded">
    {#if ip === ""}
        Get IP
    {:else}
        Your IP is {ip} (received via SSL/TLS).
    {/if}
</button>
