import { useEffect, useRef, useState } from 'react'
import './App.css'

const Icon = ({ name, size = 20 }) => {
  const paths = {
    logo: <><path d="M5 5.5A3.5 3.5 0 0 1 8.5 2h7A3.5 3.5 0 0 1 19 5.5v5A3.5 3.5 0 0 1 15.5 14H13l-4.5 4v-4H8.5A3.5 3.5 0 0 1 5 10.5z" /><path d="M9 7.75h6M9 10.75h3.5" /></>,
    search: <><circle cx="11" cy="11" r="5.5" /><path d="m15.2 15.2 3.3 3.3" /></>,
    more: <><circle cx="6" cy="12" r="1" fill="currentColor" /><circle cx="12" cy="12" r="1" fill="currentColor" /><circle cx="18" cy="12" r="1" fill="currentColor" /></>,
    smile: <><circle cx="12" cy="12" r="8.5" /><path d="M8.5 14.25s1.2 2 3.5 2 3.5-2 3.5-2M9 9.5h.01M15 9.5h.01" /></>,
    paperclip: <path d="m19 11.5-7.6 7.6a4.25 4.25 0 0 1-6-6L13 5.5a2.75 2.75 0 0 1 3.9 3.9l-7.5 7.5a1.25 1.25 0 0 1-1.8-1.8l6.8-6.8" />,
    send: <path d="m20 4-7.2 16-3.1-7.1L4 9.8zM9.7 12.9 13.2 9.4" />,
    arrow: <path d="m9 18 6-6-6-6" />,
  }
  return <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">{paths[name]}</svg>
}

function Avatar({ name, className = '' }) { return <div className={`avatar ${className}`}>{name?.slice(0, 1).toUpperCase() || 'H'}</div> }

