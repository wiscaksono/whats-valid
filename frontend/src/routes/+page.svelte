<script lang="ts">
	import type { WhatsValidRes } from '$lib/type';

	let data = $state<WhatsValidRes>();
	let number = $state<string>('');
	let loading = $state(false);
	let getDebouncedNumber = debounced(() => number, 250);

	function debounce<T>(f: (...args: T[]) => unknown, ms: number) {
		let id: null | number = null;
		return (...args: T[]) => {
			if (id) {
				clearTimeout(id);
			}
			id = setTimeout(() => {
				f(...args);
			}, ms);
		};
	}

	function debounced<T>(stateGetter: () => T, ms: number) {
		let state = $state(stateGetter());
		const update = debounce<T>((v) => (state = v), ms);
		$effect(() => update(stateGetter()));
		return () => state;
	}

	$effect(() => {
		const number = getDebouncedNumber();
		if (!number) return;

		(async () => {
			try {
				loading = true;
				const response = await fetch(`http://localhost:3000/check?number=${number}`);
				data = (await response.json()) as WhatsValidRes;
			} catch (error) {
				console.error(error);
			} finally {
				loading = false;
			}
		})();
	});
</script>

<div class="flex max-w-[36rem] flex-col items-center justify-center gap-5 text-center">
	<div class="space-y-2">
		<h1 class="text-4xl font-semibold">Whats Valid</h1>
		<p>
			Verify any phone number's WhatsApp registration status in seconds. Fast, accurate, and
			reliable validation for businesses and individuals.
		</p>
	</div>
	<div class="relative w-full">
		<input
			type="tel"
			inputmode="numeric"
			class="w-full rounded border border-zinc-600 bg-zinc-800 px-3 py-2 focus:border-transparent focus:ring-2 focus:ring-zinc-400 focus:outline-none"
			bind:value={number}
			placeholder="Enter a phone number (e.g. +6287885002327)"
			autocomplete="tel"
		/>
		{#if loading}
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="24"
				height="24"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
				class="absolute top-1/2 right-3 size-6 -translate-y-1/2 animate-spin text-zinc-400"
			>
				<path d="M21 12a9 9 0 1 1-6.219-8.56" />
			</svg>
		{/if}
	</div>
	{#if data}
		<div
			class="relative w-full overflow-hidden rounded border border-zinc-600 bg-zinc-800 p-3 text-left text-zinc-200"
		>
			{#if data?.isOnWhatsApp && !loading}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="absolute top-1/2 right-3 hidden size-12 -translate-y-1/2 text-[#46B37F] sm:block"
				>
					<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
					<path d="m9 11 3 3L22 4"></path>
				</svg>
			{:else}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="absolute top-1/2 right-3 hidden size-12 -translate-y-1/2 text-red-500 sm:block"
				>
					<circle cx="12" cy="12" r="10"></circle>
					<path d="m15 9-6 6"></path>
					<path d="m9 9 6 6"></path>
				</svg>
			{/if}
			<div>
				<p class={`font-medium ${data?.isOnWhatsApp ? 'text-[#46B37F]' : 'text-red-500'}`}>
					{data?.isOnWhatsApp
						? 'This number is registered on WhatsApp'
						: 'This number is not registered on WhatsApp'}
				</p>
				<p class="text-sm">Number: {data?.number}</p>
			</div>
		</div>
	{/if}
</div>
