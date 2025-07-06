<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>Whats Valid</title>
	<meta
		name="description"
		content="Verify any phone number's WhatsApp registration status in seconds. Fast, accurate, and reliable validation for businesses and individuals."
	/>
	<meta name="theme-color" content="#18181b" />

	<!-- Facebook Meta Tags -->
	<!-- <meta property="og:url" content="TODO" /> -->
	<meta property="og:type" content="website" />
	<meta property="og:title" content="Whats Valid" />
	<meta
		property="og:description"
		content="Verify any phone number's WhatsApp registration status in seconds. Fast, accurate, and reliable validation for businesses and individuals."
	/>
	<meta property="og:image" content="/opengraph.png" />
	<meta property="og:site_name" content="Dateng Ya!" />

	<!-- Twitter Meta Tags -->
	<meta name="twitter:card" content="summary_large_image" />
	<!-- <meta property="twitter:domain" content="TODO" /> -->
	<!-- <meta property="twitter:url" content="TODO" /> -->
	<meta name="twitter:title" content="Whats Valid" />
	<meta
		name="twitter:description"
		content="Verify any phone number's WhatsApp registration status in seconds. Fast, accurate, and reliable validation for businesses and individuals."
	/>
	<meta name="twitter:image" content="/opengraph.png" />
</svelte:head>

<main class="relative grid h-dvh place-items-center bg-zinc-900 px-5 text-zinc-200">
	{@render children()}
	<footer class="absolute bottom-3 left-1/2 -translate-x-1/2 text-sm text-zinc-400">
		<a href="https://github.com/wiscaksono" target="_blank" class="underline underline-offset-2">
			GitHub
		</a>
		•
		<a href="https://wiscaksono.com/" target="_blank" class="underline underline-offset-2">
			Wiscaksono
		</a>
	</footer>
</main>