export default function App() {
  const [token, setToken] = useState(null)
  const [mode, setMode] = useState('login')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [authError, setAuthError] = useState('')
  const [messages, setMessages] = useState([])
  const [input, setInput] = useState('')
  const [connected, setConnected] = useState(false)
  const wsRef = useRef(null)
  const messagesEndRef = useRef(null)
useEffect(() => {
	if (!token) return undefined

	const ws = new WebSocket(`ws://${window.location.host}/ws?token=${token}&room=general`)
	wsRef.current = ws

	ws.onmessage = (event) => {
		const msg = JSON.parse(event.data)
		if (msg.type !== 'error') setMessages((previous) => [...previous, msg])
	}
	ws.onopen = () => setConnected(true)
	ws.onclose = () => setConnected(false)

	return () => {
		ws.close()
	}
}, [token])
useEffect(() => {
	messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
	return undefined
}, [messages])
  async function handleAuth(event) { event.preventDefault(); setAuthError(''); try { if (mode === 'register') { const register = await fetch('/register', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ username, password }) }); if (!register.ok) { setAuthError(await register.text()); return } } const login = await fetch('/login', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ username, password }) }); if (!login.ok) { setAuthError(await login.text()); return } const data = await login.json(); setToken(data.token) } catch (error) { setAuthError(`Не удалось подключиться: ${error.message}`) } }
  function handleSend(event) { event.preventDefault(); const text = input.trim(); if (!text || wsRef.current?.readyState !== WebSocket.OPEN) return; wsRef.current.send(JSON.stringify({ type: 'chat', payload: text })); setInput('') }
  if (!token) return <AuthScreen mode={mode} setMode={setMode} username={username} setUsername={setUsername} password={password} setPassword={setPassword} error={authError} onSubmit={handleAuth} />
  return <main className="messenger-shell"><aside className="sidebar"><div className="sidebar-top"><div className="brand-mark"><Icon name="logo" size={23} /></div><span className="brand-name">hear<span>you</span></span><button className="icon-button sidebar-more" aria-label="Меню"><Icon name="more" /></button></div><label className="search"><Icon name="search" size={18} /><input placeholder="Поиск" /></label><section className="chat-list"><p className="section-label">СООБЩЕНИЯ</p><button className="conversation active"><Avatar name="H" className="room-avatar" /><span className="conversation-copy"><b>Общий чат</b><small>{connected ? 'В сети · можно писать' : 'Подключаемся...'}</small></span><span className="unread-dot" /></button></section><div className="profile"><Avatar name={username} /><span><b>{username}</b><small><i className={connected ? 'online' : ''} />{connected ? 'в сети' : 'нет соединения'}</small></span><button className="icon-button" aria-label="Настройки"><Icon name="more" /></button></div></aside><section className="chat-panel"><header className="chat-header"><div><div className="chat-title"><Avatar name="H" className="room-avatar tiny" /><span>Общий чат</span></div><p><i className={connected ? 'online' : ''} />{connected ? 'Онлайн' : 'Соединение...'}</p></div><div className="header-actions"><button className="icon-button" aria-label="Поиск в чате"><Icon name="search" /></button><button className="icon-button" aria-label="Больше действий"><Icon name="more" /></button></div></header><div className="messages-area"><div className="date-divider"><span>Сегодня</span></div>{messages.length === 0 ? <div className="empty-chat"><div className="empty-icon"><Icon name="logo" size={30} /></div><h2>Начните разговор</h2><p>Поздоровайтесь — ваше сообщение увидят все в этом чате.</p></div> : messages.map((message, index) => { const own = message.author === username; return <div className={`message-row ${own ? 'own' : ''}`} key={`${message.author}-${index}`}><Avatar name={message.author} /><div className="message"><div className="message-meta"><b>{own ? 'Вы' : message.author}</b><time>{new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</time></div><p>{message.payload}</p></div></div> })}<div ref={messagesEndRef} /></div><form className="composer" onSubmit={handleSend}><button type="button" className="icon-button" aria-label="Прикрепить файл"><Icon name="paperclip" /></button><input value={input} onChange={(event) => setInput(event.target.value)} placeholder={connected ? 'Напишите сообщение...' : 'Ожидание соединения...'} disabled={!connected} /><button type="button" className="icon-button" aria-label="Добавить эмодзи"><Icon name="smile" /></button><button className="send-button" aria-label="Отправить" disabled={!input.trim() || !connected}><Icon name="send" size={19} /></button></form></section></main>
}

function AuthScreen({ mode, setMode, username, setUsername, password, setPassword, error, onSubmit }) { const isLogin = mode === 'login'; return <main className="auth-page"><div className="auth-ambient one" /><div className="auth-ambient two" /><section className="auth-card"><div className="auth-logo"><Icon name="logo" size={31} /></div><p className="eyebrow">HEARYOU MESSENGER</p><h1>{isLogin ? 'Рады видеть' : 'Создайте аккаунт'}</h1><p className="auth-subtitle">{isLogin ? 'Войдите, чтобы продолжить общение.' : 'Пара минут — и вы в общем чате.'}</p><form onSubmit={onSubmit} className="auth-form"><label>Имя пользователя<input value={username} onChange={(event) => setUsername(event.target.value)} placeholder="Например, alex" autoComplete="username" required /></label><label>Пароль<input type="password" value={password} onChange={(event) => setPassword(event.target.value)} placeholder="Не менее 6 символов" autoComplete={isLogin ? 'current-password' : 'new-password'} required /></label>{error && <p className="auth-error">{error}</p>}<button className="primary-button" type="submit">{isLogin ? 'Войти в HearYou' : 'Зарегистрироваться'}<Icon name="arrow" size={18} /></button></form><p className="auth-switch">{isLogin ? 'Впервые здесь?' : 'Уже есть аккаунт?'} <button type="button" onClick={() => setMode(isLogin ? 'register' : 'login')}>{isLogin ? 'Создать аккаунт' : 'Войти'}</button></p></section><p className="auth-footer">Простое пространство для настоящих разговоров</p></main> }
