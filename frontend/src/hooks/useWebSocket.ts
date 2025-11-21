import { useEffect, useRef, useState } from 'react'

interface WebSocketMessage {
  type: string
  payload: any
}

export const useWebSocket = () => {
  const [connected, setConnected] = useState(false)
  const [messages, setMessages] = useState<WebSocketMessage[]>([])
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws`
    
    const connect = () => {
      const ws = new WebSocket(wsUrl)
      
      ws.onopen = () => {
        console.log('WebSocket connected')
        setConnected(true)
      }
      
      ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          setMessages((prev) => [...prev, message])
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }
      
      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }
      
      ws.onclose = () => {
        console.log('WebSocket disconnected')
        setConnected(false)
        
        // Reconnect after 3 seconds
        setTimeout(connect, 3000)
      }
      
      wsRef.current = ws
    }
    
    connect()
    
    return () => {
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [])

  const sendMessage = (type: string, payload: any) => {
    if (wsRef.current && connected) {
      wsRef.current.send(JSON.stringify({ type, payload }))
    }
  }

  return { connected, messages, sendMessage }
}
