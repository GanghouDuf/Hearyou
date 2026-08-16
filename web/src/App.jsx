// web/src/App.jsx
import { useState, useEffect, useRef } from 'react'

export default function App() {
	const [messages, setMessages] = useState([])
	const [input, setInput] = useState('')
	const [name, setName] = useState('anon')
	const wsRef = useRef(null)

	useEffect(() => {
		const ws = new WebSocket(`ws://${window.location.host}/ws`)
		wsRef.current = ws

		ws.onmessage = (event) => {
			const msg = JSON.parse(event.data)
			setMessages((prev) => [...prev, msg])
		}

		ws.onopen = () => console.log('WebSocket connected')
		ws.onclose = () => console.log('WebSocket disconnected')

		return () => ws.close()
	}, [])

	function handleSubmit(e) {
		e.preventDefault()
		if (input.trim() === '') return

		wsRef.current.send(JSON.stringify({
			type: 'chat',
			author: name || 'anon',
			payload: input,
		}))
		setInput('')
	}

	return (
		<div style={{ maxWidth: 500, margin: '40px auto', fontFamily: 'sans-serif' }}>
			<h2>Chat</h2>
			<input
				value={name}
				onChange={(e) => setName(e.target.value)}
				placeholder="Ваше имя"
				style={{ width: 120, padding: 6, marginBottom: 8 }}
			/>
			<div style={{ border: '1px solid #ccc', height: 300, overflowY: 'auto', padding: 8, marginBottom: 8 }}>
				{messages.map((msg, i) => (
					<div key={i}>{msg.author}: {msg.payload}</div>
				))}
			</div>
			<form onSubmit={handleSubmit} style={{ display: 'flex', gap: 8 }}>
				<input
					value={input}
					onChange={(e) => setInput(e.target.value)}
					placeholder="Введите сообщение..."
					style={{ flex: 1, padding: 6 }}
				/>
				<button type="submit">Отправить</button>
			</form>
		</div>
	)
}