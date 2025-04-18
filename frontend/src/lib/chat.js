import { writable } from 'svelte/store';

function createChatStore() {
	let socket;
	const messages = writable([]);
	const status = writable('disconnected');

	function connect() {
		const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:3000/ws';

		socket = new WebSocket(wsUrl);

		socket.onopen = () => {
			status.set('connected');
		};
		socket.onclose = () => {
			status.set('disconnected');
		};
		socket.onerror = () => {
			status.set('error');
		};
		socket.onmessage = (event) => {
			messages.update((msgs) => [...msgs, event.data]);
		};
	}

	function send(msg) {
		if (socket && socket.readyState === WebSocket.OPEN) {
			socket.send(msg);
		}
	}

	function disconnect() {
		if (socket) {
			socket.close();
		}
	}

	return {
		connect,
		disconnect,
		send,
		messages,
		status
	};
}

export const chat = createChatStore();