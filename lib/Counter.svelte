<script lang="ts">
    let count: number = $state(0);
    const increment = () => {
        count += 1;
    };
</script>

<button onclick={increment} class="bg-blue-500 text-white p-2 cursor-pointer hover:bg-blue-700 rounded">
    count is {count}
</button>
