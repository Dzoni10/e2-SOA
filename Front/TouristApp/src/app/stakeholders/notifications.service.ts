import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class NotificationsService {
  private ws?: WebSocket;
  private _notifications = new Subject<any>();
  private connectionCount = 0; // Broj aktivnih "slušalaca"
  private currentUserId?: string;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectTimeout?: any;

  public notifications$ = this._notifications.asObservable();

  /**
   * Kreira ili održava postojeću WebSocket konekciju
   */
  connect(userId: string): void {
    this.connectionCount++;
    
    // Ako je već povezan sa istim korisnikom, samo povećaj counter
    if (this.ws && this.ws.readyState === WebSocket.OPEN && this.currentUserId === userId) {
      console.log(`WebSocket already connected for user ${userId}. Connection count: ${this.connectionCount}`);
      return;
    }

    // Ako se pokušava povezati sa drugim korisnikom, zatvori postojeću konekciju
    if (this.currentUserId && this.currentUserId !== userId) {
      this.forceDisconnect();
    }

    this.currentUserId = userId;
    this.createWebSocketConnection();
  }

  /**
   * Smanjuje broj aktivnih slušalaca, zatvara konekciju samo ako nema više slušalaca
   */
  disconnect(): void {
    this.connectionCount = Math.max(0, this.connectionCount - 1);
    
    console.log(`Disconnect called. Remaining connections: ${this.connectionCount}`);
    
    // Zatvori konekciju samo ako nema više aktivnih slušalaca
    if (this.connectionCount === 0) {
      this.forceDisconnect();
    }
  }

  /**
   * Prisilno zatvara WebSocket konekciju bez obzira na broj slušalaca
   */
  forceDisconnect(): void {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = undefined;
    }

    if (this.ws) {
      console.log('Closing WebSocket connection');
      this.ws.onclose = null; // Spreci auto-reconnect
      this.ws.close();
      this.ws = undefined;
    }
    
    this.connectionCount = 0;
    this.currentUserId = undefined;
    this.reconnectAttempts = 0;
  }

  /**
   * Kreira novu WebSocket konekciju
   */
  private createWebSocketConnection(): void {
    if (!this.currentUserId) {
      console.error('Cannot create WebSocket connection: no userId provided');
      return;
    }

    try {
      console.log(`Creating WebSocket connection for user: ${this.currentUserId}`);
      this.ws = new WebSocket(`ws://localhost:8070/followers/ws?userId=${this.currentUserId}`);

      this.ws.onopen = () => {
        console.log('WebSocket connected successfully');
        this.reconnectAttempts = 0; // Reset na uspešnu konekciju
      };

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          console.log('Received WebSocket message:', data);
          this._notifications.next(data);
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      this.ws.onclose = (event) => {
        console.log('WebSocket closed:', event.code, event.reason);
        this.ws = undefined;
        
        // Pokušaj reconnect samo ako ima aktivnih slušalaca i nije dosegnut limit
        if (this.connectionCount > 0 && this.reconnectAttempts < this.maxReconnectAttempts) {
          this.scheduleReconnect();
        } else {
          console.log('WebSocket will not reconnect (no active listeners or max attempts reached)');
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

    } catch (error) {
      console.error('Error creating WebSocket connection:', error);
    }
  }

  /**
   * Zakazuje pokušaj ponovnog povezivanja
   */
  private scheduleReconnect(): void {
    this.reconnectAttempts++;
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000); // Exponential backoff, max 30s
    
    console.log(`Scheduling reconnect attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts} in ${delay}ms`);
    
    this.reconnectTimeout = setTimeout(() => {
      if (this.connectionCount > 0 && this.currentUserId) {
        this.createWebSocketConnection();
      }
    }, delay);
  }

  /**
   * Vraća status konekcije
   */
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  /**
   * Vraća broj aktivnih slušalaca
   */
  getConnectionCount(): number {
    return this.connectionCount;
  }
}