<script>
    import { onMount, onDestroy } from 'svelte';
    import { chat } from '$lib/chat.js';

    let input = '';
    let autoscroll;

    // Destructure Svelte stores from chat and assign to local vars
    const { messages, status } = chat;

    onMount(() => {
        chat.connect();
    });

    onDestroy(() => {
        chat.disconnect();
    });

    function sendMsg() {
        if (input.trim().length) {
            chat.send(input);
            input = '';
        }
    }

    // Auto-scroll when messages change
    $: if (autoscroll && $messages) {
        setTimeout(() => {
            autoscroll.scrollTo(0, autoscroll.scrollHeight);
        }, 0);
    }
</script>

<div class="chat-container" style="max-width: 500px; margin: auto;">
    <h2>WebSocket Chat</h2>
    <div class="status-bar">
        Status: <span>{$status}</span>
    </div>
    <div class="chat-messages" bind:this={autoscroll} style="height: 300px; overflow-y: auto; border: 1px solid #ddd; margin-bottom: 8px;">
        {#each $messages as msg}
            <div class="chat-message" style="padding: 2px 4px;">
                {msg}
            </div>
        {/each}
    </div>
    <form on:submit|preventDefault={sendMsg}>
        <input
            type="text"
            bind:value={input}
            placeholder="Type your message..."
            class="chat-input"
            autocomplete="off"
            style="width: 85%;"
        />
        <button type="submit" style="width: 13%;">Send</button>
    </form>
</div>
