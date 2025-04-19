import { writable } from 'svelte/store';
import { PUBLIC_SERVER_URL } from "$env/static/public";
function createChatStore() {
	let socket;
	let wsUrl;
	const messages = writable([]);
	const status = writable('disconnected');

	function connect() {
		if (PUBLIC_SERVER_URL) {
			wsUrl = 'ws://localhost:3000/ws';
		} else {
			wsUrl = "wss://" + PUBLIC_SERVER_URL + "/ws"
		}
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