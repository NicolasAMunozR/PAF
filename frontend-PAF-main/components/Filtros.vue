<template>
    <div class="space-y-2">
      <details class="overflow-hidden rounded border border-gray-300">
        <summary class="flex cursor-pointer items-center justify-between gap-2 bg-white p-4 text-gray-900 transition">
          <span class="text-sm font-medium">Availability</span>
          <span class="transition group-open:-rotate-180">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="size-4">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
            </svg>
          </span>
        </summary>
        <div class="border-t border-gray-200 bg-white">
          <header class="flex items-center justify-between p-4">
            <span class="text-sm text-gray-700">{{ selectedCount }} Selected</span>
            <button type="button" @click="resetFilters" class="text-sm text-gray-900 underline">Reset</button>
          </header>
          <ul class="space-y-1 border-t border-gray-200 p-4">
            <li v-for="(option, index) in availabilityOptions" :key="index">
              <label class="inline-flex items-center gap-2">
                <input type="checkbox" v-model="selectedOptions" :value="option.value" class="size-5 rounded border-gray-300" />
                <span class="text-sm font-medium text-gray-700">{{ option.label }}</span>
              </label>
            </li>
          </ul>
        </div>
      </details>
  
      <details class="overflow-hidden rounded border border-gray-300">
        <summary class="flex cursor-pointer items-center justify-between gap-2 bg-white p-4 text-gray-900 transition">
          <span class="text-sm font-medium">Price</span>
          <span class="transition group-open:-rotate-180">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="size-4">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
            </svg>
          </span>
        </summary>
        <div class="border-t border-gray-200 bg-white">
          <header class="flex items-center justify-between p-4">
            <span class="text-sm text-gray-700">The highest price is $600</span>
            <button type="button" @click="resetPrice" class="text-sm text-gray-900 underline">Reset</button>
          </header>
          <div class="border-t border-gray-200 p-4">
            <div class="flex justify-between gap-4">
              <input v-model.number="priceRange.from" type="number" placeholder="From" class="w-full rounded-md border-gray-200 shadow-sm sm:text-sm" />
              <input v-model.number="priceRange.to" type="number" placeholder="To" class="w-full rounded-md border-gray-200 shadow-sm sm:text-sm" />
            </div>
          </div>
        </div>
      </details>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        selectedOptions: [],
        availabilityOptions: [
          { label: 'In Stock (5+)', value: 'inStock' },
          { label: 'Pre Order (3+)', value: 'preOrder' },
          { label: 'Out of Stock (10+)', value: 'outOfStock' }
        ],
        priceRange: { from: '', to: '' }
      };
    },
    computed: {
      selectedCount() {
        return this.selectedOptions.length;
      }
    },
    methods: {
      resetFilters() {
        this.selectedOptions = [];
      },
      resetPrice() {
        this.priceRange = { from: '', to: '' };
      }
    },
    watch: {
      selectedOptions(newVal) {
        this.$emit('filter-changed', { availability: newVal, price: this.priceRange });
      },
      priceRange: {
        deep: true,
        handler(newVal) {
          this.$emit('filter-changed', { availability: this.selectedOptions, price: newVal });
        }
      }
    }
  };
  </script>
  