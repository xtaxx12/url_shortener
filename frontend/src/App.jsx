import { useState, useEffect } from 'react'

const API_URL = import.meta.env.VITE_API_URL || '/api'

function App() {
    const [url, setUrl] = useState('')
    const [shortUrl, setShortUrl] = useState(null)
    const [stats, setStats] = useState(null)
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState(null)
    const [copied, setCopied] = useState(false)
    const [recentUrls, setRecentUrls] = useState([])

    useEffect(() => {
        const saved = localStorage.getItem('recentUrls')
        if (saved) {
            setRecentUrls(JSON.parse(saved))
        }
    }, [])

    const shortenUrl = async (e) => {
        e.preventDefault()
        if (!url.trim()) return

        setLoading(true)
        setError(null)
        setShortUrl(null)
        setStats(null)

        try {
            const response = await fetch(`${API_URL}/shorten`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url: url.trim() }),
            })

            if (!response.ok) {
                const errorData = await response.json()
                throw new Error(errorData.error || 'Error al acortar la URL')
            }

            const data = await response.json()
            setShortUrl(data)

            const newRecent = [
                { ...data, created_at: new Date().toISOString() },
                ...recentUrls.filter(r => r.short_code !== data.short_code).slice(0, 4)
            ]
            setRecentUrls(newRecent)
            localStorage.setItem('recentUrls', JSON.stringify(newRecent))

            fetchStats(data.short_code)
        } catch (err) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    const fetchStats = async (code) => {
        try {
            const response = await fetch(`${API_URL}/stats/${code}`)
            if (response.ok) {
                const data = await response.json()
                setStats(data)
            }
        } catch (err) {
            console.error('Error fetching stats:', err)
        }
    }

    const copyToClipboard = async () => {
        if (!shortUrl?.short_url) return

        try {
            await navigator.clipboard.writeText(shortUrl.short_url)
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        } catch (err) {
            console.error('Failed to copy:', err)
        }
    }

    const refreshStats = (code) => {
        fetchStats(code)
    }

    return (
        <div className="min-h-screen flex flex-col">
            <div className="absolute inset-0 overflow-hidden pointer-events-none">
                <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary-500/20 rounded-full filter blur-3xl animate-pulse-slow" />
                <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-purple-500/20 rounded-full filter blur-3xl animate-pulse-slow" style={{ animationDelay: '1s' }} />
            </div>

            <header className="relative z-10 py-6 px-4">
                <div className="max-w-4xl mx-auto flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-gradient-to-br from-primary-400 to-primary-600 rounded-xl flex items-center justify-center">
                            <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                            </svg>
                        </div>
                        <h1 className="text-2xl font-bold text-white">URL Shortener</h1>
                    </div>
                    <div className="flex items-center gap-2">
                        <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-green-500/20 text-green-400 border border-green-500/30">
                            <span className="w-2 h-2 bg-green-400 rounded-full mr-2 animate-pulse" />
                            API Online
                        </span>
                    </div>
                </div>
            </header>

            <main className="relative z-10 flex-1 flex flex-col items-center justify-center px-4 py-12">
                <div className="max-w-2xl w-full text-center mb-12">
                    <h2 className="text-4xl md:text-5xl font-extrabold text-white mb-4">
                        Acorta tus enlaces
                        <span className="block text-transparent bg-clip-text bg-gradient-to-r from-primary-400 to-purple-400">
                            en segundos
                        </span>
                    </h2>
                    <p className="text-gray-400 text-lg">
                        Crea enlaces cortos profesionales y rastrea sus estadísticas en tiempo real
                    </p>
                </div>

                <form onSubmit={shortenUrl} className="w-full max-w-2xl mb-8">
                    <div className="glass-card p-2 flex flex-col sm:flex-row gap-2">
                        <input
                            type="url"
                            value={url}
                            onChange={(e) => setUrl(e.target.value)}
                            placeholder="https://ejemplo.com/mi-url-muy-larga"
                            className="input-glass flex-1 bg-transparent border-0"
                            required
                        />
                        <button
                            type="submit"
                            disabled={loading || !url.trim()}
                            className="btn-primary whitespace-nowrap flex items-center justify-center gap-2"
                        >
                            {loading ? (
                                <>
                                    <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                                    </svg>
                                    Acortando
                                </>
                            ) : (
                                <>
                                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                                    </svg>
                                    Acortar URL
                                </>
                            )}
                        </button>
                    </div>
                </form>

                {error && (
                    <div className="w-full max-w-2xl mb-6 p-4 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 flex items-center gap-3">
                        <svg className="w-5 h-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        {error}
                    </div>
                )}

                {shortUrl && (
                    <div className="w-full max-w-2xl glass-card p-6 mb-8 animate-in">
                        <div className="flex items-center justify-between mb-4">
                            <h3 className="text-lg font-semibold text-white flex items-center gap-2">
                                <svg className="w-5 h-5 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                                ¡URL Acortada!
                            </h3>
                            <button
                                onClick={() => refreshStats(shortUrl.short_code)}
                                className="text-gray-400 hover:text-white transition-colors p-2 rounded-lg hover:bg-white/10"
                                title="Actualizar estadísticas"
                            >
                                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                                </svg>
                            </button>
                        </div>

                        <div className="flex items-center gap-3 p-4 bg-white/5 rounded-xl mb-4">
                            <input
                                type="text"
                                value={shortUrl.short_url}
                                readOnly
                                className="flex-1 bg-transparent text-primary-400 font-mono text-lg focus:outline-none"
                            />
                            <button
                                onClick={copyToClipboard}
                                className={`p-2 rounded-lg transition-all duration-200 ${copied
                                        ? 'bg-green-500/20 text-green-400'
                                        : 'bg-white/10 text-white hover:bg-white/20'
                                    }`}
                            >
                                {copied ? (
                                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                    </svg>
                                ) : (
                                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                                    </svg>
                                )}
                            </button>
                        </div>

                        <p className="text-gray-400 text-sm truncate">
                            <span className="text-gray-500">Original:</span> {shortUrl.original_url}
                        </p>

                        {stats && (
                            <div className="grid grid-cols-3 gap-4 mt-6 pt-6 border-t border-white/10">
                                <div className="stat-card !p-4">
                                    <div className="text-3xl font-bold text-white mb-1">{stats.clicks}</div>
                                    <div className="text-gray-400 text-sm">Clics</div>
                                </div>
                                <div className="stat-card !p-4">
                                    <div className="text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-primary-400 to-purple-400">
                                        {shortUrl.short_code}
                                    </div>
                                    <div className="text-gray-400 text-sm">Código</div>
                                </div>
                                <div className="stat-card !p-4">
                                    <div className="text-lg font-bold text-white mb-1">
                                        {new Date(stats.created_at).toLocaleDateString('es-ES', { month: 'short', day: 'numeric' })}
                                    </div>
                                    <div className="text-gray-400 text-sm">Creado</div>
                                </div>
                            </div>
                        )}
                    </div>
                )}

                {recentUrls.length > 0 && !shortUrl && (
                    <div className="w-full max-w-2xl">
                        <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                            <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            URLs Recientes
                        </h3>
                        <div className="space-y-3">
                            {recentUrls.map((item) => (
                                <div key={item.short_code} className="glass-card p-4 flex items-center justify-between">
                                    <div className="flex-1 min-w-0">
                                        <p className="text-primary-400 font-mono">{item.short_url}</p>
                                        <p className="text-gray-500 text-sm truncate">{item.original_url}</p>
                                    </div>
                                    <button
                                        onClick={() => {
                                            setShortUrl(item)
                                            fetchStats(item.short_code)
                                        }}
                                        className="ml-4 p-2 rounded-lg bg-white/10 text-white hover:bg-white/20 transition-colors"
                                    >
                                        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                                        </svg>
                                    </button>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </main>

            <footer className="relative z-10 py-6 px-4 border-t border-white/10">
                <div className="max-w-4xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4 text-gray-500 text-sm">
                    <div className="flex items-center gap-2">
                        <span>Construido con</span>
                        <span className="text-red-400">♥</span>
                        <span>usando Go, React & Docker</span>
                    </div>
                    <div className="flex items-center gap-4">
                        <a href="https://github.com" className="hover:text-white transition-colors flex items-center gap-2">
                            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                                <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
                            </svg>
                            GitHub
                        </a>
                    </div>
                </div>
            </footer>
        </div>
    )
}

export default App
