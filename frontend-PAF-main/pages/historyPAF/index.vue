<script setup>
    const {data: posts} = await useFetch('https://jsonplaceholder.typicode.com/posts')
</script>

<template>
    <h1 class="text-3xl text-center">Historial de PAF</h1>
    <section>
        <div v-for="{ title, id } in posts" :key="id">
            {{ title }}
        </div>
    </section>
</template>