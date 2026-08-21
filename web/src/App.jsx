// web/src/App.jsx
import { useState, useEffect, useRef } from 'react'

export default function App() {
	const [token, setToken] = useState(null)
	const [mode, setMode] = useState('login') // 'login' | 'register'
	const [username, setUsername] = useState('')
	const [password, setPassword] = useState('')
	const [authError, setAuthError] = useState('')

	const [messages, setMessages] = useState([])
	const [input, setInput] = useState('')
	const wsRef = useRef(null)

	// Подключаемся к WS только когда есть токен
	useEffect(() => {
		if (!token) return

		const ws = new WebSocket(`ws://${window.location.host}/ws?token=${token}`)
		wsRef.current = ws

		ws.onmessage = (event) => {
			const msg = JSON.parse(event.data)
			if (msg.type === 'error') {
				console.error('server error:', msg.payload)
				return
			}
			setMessages((prev) => [...prev, msg])
		}

		ws.onopen = () => console.log('WebSocket connected')
		ws.onclose = () => console.log('WebSocket disconnected')

		return () => ws.close()
	}, [token])

	async function handleAuth(e) {
		e.preventDefault()
		setAuthError('')

		try {
			if (mode === 'register') {
				const res = await fetch('/register', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ username, password }),
				})
				if (!res.ok) {
					setAuthError(await res.text())
					return
				}
				// после регистрации сразу логинимся
			}

			const res = await fetch('/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, password }),
			})
			if (!res.ok) {
				setAuthError(await res.text())
				return
			}
			const data = await res.json()
			setToken(data.token)
		} catch (err) {
			setAuthError('network error: ' + err.message)
		}
	}

	function handleSend(e) {
	e.preventDefault()
	if (input.trim() === '') return

	wsRef.current.send(JSON.stringify({
		type: 'chat',
		payload: input, // без author
	}))
	setInput('')
}

	// --- Экран логина/регистрации ---
	if (!token) {
		return (
			<div style={{ maxWidth: 350, margin: '80px auto', fontFamily: 'sans-serif' }}>
				<h2>{mode === 'login' ? 'Вход' : 'Регистрация'}</h2>
				<form onSubmit={handleAuth} style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
					<input
						value={username}
						onChange={(e) => setUsername(e.target.value)}
						placeholder="Имя пользователя"
					/>
					<input
						type="password"
						value={password}
						onChange={(e) => setPassword(e.target.value)}
						placeholder="Пароль"
					/>
					<button type="submit">{mode === 'login' ? 'Войти' : 'Зарегистрироваться'}</button>
				</form>
				{authError && <p style={{ color: 'red' }}>{authError}</p>}
				<button
					onClick={() => setMode(mode === 'login' ? 'register' : 'login')}
					style={{ marginTop: 8 }}
				>
					{mode === 'login' ? 'Нет аккаунта? Зарегистрироваться' : 'Уже есть аккаунт? Войти'}
				</button>
			</div>
		)
	}

	// --- Экран чата ---
	return (
		<div style={{ maxWidth: 500, margin: '40px auto', fontFamily: 'sans-serif' }}>
			<h2>Chat — вы вошли как {username}</h2>
			<div style={{ border: '1px solid #ccc', height: 300, overflowY: 'auto', padding: 8, marginBottom: 8 }}>
				{messages.map((msg, i) => (
					<div key={i}>{msg.author}: {msg.payload}</div>
				))}
			</div>
			<form onSubmit={handleSend} style={{ display: 'flex', gap: 8 }}>
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